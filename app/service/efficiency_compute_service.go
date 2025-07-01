package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/global/error_handler"
	"wagner/app/http/vo"
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/service/calc/calc_node"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/lock"
	"wagner/app/utils/log"
	"wagner/infrastructure/persistence/entity"
)

// 人效计算服务
type EfficiencyComputeService struct {
}

// 人效计算参数
type ComputeParams struct {
}

func CreateEfficiencyComputeService() *EfficiencyComputeService {
	return &EfficiencyComputeService{}
}

func (service *EfficiencyComputeService) TimeOnTask(employeeNumber string, operateDay time.Time) (*vo.TimeOnTaskVO, *business_error.BusinessError) {
	ctx, calcParam, err := service.createAndComputeCtx(employeeNumber, operateDay)
	if err != nil {
		return nil, err
	}
	timeOnTaskVO := service.buildTimeOnTask(ctx, calcParam.CalcOtherParam.Work.WorkLoadUnits)

	return timeOnTaskVO, nil
}

func (service *EfficiencyComputeService) ComputeEmployee(employeeNumber string, operateDay time.Time) (bool, *business_error.BusinessError) {
	// 根据计算粒度分布式加锁
	if lockSuccess, err := lock.Lock(employeeNumber, operateDay, 2); err != nil {
		error_handler.LogAndPanic(business_error.LockFailureBySystemError(err))
		return false, nil
	} else if !lockSuccess {
		log.ComputeLogger.Info("员工人效数据计算，加锁失败", "number", employeeNumber)
		return false, business_error.LockFailure()
	} else {
		log.ComputeLogger.Info("员工人效数据计算，加锁成功", "number", employeeNumber)
	}

	ctx, calcParam, err := service.createAndComputeCtx(employeeNumber, operateDay)

	if err != nil {
		return false, err
	}

	// 处理聚合
	for _, storage := range calcParam.SinkStorages {
		switch storage.SinkType {
		case calc_dynamic_param.SUMMARY:
			service.handleSummary(ctx, storage, calcParam.CalcOtherParam)
		case calc_dynamic_param.EMPLOYEE_STATUS:
			service.handleEmployeeStatus(ctx)
		}
	}

	// 根据计算粒度分布式解锁
	if unlockSuccess, err := lock.Unlock(employeeNumber, operateDay); err != nil {
		error_handler.LogAndPanic(business_error.UnlockFailureBySystemError(err))
		return false, nil
	} else if !unlockSuccess {
		log.ComputeLogger.Info("员工人效数据计算，解锁失败", "number", employeeNumber)
		return false, business_error.UnlockFailure()
	} else {
		log.ComputeLogger.Info("员工人效数据计算，解锁成功", "number", employeeNumber)
		return true, nil
	}
}

func (service *EfficiencyComputeService) buildTimeOnTask(ctx *domain.ComputeContext, workLoadUnits []calc_dynamic_param.WorkLoadUnit) *vo.TimeOnTaskVO {
	timeOnTaskVO := vo.TimeOnTaskVO{
		EmployeeNumber: ctx.Employee.Number,
		EmployeeName:   ctx.Employee.Name,
		WorkplaceName:  ctx.Workplace.Name,
		RegionCode:     ctx.Workplace.RegionCode,
		OperateDay:     ctx.OperateDay,
	}

	if ctx.TodayAttendance != nil {
		timeOnTaskVO.Attendance = &vo.AttendanceVO{
			ActionType: domain.ATTENDANCE,
			StartTime:  *ctx.TodayAttendance.ComputedStartTime,
			EndTime:    *ctx.TodayAttendance.ComputedEndTime,
		}
	}

	if ctx.TodayScheduling != nil {
		timeOnTaskVO.Scheduling = &vo.SchedulingVO{
			ActionType: domain.SCHEDULING,
			StartTime:  *ctx.TodayScheduling.ComputedStartTime,
			EndTime:    *ctx.TodayScheduling.ComputedEndTime,
		}

		if ctx.TodayRestList != nil && len(ctx.TodayRestList) > 0 {
			restList := make([]vo.RestVO, len(ctx.TodayRestList))
			for _, rest := range ctx.TodayRestList {
				restList = append(restList, vo.RestVO{
					ActionType: domain.REST,
					StartTime:  *rest.StartTime,
					EndTime:    *rest.EndTime,
				})
			}

			timeOnTaskVO.Scheduling.RestList = restList
		}
	}

	workLoadUnitCode2Name := make(map[string]string)
	workLoadCodeList := make([]string, 0)
	for _, workLoadUnit := range workLoadUnits {
		workLoadUnitCode2Name[workLoadUnit.Code] = workLoadUnit.Name
		workLoadCodeList = append(workLoadCodeList, workLoadUnit.Code)
	}

	processDurationList := service.buildProcessDurationList(ctx.TodayWorkList, ctx.Workplace.Name, workLoadUnitCode2Name, workLoadCodeList)
	timeOnTaskVO.ProcessDurationList = processDurationList

	return &timeOnTaskVO
}

