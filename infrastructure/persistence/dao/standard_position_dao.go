package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type StandardPositionDao struct {
	db *gorm.DB
}

func CreateStandardPositionDao(client *gorm.DB) *StandardPositionDao {
	return &StandardPositionDao{client}
}

func (dao *StandardPositionDao) FindByCode(code string) entity.StandardPositionEntity {
	standardPosition := entity.StandardPositionEntity{}
	dao.db.Where("code = ?", code).Find(&standardPosition)
	return standardPosition
}

func (dao *StandardPositionDao) FindByIndustry(industryCode string, subIndustryCode string, version int) []*entity.StandardPositionEntity {
	array := make([]*entity.StandardPositionEntity, 0)
	if subIndustryCode != "" {
		dao.db.Where("industry_code = ? and sub_industry_code = ? and version = ?", industryCode, subIndustryCode, version).Find(&array)
		if len(array) > 0 {
			return array
		}
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).Find(&array)
		return array
	} else {
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).Find(&array)
		return array
	}
}

// 根据行业查找最大版本（在没有实施功能时，临时获取最大版本方法）
// Parameters: industryCode 行业, subIndustryCode 子行业
// Returns: version的最大值
func (dao *StandardPositionDao) FindMaxVersionByIndustry(industryCode string, subIndustryCode string) int {
	var maxVersion int
	if subIndustryCode != "" {

		dao.db.Model(&entity.StandardPositionEntity{}).
			Where("industry_code = ? and sub_industry_code = ?", industryCode, subIndustryCode).
			Select("max(version)").
			First(&maxVersion)
		if maxVersion > 0 {
			return maxVersion
		}
		dao.db.Model(&entity.StandardPositionEntity{}).
			Where("industry_code = ? ", industryCode).
			Select("max(version)").
			First(&maxVersion)
		return maxVersion
	} else {
		dao.db.Model(&entity.StandardPositionEntity{}).
			Where("industry_code = ? ", industryCode).
			Select("max(version)").
			First(&maxVersion)
		return maxVersion
	}
}

// 根据版本查找所有StandardPosition
func (dao *StandardPositionDao) FindByVersion(version int) []*entity.StandardPositionEntity {
	array := make([]*entity.StandardPositionEntity, 0)
	dao.db.Where("version = ?", version).Find(&array)
	return array
}
