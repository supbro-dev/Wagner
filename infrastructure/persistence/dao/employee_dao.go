package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/query"
)

type EmployeeDao struct {
	db *gorm.DB
}

func CreateEmployeeDao(client *gorm.DB) *EmployeeDao {
	return &EmployeeDao{client}
}

func (dao *EmployeeDao) FindByNumber(number string) *entity.EmployeeEntity {
	employee := entity.EmployeeEntity{}
	dao.db.Where("number = ?", number).Find(&employee)
	return &employee
}

func (dao *EmployeeDao) FindByWorkplaceCode(code string) []*entity.EmployeeEntity {
	employeeList := make([]*entity.EmployeeEntity, 0)
	dao.db.Where("workplace_code = ?", code).Find(&employeeList)

	return employeeList
}

func (dao *EmployeeDao) FindByQuery(q *query.EmployeeQuery) *entity.EmployeeEntity {
	list := make([]*entity.EmployeeEntity, 0)
	dao.db.Model(entity.EmployeeEntity{}).Where("name = ? and work_group_code = ? and workplace_code = ?", q.Name, q.WorkGroupCode, q.WorkplaceCode).Find(&list)
	if len(list) > 0 {
		return list[0]
	} else {
		return nil
	}
}
