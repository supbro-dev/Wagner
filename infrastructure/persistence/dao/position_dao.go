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

func (dao *PositionDao) FindAll(industryCode, subIndustryCode string) []*entity.PositionEntity {
	var positions []*entity.PositionEntity
	dao.db.Model(entity.PositionEntity{}).Where("industry_code = ? and sub_industry_code = ?", industryCode, subIndustryCode).Find(&positions)
	return positions
}

func (dao *PositionDao) FindByCodeAndIndustry(positionCode string, industryCode string, subIndustryCode string) *entity.PositionEntity {
	var positions []*entity.PositionEntity
	dao.db.Model(entity.PositionEntity{}).Where("code = ? and industry_code = ? and sub_industry_code = ?", positionCode, industryCode, subIndustryCode).Find(&positions)

	if len(positions) > 0 {
		return positions[0]
	} else {
		return nil
	}
}
