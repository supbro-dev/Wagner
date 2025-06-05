package dao

import (
	"gorm.io/gorm"
	"time"
	"wagner/infrastructure/persistence/common"
	"wagner/infrastructure/persistence/entity"
)

type ActionDao struct {
	common.BaseDao
	db *gorm.DB
}

func CreateActionRepository(client *gorm.DB) *ActionDao {
	return &ActionDao{db: client}
}

func (dao *ActionDao) FindBy(employeeNumber string, operateDayList []time.Time) []entity.ActionEntity {
	var actions []entity.ActionEntity
	dao.db.Where("employee_number = ? and operate_day in ?", employeeNumber, dao.TimeList2DateList(operateDayList)).Order("operate_day asc").Find(&actions)
	return actions
}
