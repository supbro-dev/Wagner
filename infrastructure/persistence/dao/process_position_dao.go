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

func (dao *ProcessPositionDao) FindByIndustry(industryCode string, subIndustryCode string, version int) []*entity.ProcessPositionEntity {
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

// 根据行业查找最大版本（在没有实施功能时，临时获取最大版本方法）
// Parameters: industryCode 行业, subIndustryCode 子行业
// Returns: version的最大值
func (dao *ProcessPositionDao) FindMaxVersionByIndustry(industryCode string, subIndustryCode string) int {
	var maxVersion int
	if subIndustryCode != "" {

		dao.db.Model(&entity.ProcessPositionEntity{}).
			Where("industry_code = ? and sub_industry_code = ?", industryCode, subIndustryCode).
			Select("max(version)").
			First(&maxVersion)
		if maxVersion > 0 {
			return maxVersion
		}
		dao.db.Model(&entity.ProcessPositionEntity{}).
			Where("industry_code = ? ", industryCode).
			Select("max(version)").
			First(&maxVersion)
		return maxVersion
	} else {
		dao.db.Model(&entity.ProcessPositionEntity{}).
			Where("industry_code = ? ", industryCode).
			Select("max(version)").
			First(&maxVersion)
		return maxVersion
	}
}

// 根据版本查找所有StandardPosition
func (dao *ProcessPositionDao) FindByVersion(version int) []*entity.ProcessPositionEntity {
	array := make([]*entity.ProcessPositionEntity, 0)
	dao.db.Where("version = ?", version).Find(&array)
	return array
}