func (service *EfficiencyComputeService) buildProcessDurationList(todayWorkList []domain.Actionable, workplaceName string, workLoadUnitCode2Name map[string]string, workLoadCodeList []string) []*vo.ProcessDurationVO {
	workList := todayWorkList
	processDurationList := make([]*vo.ProcessDurationVO, 0)
	if workList == nil || len(workList) == 0 {
		return processDurationList
	}

	var currentDuration = service.initProcessDuration(workList[0], workplaceName)

	for i := 1; i < len(workList); i++ {
		currentWork := workList[i]

		if currentWork.GetAction().ActionType == domain.IDLE {
			service.buildWorkLoadDesc(currentDuration, workLoadUnitCode2Name, workLoadCodeList)
			processDurationList = append(processDurationList, currentDuration)
			currentDuration = service.initProcessDuration(currentWork, workplaceName)
		} else {
			// 当currentWork的环节与currentDuration的环节相同时
			if currentWork.GetAction().ProcessCode == currentDuration.ProcessCode && currentWork.GetAction().ActionType == currentDuration.ActionType {
				service.mergeProcessDuration(currentDuration, currentWork)
			} else {
				// 当currentWork的环节与currentDuration的环节不同时
				service.buildWorkLoadDesc(currentDuration, workLoadUnitCode2Name, workLoadCodeList)
				processDurationList = append(processDurationList, currentDuration)

				currentDuration = service.initProcessDuration(currentWork, workplaceName)
			}
		}

	}

	if currentDuration != nil {
		service.buildWorkLoadDesc(currentDuration, workLoadUnitCode2Name, workLoadCodeList)
		processDurationList = append(processDurationList, currentDuration)
	}

	return processDurationList
}

func (service *EfficiencyComputeService) initProcessDuration(actionable domain.Actionable, workplaceName string) *vo.ProcessDurationVO {
	processDurationVO := vo.ProcessDurationVO{
		Id:            fmt.Sprintf("%v-%v", datetime_util.FormatDatetime(*actionable.GetAction().ComputedStartTime), actionable.GetAction().ProcessCode),
		ProcessCode:   actionable.GetAction().ProcessCode,
		ProcessName:   actionable.GetAction().Process.Name,
		ActionType:    actionable.GetAction().ActionType,
		StartTime:     *actionable.GetAction().ComputedStartTime,
		EndTime:       *actionable.GetAction().ComputedEndTime,
		WorkplaceName: workplaceName, // 为当天多工作点作业做准备
		WorkLoad:      make(map[string]float64),
		Details:       make([]vo.ProcessDurationDetailVO, 0),
	}
	service.mergeProcessDuration(&processDurationVO, actionable)
	return &processDurationVO
}

func (service *EfficiencyComputeService) mergeProcessDuration(current *vo.ProcessDurationVO, work domain.Actionable) {
	diff := work.GetAction().ComputedEndTime.Sub(*work.GetAction().ComputedStartTime)
	duration := math.Round(diff.Seconds() / 60)
	current.Duration += duration

	detail := vo.ProcessDurationDetailVO{
		StartTime: datetime_util.FormatDatetime(*work.GetAction().ComputedStartTime),
		EndTime:   datetime_util.FormatDatetime(*work.GetAction().ComputedEndTime),
		Duration:  duration,
	}
	if work.GetAction().OperationMsgList != nil {
		detail.OperationMessage = strings.Join(work.GetAction().OperationMsgList, "\n")
	}
	current.Details = append(current.Details, detail)

	if work.GetAction().ActionType == domain.DIRECT_WORK {
		directWork := work.(*domain.DirectWork)
		for key, thisValue := range current.WorkLoad {
			thatValue := directWork.WorkLoad[key]
			current.WorkLoad[key] = thisValue + thatValue
			delete(directWork.WorkLoad, key)
		}

		for key, thatValue := range directWork.WorkLoad {
			current.WorkLoad[key] = thatValue
			delete(directWork.WorkLoad, key)
		}
	}
}

