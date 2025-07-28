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

func (dao *EmployeeDao) FindByQuery(q *query.EmployeeQuery) []*entity.EmployeeEntity {
	list := make([]*entity.EmployeeEntity, 0)
	tx := dao.db.Model(entity.EmployeeEntity{}).Where("work_group_code = ? and workplace_code = ?", q.WorkGroupCode, q.WorkplaceCode)
	if q.Name != "" {
		tx = tx.Where("name = ?", q.Name)
	}

	tx.Find(&list)
	return list
}
