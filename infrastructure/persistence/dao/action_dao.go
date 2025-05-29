package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type ActionDao struct {
	db *gorm.DB
}

func CreateActionRepository(client *gorm.DB) *ActionDao {
	return &ActionDao{client}
}

func (dao *ActionDao) FindBy(employeeNumber, operateDayList []string) []entity.ActionEntity {
	var actions []entity.ActionEntity
	dao.db.Where("employee_number = ? and operate_day in ?", employeeNumber, operateDayList).Order("operate_day asc").Find(&actions)
	return actions
}
