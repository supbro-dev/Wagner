package employee_snapshot

import (
	"github.com/jinzhu/copier"
	"time"
	"wagner/app/domain"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/query"
)

type EmployeeSnapshotService struct {
	employeeDao *dao.EmployeeDao
}

// 通过构造函数注入 DAO
func CreateEmployeeSnapshotService(employeeDao *dao.EmployeeDao) *EmployeeSnapshotService {
	return &EmployeeSnapshotService{employeeDao: employeeDao}
}

// 根据员工工号和日期查找这天的员工快照列表
func (service *EmployeeSnapshotService) FindEmployeeSnapshot(employeeNumber string, operateDay time.Time) *domain.EmployeeSnapshot {
	// 生产环境需要根据员工一段时间的履历，获取在某个工作点某天的人员快照，这里简单使用人员信息代替
	employee := service.employeeDao.FindByNumber(employeeNumber)
	return convertEmployee(employee)
}

// 根据工作点和日期查找这天在工作点工作的员工快照列表
func (service *EmployeeSnapshotService) FindWorkplaceEmployeeSnapshot(workplace *domain.Workplace, operateDay time.Time) []*domain.EmployeeSnapshot {
	// 生产环境需要根据员工一段时间的履历，获取在某个工作点某天的人员快照，这里简单使用人员信息代替
	employeeEntities := service.employeeDao.FindByWorkplaceCode(workplace.Code)

	employeeSnapshotList := make([]*domain.EmployeeSnapshot, 0)
	for _, employeeEntity := range employeeEntities {
		employeeSnapshot := convertEmployee(employeeEntity)
		employeeSnapshotList = append(employeeSnapshotList, employeeSnapshot)
	}

	return employeeSnapshotList
}

func (service *EmployeeSnapshotService) FindByInfo(name string, workGroupCode string, workplaceCode string) *domain.EmployeeSnapshot {
	q := query.EmployeeQuery{
		name, workGroupCode, workplaceCode,
	}
	employeeEntityList := service.employeeDao.FindByQuery(&q)

	if len(employeeEntityList) > 0 {
		return convertEmployee(employeeEntityList[0])
	} else {
		return nil
	}
}

func (service *EmployeeSnapshotService) FindByWorkGroupCode(workGroupCode string, workplaceCode string) []*domain.EmployeeSnapshot {
	q := query.EmployeeQuery{
		WorkGroupCode: workGroupCode,
		WorkplaceCode: workplaceCode,
	}
	employeeEntityList := service.employeeDao.FindByQuery(&q)

	domainList := make([]*domain.EmployeeSnapshot, 0)
	for _, employeeEntity := range employeeEntityList {
		domainList = append(domainList, convertEmployee(employeeEntity))
	}

	return domainList
}

func convertEmployee(employee *entity.EmployeeEntity) *domain.EmployeeSnapshot {
	employeeSnapshot := domain.EmployeeSnapshot{}
	copier.Copy(&employeeSnapshot, employee)

	properties := make(map[string]string)
	properties["employeeName"] = employee.Name
	employeeSnapshot.Properties = properties

	return &employeeSnapshot
}