func (service *EfficiencyComputeService) buildWorkLoadDesc(current *vo.ProcessDurationVO, workLoadUnitCode2Name map[string]string, workLoadCodeList []string) {
	if len(current.WorkLoad) == 0 {
		return
	}
	workLoadDescList := make([]string, 0)
	for _, code := range workLoadCodeList {
		if value, exists := current.WorkLoad[code]; exists {
			name := workLoadUnitCode2Name[code]
			workLoadDescList = append(workLoadDescList, fmt.Sprintf("%s:%v", name, value))
		}
	}

	if len(workLoadDescList) > 0 {
		current.WorkLoadDesc = strings.Join(workLoadDescList, ",")
	}
}

func (service *EfficiencyComputeService) createAndComputeCtx(employeeNumber string, operateDay time.Time) (*domain.ComputeContext, *calc_dynamic_param.CalcParam, *business_error.BusinessError) {
	employeeSnapshotService := DomainHolder.EmployeeSnapshotService
	calcDynamicParamService := DomainHolder.CalcDynamicParamService
	processService := DomainHolder.ProcessService
	workplaceService := DomainHolder.WorkplaceService

	// 1.获取员工当天快照和工作点信息
	employee := employeeSnapshotService.FindEmployeeSnapshot(employeeNumber, operateDay)
	workplace := workplaceService.FindByCode(employee.WorkplaceCode)

	// 2.初始化计算参数
	// 包括动态维度，计算聚合结果，工序加工节点列表，工序映射服务
	calcParam, err := calcDynamicParamService.FindParamsByWorkplace(employee.WorkplaceCode)
	if err != nil {
		return nil, nil, err
	}

	// 3. 查询工序映射关系
	processPositionList := processService.FindProcessPositionByWorkplace(employee.WorkplaceCode)

	ctx := domain.ComputeContext{
		Employee:       employee,
		Workplace:      workplace,
		OperateDay:     operateDay,
		ProcessList:    processPositionList,
		CalcOtherParam: calcParam.CalcOtherParam,
	}

	// 4. 注入原始数据
	if _, err := injectActions(&ctx, *calcParam); err != nil {
		return nil, nil, business_error.InjectDataError(err)
	}

	// 4.循环执行所有节点
	ctx.CalcStartTime = time.Now()
	ctxPointer := &ctx
	for _, node := range calcParam.CalcNodeList.List {
		if f, exists := calc_node.GetFunction(node.NodeName); exists {
			ctxRes := f(ctxPointer)
			ctxPointer = ctxRes
		} else {
			return nil, nil, business_error.NoCalcNodeError()
		}
	}
	ctx.CalcEndTime = time.Now()

	return ctxPointer, calcParam, nil
}

// 处理聚合存储逻辑
func (service *EfficiencyComputeService) handleSummary(ctx *domain.ComputeContext, storageParam calc_dynamic_param.SinkStorage, otherParam calc_dynamic_param.CalcOtherParam) {
	summarySinkService := Holder.SummarySinkService
	hourSummaryResultList := service.efficientAggregateActions(ctx.TodayWorkList, storageParam, otherParam.HourSummary, otherParam.Work)
	summarySinkService.BatchInsertSummaryResult(hourSummaryResultList, ctx.Employee, ctx.Workplace, ctx.OperateDay)
}

