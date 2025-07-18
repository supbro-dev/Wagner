/*
* @Author: supbro
* @Date:   2025/7/17 14:30
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/17 14:30
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/common"
	"wagner/infrastructure/persistence/entity"
)

type WorkGroupDao struct {
	common.BaseDao
	db *gorm.DB
}

func CreateWorkGroupDao(client *gorm.DB) *WorkGroupDao {
	return &WorkGroupDao{db: client}
}

func (dao *WorkGroupDao) FindByWorkplaceCode(workplaceCode string) []*entity.WorkGroupEntity {
	workGroups := make([]*entity.WorkGroupEntity, 0)
	dao.db.Model(entity.WorkGroupEntity{}).Where("workplace_code = ?", workplaceCode).Find(&workGroups)

	return workGroups
}

func (dao *WorkGroupDao) FindByCode(workGroupCode, workplaceCode string) *entity.WorkGroupEntity {
	list := make([]*entity.WorkGroupEntity, 0)
	dao.db.Model(entity.WorkGroupEntity{}).Where("code = ? and workplace_code = ?", workGroupCode, workplaceCode).First(&list)
	if len(list) > 0 {
		return list[0]
	}
	return nil
}
