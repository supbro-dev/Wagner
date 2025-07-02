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
		d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'subIndustry' and target_code=?", subIndustryCode).First(processImplementation)
		if processImplementation != nil {
			return processImplementation
		}
	}
	d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'industry' and target_code=?", industryCode).First(processImplementation)
	return processImplementation
}

func (d *ProcessImplementationDao) FindByWorkplaceCode(workplaceCode string) *entity.ProcessImplementationEntity {
	processImplementation := &entity.ProcessImplementationEntity{}
	d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'workplace' and target_code=? ", workplaceCode).First(processImplementation)
	return processImplementation
}