// 高效聚合算法
func (service *EfficiencyComputeService) efficientAggregateActions(works []domain.Actionable,
	storageParam calc_dynamic_param.SinkStorage,
	summaryParam calc_dynamic_param.HourSummaryParam,
	workParam calc_dynamic_param.WorkParam) []*domain.HourSummaryResult {
	// 1. 对Action按开始时间排序
	sort.Slice(works, func(i, j int) bool {
		return works[i].GetAction().ComputedStartTime.Before(*(works[j].GetAction().ComputedStartTime))
	})

	// 2. 创建桶收集聚合结果
	buckets := make(map[domain.HourSummaryAggregateKey]*domain.HourSummaryResult)

	// 3. 处理每个Action
	for _, work := range works {
		start := *(work.GetAction().ComputedStartTime)
		end := *(work.GetAction().ComputedEndTime)

		// 处理开始和结束时间相等的情况
		if start.Equal(end) {
			hourStart := start.Truncate(time.Hour)
			bucket := service.getOrCreateBucket(buckets, hourStart, work, storageParam.AggregateFields, storageParam.FieldName2ColumnName)
			if directWork, ok := work.(*domain.DirectWork); ok {
				bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, 1)
			}
			// 这里不处理工时
			continue
		}

		// 跳过无效时间段
		if start.After(end) {
			continue
		}

		// 计算总持续时间
		totalDuration := end.Sub(start).Seconds()

		// 确定开始和结束的小时桶
		startHour := start.Truncate(time.Hour)
		var endHour time.Time
		// 如果作业的结束时间正好等于整小时，按上一个小时计算
		if end == end.Truncate(time.Hour) {
			endHour = end.Truncate(time.Hour).Add(-time.Hour)
		} else {
			endHour = end.Truncate(time.Hour)
		}

		// 处理开始小时（可能不是完整小时）
		if start.Before(startHour.Add(time.Hour)) {
			segmentEnd := startHour.Add(time.Hour)
			if segmentEnd.After(end) {
				segmentEnd = end
			}
			duration := segmentEnd.Sub(start).Seconds()
			bucket := service.getOrCreateBucket(buckets, startHour, work, storageParam.AggregateFields, storageParam.FieldName2ColumnName)
			bucket.MergeTime(work, duration)

			// 根据策略处理物品数量
			switch summaryParam.WorkLoadAggregateType {
			case calc_dynamic_param.AggregateEndHour:
				if segmentEnd.Equal(end) {
					if directWork, ok := work.(*domain.DirectWork); ok {
						bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, 1)
					}
				}
			case calc_dynamic_param.AggregateProportion:
				if totalDuration > 0 {
					proportion := duration / totalDuration
					if directWork, ok := work.(*domain.DirectWork); ok {
						bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
					}
				} else {
					if directWork, ok := work.(*domain.DirectWork); ok {
						bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, 1)
					}
				}
			}
		}

		// 处理完整的小时
		if endHour.After(startHour.Add(time.Hour)) {
			currentHour := startHour.Add(time.Hour)
			for currentHour.Before(endHour) {
				duration := 3600.0
				bucket := service.getOrCreateBucket(buckets, currentHour, work, storageParam.AggregateFields, storageParam.FieldName2ColumnName)
				bucket.MergeTime(work, duration)

				// 根据策略处理物品数量
				switch summaryParam.WorkLoadAggregateType {
				case calc_dynamic_param.AggregateEndHour:
					// 完整小时不累加物品数量（只累加到结束小时）
				case calc_dynamic_param.AggregateProportion:
					if totalDuration > 0 {
						proportion := duration / totalDuration
						if directWork, ok := work.(*domain.DirectWork); ok {
							bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
						}
					}
				}

				currentHour = currentHour.Add(time.Hour)
			}

		}

		// 处理结束小时（可能不是完整小时）
		if end.After(endHour) && !startHour.Equal(endHour) {
			duration := end.Sub(endHour).Seconds()
			bucket := service.getOrCreateBucket(buckets, endHour, work, storageParam.AggregateFields, storageParam.FieldName2ColumnName)
			bucket.MergeTime(work, duration)

			// 根据策略处理物品数量
			switch summaryParam.WorkLoadAggregateType {
			case calc_dynamic_param.AggregateEndHour:
				if directWork, ok := work.(*domain.DirectWork); ok {
					bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, 1)
				}

			case calc_dynamic_param.AggregateProportion:
				if totalDuration > 0 {
					proportion := duration / totalDuration
					if directWork, ok := work.(*domain.DirectWork); ok {
						bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
					}
				} else {
					if directWork, ok := work.(*domain.DirectWork); ok {
						bucket.MergeWorkLoad(directWork.GetWorkLoad(), workParam.WorkLoadUnits, 1)
					}
				}
			}
		}
	}

	// 4. 将桶转换为有序切片
	result := make([]*domain.HourSummaryResult, 0, len(buckets))
	for _, bucket := range buckets {
		result = append(result, bucket)
	}

	// 5. 按小时排序结果
	sort.Slice(result, func(i, j int) bool {
		return result[i].AggregateKey.OperateTime.Before(result[j].AggregateKey.OperateTime)
	})

	return result
}

