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

// 根据行业/子行业查询计算参数（如果查到了子行业，优先选取子行业的计算参数，如果查不到，选用主行业的计算参数）
// Parameters: 行业、子行业
// Returns: 参数列表
func (dao CalcDynamicParamDao) FindByIndustry(industryCode string, subIndustryCode string) []entity.CalcDynamicParamEntity {
	paramList := make([]entity.CalcDynamicParamEntity, 0)
	if subIndustryCode != "" {
		dao.db.Where("industry_code = ? and sub_industry_code = ?", industryCode, subIndustryCode).Find(&paramList)
		if len(paramList) > 0 {
			return paramList
		}
	}

	dao.db.Where("industry_code = ?", industryCode).Find(&paramList)
	return paramList
}

func CreateCalcDynamicParamDao(client *gorm.DB) *CalcDynamicParamDao {
	return &CalcDynamicParamDao{client}
}
