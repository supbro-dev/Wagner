package service

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
	"wagner/app/domain"
	"wagner/app/global/my_const"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/utils/lock"
	"wagner/app/utils/script_util"
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

func (service *EfficiencyComputeService) ComputeEmployee(employeeNumber string, operateDay time.Time) {
	employeeSnapshotService := DomainHolder.EmployeeSnapshotService
	calcDynamicParamService := DomainHolder.CalcDynamicParamService
	standardPositionService := DomainHolder.StandardPositionService
	// 1.获取员工当天快照
	employee := employeeSnapshotService.FindEmployeeSnapshot(employeeNumber, operateDay)

	// 2.初始化计算参数
	// 包括动态维度，计算聚合结果，工序加工节点列表，工序映射服务
	calcParam := calcDynamicParamService.FindParamsByWorkplace(employee.WorkplaceCode)

	ctx := domain.ComputeContext{
		Employee:   employee,
		OperateDay: operateDay,
	}

	// 3. 查询工序映射关系
	standardPositionList := standardPositionService.FindStandardPositionByWorkplace(employee.WorkplaceCode)
	fmt.Println(standardPositionList)

	// 4. 注入原始数据
	injectActions(&ctx, calcParam)

	// 3.根据计算粒度分布式加锁
	lockSuccess, err := lock.Lock(employeeNumber)
	if err != nil {
		panic(err)
	}
	if !lockSuccess {
		return
	}

	// 4.循环执行所有节点
	ctxPointer := &ctx
	for _, node := range *calcParam.CalcNodeList.List {
		ctxRes, err := script_util.Run[*domain.ComputeContext, *domain.ComputeContext](node.NodeName, node.Script, node.NodeType, ctxPointer, "ctx")
		ctxPointer = ctxRes
		if err != nil {
			// todo 先继续执行
			//panic(err)
		}
	}

	// 5.处理聚合
	for _, storage := range *calcParam.SinkStorages {
		switch storage.SinkType {
		case my_const.SUMMARY:
			handleSummary(ctx, storage, calcParam.CalcOtherParam.HourSummary)
		}
	}

	unlockSuccess, err := lock.Unlock(employeeNumber)
	if err != nil || !unlockSuccess {
		panic(err)
	}
}

// 处理聚合存储逻辑
func handleSummary(ctx domain.ComputeContext, storage calc_dynamic_param.SinkStorage, summaryParam calc_dynamic_param.HourSummaryParam) {

}

