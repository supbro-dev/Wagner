package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
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
