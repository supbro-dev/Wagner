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

func (dao *ProcessPositionDao) FindByCode(code string, version int64) *entity.ProcessPositionEntity {
	array := make([]*entity.ProcessPositionEntity, 0)
	dao.db.Where("code = ? and version = ?", code, version).Find(&array)

	if len(array) > 0 {
		return array[0]
	} else {
		return nil
	}
}

func (dao *ProcessPositionDao) FindByIndustry(industryCode string, subIndustryCode string, version int64) []*entity.ProcessPositionEntity {
	array := make([]*entity.ProcessPositionEntity, 0)
	if subIndustryCode != "" {
		dao.db.Where("industry_code = ? and sub_industry_code = ? and version = ?", industryCode, subIndustryCode, version).
			Order("level, sort_index").
			Find(&array)
		if len(array) > 0 {
			return array
		}
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).
			Order("level, sort_index").
			Find(&array)
		return array
	} else {
		dao.db.Where("industry_code = ? and version = ?", industryCode, version).
			Order("level, sort_index").
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

func (dao *ProcessPositionDao) Insert(e *entity.ProcessPositionEntity) {
	dao.db.Omit("gmt_create", "gmt_modified").Create(e)
}

func (dao *ProcessPositionDao) Update(e *entity.ProcessPositionEntity) {
	dao.db.Omit("gmt_create", "gmt_modified").Model(entity.ProcessPositionEntity{}).Where("id = ?", e.Id).Updates(e)
}

func (dao *ProcessPositionDao) DeleteById(id int64) {
	dao.db.Model(entity.ProcessPositionEntity{}).Where("id = ?", id).Delete(id)
}
