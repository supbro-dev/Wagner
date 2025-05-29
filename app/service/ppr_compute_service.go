package service

import (
	"fmt"
	"time"
)

type PprComputeService struct {
}

func CreatePprComputeService() *PprComputeService {
	return &PprComputeService{}
}

func (service *PprComputeService) ComputeEmployee(employeeNumber string, operateDay time.Time) {
	employeeSnapshotService := DomainHolder.EmployeeSnapshotService
	actionService := DomainHolder.ActionService
	// 1.获取员工当天快照
	employee := employeeSnapshotService.FindEmployeeSnapshot(employeeNumber, operateDay)

	// 根据工作点查找工序实施编码
	//holder.ServiceHolder.StandardPositionService
	fmt.Println(actionService, employee)
}
