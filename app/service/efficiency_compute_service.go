package service

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/utils/lock"
	"wagner/app/utils/script_util"
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

func (service *EfficiencyComputeService) ComputeEmployee(employeeNumber string, operateDay time.Time) {
	employeeSnapshotService := DomainHolder.EmployeeSnapshotService
	actionService := DomainHolder.ActionService
	calcDynamicParamService := DomainHolder.CalcDynamicParamService
	// 1.获取员工当天快照
	employee := employeeSnapshotService.FindEmployeeSnapshot(employeeNumber, operateDay)

	// 2.初始化计算参数
	// 包括动态维度，计算聚合结果，工序加工节点列表，工序映射服务
	calcParam := calcDynamicParamService.FindParamsByWorkplace(employee.WorkplaceCode)

	actions := actionService.FindEmployeeActions(employeeNumber, []time.Time{operateDay})

	script := ` fmt := import("fmt")
    fmt.println("myCtx:", ctx)
ctx.Name = "123"
ctxResult := ctx
`
	ctx := domain.ComputeContext{
		Employee:        employee,
		TodayActionList: actions,
		OperateDay:      operateDay,
	}
	ret, err2 := script_util.Run[*domain.ComputeContext, *domain.ComputeContext](script, entity.GOLANG, &ctx, "ctx")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println(ret)

	fmt.Println(actionService, calcParam)

	// 3.根据工作点查找工序实施编码
	//holder.ServiceHolder.StandardPositionService

	// 3.根据计算粒度分布式加锁

	b, err := lock.Lock(employeeNumber)
	fmt.Println(b, err)

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
