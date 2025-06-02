/*
* @Author: supbro
* @Date:   2025/6/2 11:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 11:29
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type WorkplaceDao struct {
	db *gorm.DB
}

func (dao *WorkplaceDao) FindByCode(code string) entity.WorkplaceEntity {
	workplace := entity.WorkplaceEntity{}
	dao.db.Where("code = ?", code).Find(&workplace)
	return workplace
}

func CreateWorkplaceDao(client *gorm.DB) *WorkplaceDao {
	return &WorkplaceDao{client}
}