// 辅助函数：获取或创建桶
func (service *EfficiencyComputeService) getOrCreateBucket(buckets map[domain.HourSummaryAggregateKey]*domain.HourSummaryResult,
	operateTime time.Time,
	work domain.Actionable,
	aggregateFields []string, field2Column map[string]string) *domain.HourSummaryResult {
	key := service.buildAggregateKey(operateTime, work, aggregateFields)

	if bucket, exists := buckets[key]; exists {
		return bucket
	}
	bucket := domain.MakeHourSummaryResult(key, work, field2Column)
	buckets[key] = &bucket
	return &bucket
}

func (service *EfficiencyComputeService) buildAggregateKey(operateTime time.Time, work domain.Actionable, aggregateFields []string) domain.HourSummaryAggregateKey {
	key := domain.HourSummaryAggregateKey{
		EmployeeNumber: work.GetAction().EmployeeNumber,
		WorkplaceCode:  work.GetAction().WorkplaceCode,
		ProcessCode:    work.GetAction().ProcessCode,
		OperateTime:    operateTime,
	}

	var str = ""
	for _, field := range aggregateFields {
		value := work.GetAction().Properties[field]
		str += fmt.Sprint(value) + "|"
	}

	key.PropertyValues = str

	return key
}

func injectActions(ctx *domain.ComputeContext, param calc_dynamic_param.CalcParam) (bool, error) {
	actionService := DomainHolder.ActionService

	yesterday := ctx.OperateDay.AddDate(0, 0, -1)
	tomorrow := ctx.OperateDay.AddDate(0, 0, 1)
	operateDayRange := []time.Time{yesterday, ctx.OperateDay, tomorrow}

	if day2ActionData, err := actionService.FindEmployeeActions(ctx.Employee.Number, operateDayRange, param.InjectSource); err != nil {
		return false, err
	} else {
		ctx.YesterdayWorkList = day2ActionData.Day2WorkList[yesterday]
		ctx.TodayWorkList = day2ActionData.Day2WorkList[ctx.OperateDay]
		ctx.TomorrowWorkList = day2ActionData.Day2WorkList[tomorrow]

		ctx.YesterdayAttendance = day2ActionData.Day2Attendance[yesterday]
		ctx.TodayAttendance = day2ActionData.Day2Attendance[ctx.OperateDay]
		ctx.TomorrowAttendance = day2ActionData.Day2Attendance[tomorrow]

		ctx.TodayScheduling = day2ActionData.Day2Scheduling[ctx.OperateDay]
		ctx.TodayRestList = day2ActionData.Day2RestList[ctx.OperateDay]

		return true, nil
	}
}

// 根据工作点获取人效计算参数
// Parameters: 工作点编码
// Returns: 人效计算参数
func (service *EfficiencyComputeService) initComputeParams(workplaceCode string) ComputeParams {
	calcDynamicParamService := DomainHolder.CalcDynamicParamService

	calcDynamicParamService.FindParamsByWorkplace(workplaceCode)

	computeParams := ComputeParams{}

	return computeParams
}

func (service *EfficiencyComputeService) initComputeContext() domain.ComputeContext {
	ctx := domain.ComputeContext{}

	return ctx
}

func (service *EfficiencyComputeService) handleEmployeeStatus(ctx *domain.ComputeContext) {
	employeeStatusSinkService := Holder.EmployeeStatusSinkService
	employeeStatusEntity := service.filterEmployeeStatus(ctx)
	employeeStatusSinkService.InsertOrUpdateEmployeeStatus(employeeStatusEntity)
}

