package service

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/utils/lock"
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
	b, err := lock.Lock(employeeNumber)
	fmt.Println(b, err)

}

func injectActions(ctx *domain.ComputeContext, param *calc_dynamic_param.CalcParam) {
	actionService := DomainHolder.ActionService

	yesterday := ctx.OperateDay.AddDate(0, 0, -1)
	tomorrow := ctx.OperateDay.AddDate(0, 0, 1)
	operateDayRange := []time.Time{yesterday, ctx.OperateDay, tomorrow}

	day2WorkList, day2Attendance, day2Scheduling := actionService.FindEmployeeActions(ctx.Employee.Number, operateDayRange, param.OriginalField)

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
