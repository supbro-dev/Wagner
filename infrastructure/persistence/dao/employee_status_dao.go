/*
* @Author: supbro
* @Date:   2025/6/18 14:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/18 14:06
 */
package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"wagner/infrastructure/persistence/entity"
)

type EmployeeStatusDao struct {
	db *gorm.DB
}

func (dao EmployeeStatusDao) InsertOrUpdate(entity *entity.EmployeeStatusEntity) {
	dao.db.Omit("gmt_create", "gmt_modified").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "employee_number"}, {Name: "operate_day"}, {Name: "workplace_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "last_action_time", "last_action_code"}), // 更新字段
	}).Create(entity)
}

func (dao EmployeeStatusDao) FindByWorkplaceAndDate(workplaceCode string, operateDay time.Time) []*entity.EmployeeStatusEntity {
	result := make([]*entity.EmployeeStatusEntity, 0)
	dao.db.
		Where("workplace_code = ? and operate_day = ?", workplaceCode, operateDay).Find(&result)
	return result
}

func CreateEmployeeStatusDao(client *gorm.DB) *EmployeeStatusDao {
	return &EmployeeStatusDao{client}
}
