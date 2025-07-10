/*
* @Author: supbro
* @Date:   2025/7/9 11:24
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/9 11:24
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/common"
	"wagner/infrastructure/persistence/entity"
)

type PositionDao struct {
	common.BaseDao
	db *gorm.DB
}

func CreatePositionDao(client *gorm.DB) *PositionDao {
	return &PositionDao{db: client}
}

func (dao *PositionDao) FindAll() []*entity.PositionEntity {
	var positions []*entity.PositionEntity
	dao.db.Find(&positions)
	return positions
}
