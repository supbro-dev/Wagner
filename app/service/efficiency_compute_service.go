package service

import (
	"fmt"
	"time"
	"wagner/app/utils/lock"
)

// 人效计算服务
type EfficiencyComputeService struct {
}

// 人效计算上下文
type ComputeContext struct {
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
	// 1.获取员工当天快照
	employee := employeeSnapshotService.FindEmployeeSnapshot(employeeNumber, operateDay)

	// 根据工作点查找工序实施编码
	//holder.ServiceHolder.StandardPositionService
	fmt.Println(actionService, employee)

	// 3.根据计算粒度分布式加锁

	b, err := lock.Lock(employeeNumber)
	fmt.Println(b, err)

	// 初始化计算参数
	// 包括动态维度，计算聚合结果，工序加工节点列表，工序映射服务
}

// 根据工作点获取人效计算参数
// Parameters: 工作点编码
// Returns: 人效计算参数
func (service *EfficiencyComputeService) initComputeParams(workplaceCode string) ComputeParams {
	computeParams := ComputeParams{}

	return computeParams
}

func (service *EfficiencyComputeService) initComputeContext() ComputeContext {
	ctx := ComputeContext{}

	return ctx
}
