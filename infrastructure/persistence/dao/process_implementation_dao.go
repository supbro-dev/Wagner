/*
* @Author: supbro
* @Date:   2025/7/1 13:03
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 13:03
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type ProcessImplementationDao struct {
	db *gorm.DB
}

func CreateProcessImplementationDao(client *gorm.DB) *ProcessImplementationDao {
	return &ProcessImplementationDao{client}
}

func (d *ProcessImplementationDao) FindByIndustry(industryCode string, subIndustryCode string) *entity.ProcessImplementationEntity {
	processImplementation := &entity.ProcessImplementationEntity{}
	if subIndustryCode != "" {
		d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and industry_code=? and sub_industry_code = ?", industryCode, subIndustryCode).First(processImplementation)
		if processImplementation == nil {
			d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and industry_code=?", industryCode).First(processImplementation)
		}
	} else {
		d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and industry_code=?", industryCode).First(processImplementation)
	}
	return processImplementation
}
