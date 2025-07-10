/*
* @Author: supbro
* @Date:   2025/6/2 10:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:47
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type CalcDynamicParamDao struct {
	db *gorm.DB
}

func (dao CalcDynamicParamDao) FindByMode(industryCode string, subIndustryCode string, workplaceCode string, mode entity.ParamMode) []*entity.CalcDynamicParamEntity {
	paramList := make([]*entity.CalcDynamicParamEntity, 0)
	switch mode {
	case entity.IndustryMode:
		dao.db.Where("industry_code = ? and mode = ?", industryCode, mode).Find(&paramList)
	case entity.SubIndustryMode:
		dao.db.Where("industry_code = ? and sub_industry_code = ? and mode = ?", industryCode, subIndustryCode, mode).Find(&paramList)
	case entity.WorkplaceMode:
		dao.db.Where("workplace_code = ? and mode = ? ", workplaceCode, mode).Find(&paramList)
	}
	return paramList
}

func (dao CalcDynamicParamDao) UpdateContentById(content string, id int64) {
	dao.db.Model(&entity.CalcDynamicParamEntity{}).Where("id = ?", id).UpdateColumn("content", content)
}

func CreateCalcDynamicParamDao(client *gorm.DB) *CalcDynamicParamDao {
	return &CalcDynamicParamDao{client}
}
