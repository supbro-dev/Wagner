package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type ProcessPositionDao struct {
	db *gorm.DB
}

func CreateProcessPositionDao(client *gorm.DB) *ProcessPositionDao {
	return &ProcessPositionDao{client}
}

func (dao *ProcessPositionDao) FindByCode(code string) entity.ProcessPositionEntity {
	standardPosition := entity.ProcessPositionEntity{}
	dao.db.Where("code = ?", code).Find(&standardPosition)
	return standardPosition
}

func (dao *ProcessPositionDao) FindByIndustry(industryCode string, subIndustryCode string, version int64) []*entity.ProcessPositionEntity {
	array := make([]*entity.ProcessPositionEntity, 0)
	if subIndustryCode != "" {
		dao.db.Where("industry_code = ? and sub_industry_code = ? and version = ?", industryCode, subIndustryCode, version).
			Order("level").
			Find(&array)
		if len(array) > 0 {
			return array
		}
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).
			Order("level").
			Find(&array)
		return array
	} else {
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).
			Order("level").
			Find(&array)
		return array
	}
}

// 根据父编码和版本号查找
func (dao *ProcessPositionDao) FindByParentCodeAndVersion(parentCode string, version int64) []*entity.ProcessPositionEntity {
	array := make([]*entity.ProcessPositionEntity, 0)
	dao.db.Where("parent_code = ? and version = ?", parentCode, version).
		Order("sort_index").
		Find(&array)

	return array
}