func (service *EfficiencyComputeService) filterEmployeeStatus(ctx *domain.ComputeContext) *entity.EmployeeStatusEntity {
	if ctx.TodayAttendance != nil && ctx.TodayAttendance.EndTime != nil {
		attendanceStatus := entity.EmployeeStatusEntity{
			EmployeeNumber: ctx.Employee.Number,
			EmployeeName:   ctx.Employee.Name,
			OperateDay:     &ctx.OperateDay,
			WorkplaceCode:  ctx.Workplace.Code,
			Status:         entity.OffDuty,
			LastActionTime: ctx.TodayAttendance.EndTime,
			LastActionCode: ctx.TodayAttendance.ActionCode,
			WorkGroupCode:  ctx.Employee.WorkGroupCode,
		}

		if ctx.TodayWorkList != nil && len(ctx.TodayWorkList) > 0 {
			for _, actionable := range ctx.TodayWorkList {
				if actionable.GetAction().ComputedStartTime.After(*ctx.TodayAttendance.EndTime) {
					log.ComputeLogger.Warn(fmt.Sprintf("下班后仍然有其他动作:%v", actionable.GetAction().ActionCode))
					break
				}
			}
		}

		return &attendanceStatus
	} else {
		if ctx.TodayWorkList != nil && len(ctx.TodayWorkList) > 0 {
			var status *entity.EmployeeStatusEntity

			now := ctx.CalcStartTime
			// 从后往前遍历
			for i := len(ctx.TodayWorkList) - 1; i >= 0; i-- {
				action := ctx.TodayWorkList[i]

				// 动作的结束时间为当前时间，说明动作正在进行中
				if action.GetAction().EndTime == nil && action.GetAction().ComputedEndTime.Equal(now) {
					status = &entity.EmployeeStatusEntity{
						EmployeeNumber: ctx.Employee.Number,
						EmployeeName:   ctx.Employee.Name,
						OperateDay:     &ctx.OperateDay,
						WorkplaceCode:  ctx.Workplace.Code,
						WorkGroupCode:  ctx.Employee.WorkGroupCode,
					}

					if action.GetAction().ActionType == domain.IDLE {
						status.Status = entity.Idle
						if i == 0 && ctx.TodayAttendance != nil {
							status.LastActionTime = ctx.TodayAttendance.StartTime
							status.LastActionCode = ctx.TodayAttendance.ActionCode
						} else if i > 0 {
							status.LastActionTime = ctx.TodayWorkList[i-1].GetAction().ComputedEndTime
							status.LastActionCode = ctx.TodayWorkList[i-1].GetAction().ActionCode
						}
					} else {
						status.LastActionTime = action.GetAction().StartTime
						status.LastActionCode = action.GetAction().ActionCode
						switch action.GetAction().ActionType {
						case domain.REST:
							status.Status = entity.Rest
						case domain.DIRECT_WORK:
							status.Status = entity.DirectWorking
						case domain.INDIRECT_WORK:
							status.Status = entity.IndirectWorking
						}
					}

					break
				}
			}

			// 如果没有正在进行中的作业，且没有考勤，认为OffDutyWithoutEndTime

			if status == nil {
				lastAction := ctx.TodayWorkList[len(ctx.TodayWorkList)-1]

				status = &entity.EmployeeStatusEntity{
					EmployeeNumber: ctx.Employee.Number,
					EmployeeName:   ctx.Employee.Name,
					OperateDay:     &ctx.OperateDay,
					WorkplaceCode:  ctx.Workplace.Code,
					WorkGroupCode:  ctx.Employee.WorkGroupCode,
					Status:         entity.OffDutyWithoutEndTime,
					LastActionTime: lastAction.GetAction().ComputedEndTime,
					LastActionCode: lastAction.GetAction().ActionCode,
				}

				return status
			}

		}
	}

	return nil
}

type ComputeResult struct {
	Success        bool
	EmployeeNumber string
	Ctx            *domain.ComputeContext
	Err            *business_error.BusinessError
}