// 高效聚合算法
func efficientAggregateActions(works []domain.Work, summaryParam *calc_dynamic_param.HourSummaryParam, aggregateFields []string, workParam *calc_dynamic_param.WorkParam) []domain.HourSummaryResult {
	// 1. 对Action按开始时间排序
	sort.Slice(works, func(i, j int) bool {
		return works[i].GetComputedStartTime().Before(works[j].GetComputedStartTime())
	})

	// 2. 创建桶收集聚合结果
	buckets := make(map[time.Time]*domain.HourSummaryResult)

	// 3. 处理每个Action
	for _, work := range works {
		start := work.GetComputedStartTime()
		end := work.GetComputedEndTime()

		// 处理开始和结束时间相等的情况
		if start.Equal(end) {
			hourStart := start.Truncate(time.Hour)
			bucket := getOrCreateBucket(buckets, hourStart)
			bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, 1)
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
		endHour := end.Truncate(time.Hour)

		// 处理开始小时（可能不是完整小时）
		if start.Before(startHour.Add(time.Hour)) {
			segmentEnd := startHour.Add(time.Hour)
			if segmentEnd.After(end) {
				segmentEnd = end
			}
			duration := segmentEnd.Sub(start).Seconds()
			bucket := getOrCreateBucket(buckets, startHour)
			bucket.TotalDuration += duration

			// 根据策略处理物品数量
			switch summaryParam.WorkLoadAggregate {
			case calc_dynamic_param.AggregateEndHour:
				if segmentEnd.Equal(end) {
					bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, 1)
				}
			case calc_dynamic_param.AggregateProportion:
				if totalDuration > 0 {
					proportion := duration / totalDuration
					bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
				} else {
					bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, 1)
				}
			}
		}

		// 处理完整的小时
		if endHour.After(startHour.Add(time.Hour)) {
			currentHour := startHour.Add(time.Hour)
			for currentHour.Before(endHour) {
				duration := 3600.0
				bucket := getOrCreateBucket(buckets, currentHour)
				bucket.TotalDuration += duration

				// 根据策略处理物品数量
				switch summaryParam.WorkLoadAggregate {
				case calc_dynamic_param.AggregateEndHour:
					// 完整小时不累加物品数量（只累加到结束小时）
				case calc_dynamic_param.AggregateProportion:
					if totalDuration > 0 {
						proportion := duration / totalDuration
						bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
					}
				}

				currentHour = currentHour.Add(time.Hour)
			}
		}

		// 处理结束小时（可能不是完整小时）
		if end.After(endHour) {
			duration := end.Sub(endHour).Seconds()
			bucket := getOrCreateBucket(buckets, endHour)
			bucket.TotalDuration += duration

			// 根据策略处理物品数量
			switch summaryParam.WorkLoadAggregate {
			case calc_dynamic_param.AggregateEndHour:
				bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, 1)
			case calc_dynamic_param.AggregateProportion:
				if totalDuration > 0 {
					proportion := duration / totalDuration
					bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, proportion)
				} else {
					bucket.TotalWorkLoad = bucket.MergeWorkLoad(work.GetWorkLoad(), workParam.WorkLoadUnits, 1)
				}
			}
		}
	}

	// 4. 将桶转换为有序切片
	result := make([]domain.HourSummaryResult, 0, len(buckets))
	for _, bucket := range buckets {
		result = append(result, *bucket)
	}

	// 5. 按小时排序结果
	sort.Slice(result, func(i, j int) bool {
		return result[i].OperateTime.Before(result[j].OperateTime)
	})

	return result
}

// 辅助函数：获取或创建桶
func (service *EfficiencyComputeService) getOrCreateBucket(buckets map[interface{}]*domain.HourSummaryResult, work *domain.Work, aggregateFields []string) *domain.HourSummaryResult {
	key := service.keyFromStruct(work, aggregateFields)

	if bucket, exists := buckets[key]; exists {
		return bucket
	}
	bucket := domain.MakeHourSummaryResult(key, aggregateFields)
	buckets[key] = &bucket
	return &bucket
}

func (service *EfficiencyComputeService) keyFromStruct(s interface{}, aggregateFields []string) interface{} {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 创建动态结构体类型
	fieldTypes := make([]reflect.StructField, len(aggregateFields))
	fieldValues := make([]reflect.Value, len(aggregateFields))

	for i, field := range aggregateFields {
		f := v.FieldByName(field)
		if !f.IsValid() {
			continue
		}

		fieldTypes[i] = reflect.StructField{
			Name: strings.Title(field),
			Type: f.Type(),
		}
		fieldValues[i] = f
	}

	// 创建动态结构体实例
	keyType := reflect.StructOf(fieldTypes)
	keyValue := reflect.New(keyType).Elem()

	for i, val := range fieldValues {
		if val.IsValid() {
			keyValue.Field(i).Set(val)
		}
	}

	return keyValue.Interface()
}

func injectActions(ctx *domain.ComputeContext, param *calc_dynamic_param.CalcParam) {
	actionService := DomainHolder.ActionService

	yesterday := ctx.OperateDay.AddDate(0, 0, -1)
	tomorrow := ctx.OperateDay.AddDate(0, 0, 1)
	operateDayRange := []time.Time{yesterday, ctx.OperateDay, tomorrow}

	day2WorkList, day2Attendance, day2Scheduling := actionService.FindEmployeeActions(ctx.Employee.Number, operateDayRange, param.InjectSource)

	ctx.YesterdayWorkList = day2WorkList[yesterday]
	ctx.TodayWorkList = day2WorkList[ctx.OperateDay]
	ctx.TodayWorkList = day2WorkList[tomorrow]

	ctx.YesterdayAttendance = day2Attendance[yesterday]
	ctx.TodayAttendance = day2Attendance[ctx.OperateDay]
	ctx.TomorrowAttendance = day2Attendance[tomorrow]

	ctx.TodayScheduling = day2Scheduling[ctx.OperateDay]
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