func (service *EfficiencyComputeService) ComputeWorkplace(workplaceCode string, operateDay time.Time) (bool, *business_error.BusinessError) {
	workplaceService := DomainHolder.WorkplaceService
	employeeSnapshotService := DomainHolder.EmployeeSnapshotService
	calcDynamicParamService := DomainHolder.CalcDynamicParamService
	processPositionService := DomainHolder.ProcessService

	workplace := workplaceService.FindByCode(workplaceCode)
	if workplace == nil {
		return false, business_error.WorkplaceDoseNotExist(workplaceCode)
	}

	// 1.获取工作点下所有员工当天快照
	employeeSnapshots := employeeSnapshotService.FindWorkplaceEmployeeSnapshot(workplace, operateDay)

	// 2.初始化计算参数
	// 包括动态维度，计算聚合结果，工序加工节点列表，工序映射服务
	calcParam, err := calcDynamicParamService.FindParamsByWorkplace(workplaceCode)
	if err != nil {
		return false, err
	}

	// 3. 查询工序映射关系
	processPositionList := processPositionService.FindProcessPositionByWorkplace(workplaceCode)

	var wg sync.WaitGroup
	ctxChannel := make(chan *ComputeResult, len(employeeSnapshots))

	for _, employeeSnapshot := range employeeSnapshots {
		wg.Add(1)
		go runComputeEmployee(service, employeeSnapshot, workplace, operateDay, calcParam, processPositionList, &wg, ctxChannel)
	}

	// 等待协程完成
	wg.Wait()

	isAllSuccess := false
	successNum := 0
	for computeResult := range ctxChannel { // 从通道接收数据
		if computeResult.Success {
			// 处理聚合
			for _, storage := range calcParam.SinkStorages {
				switch storage.SinkType {
				case calc_dynamic_param.SUMMARY:
					service.handleSummary(computeResult.Ctx, storage, calcParam.CalcOtherParam)
				case calc_dynamic_param.EMPLOYEE_STATUS:
					service.handleEmployeeStatus(computeResult.Ctx)
				}
			}
		} else {
			log.LogBusinessError(computeResult.Err)
		}

		successNum++

		// 根据计算粒度分布式解锁
		if unlockSuccess, err := lock.Unlock(computeResult.EmployeeNumber, operateDay); err != nil {
			businessError := business_error.UnlockFailureBySystemError(err)
			log.LogBusinessError(businessError)
		} else if !unlockSuccess {
			log.ComputeLogger.Info("员工人效数据计算，解锁失败", "number", computeResult.EmployeeNumber)
		} else {
			log.ComputeLogger.Info("员工人效数据计算，解锁成功", "number", computeResult.EmployeeNumber)
		}

		if successNum == len(employeeSnapshots) {
			break
		}
	}

	if successNum == len(employeeSnapshots) {
		isAllSuccess = true
	}

	return isAllSuccess, nil
}

func (service *EfficiencyComputeService) computeEachEmployee(employeeSnapshot *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time,
	calcParam *calc_dynamic_param.CalcParam, standardPositionList []*domain.StandardPosition) (*domain.ComputeContext, *business_error.BusinessError) {
	// 根据计算粒度分布式加锁
	if lockSuccess, err := lock.Lock(employeeSnapshot.Number, operateDay, 2); err != nil {
		businessError := business_error.LockFailureBySystemError(err)
		log.LogBusinessError(businessError)
		return nil, businessError
	} else if !lockSuccess {
		log.ComputeLogger.Info("员工人效数据计算，加锁失败", "number", employeeSnapshot.Number)
		return nil, business_error.LockFailure()
	} else {
		log.ComputeLogger.Info("员工人效数据计算，加锁成功", "number", employeeSnapshot.Number)
	}

	ctx := domain.ComputeContext{
		Employee:       employeeSnapshot,
		Workplace:      workplace,
		OperateDay:     operateDay,
		ProcessList:    standardPositionList,
		CalcOtherParam: calcParam.CalcOtherParam,
	}

	// 4. 注入原始数据
	if _, err := injectActions(&ctx, *calcParam); err != nil {
		return nil, business_error.InjectDataError(err)
	}

	// 4.循环执行所有节点
	ctx.CalcStartTime = time.Now()
	ctxPointer := &ctx
	for _, node := range calcParam.CalcNodeList.List {
		if f, exists := calc_node.GetFunction(node.NodeName); exists {
			ctxRes := f(ctxPointer)
			ctxPointer = ctxRes
		} else {
			return nil, business_error.NoCalcNodeError()
		}
	}
	ctx.CalcEndTime = time.Now()
	log.ComputeLogger.Info("员工人效数据计算完成", "number", employeeSnapshot.Number, "耗时(ms)", ctx.CalcEndTime.Sub(ctx.CalcStartTime).Milliseconds())

	return &ctx, nil
}

func runComputeEmployee(service *EfficiencyComputeService, employee *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time,
	calcParam *calc_dynamic_param.CalcParam, standardPositionList []*domain.StandardPosition, wg *sync.WaitGroup, channel chan *ComputeResult) {
	// 协程结束时通知WaitGroup
	defer wg.Done()
	if eachCtx, err := service.computeEachEmployee(employee, workplace, operateDay, calcParam, standardPositionList); err == nil {
		channel <- &ComputeResult{true, employee.Number, eachCtx, nil}
	} else {
		log.LogBusinessError(err)
		channel <- &ComputeResult{false, employee.Number, nil, err}
	}
}
