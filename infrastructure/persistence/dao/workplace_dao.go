/*
* @Author: supbro
* @Date:   2025/6/2 11:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 11:29
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type WorkplaceDao struct {
	db *gorm.DB
}

func (dao *WorkplaceDao) FindByCode(code string) *entity.WorkplaceEntity {
	workplace := entity.WorkplaceEntity{}
	dao.db.Where("code = ?", code).Find(&workplace)
	return &workplace
}

func (dao *WorkplaceDao) FindAll() []*entity.WorkplaceEntity {
	workplaceList := make([]*entity.WorkplaceEntity, 0)
	dao.db.Find(&workplaceList)
	return workplaceList
}

// 暂时先这么使用，实际需要有单独的行业元数据管理
func (dao *WorkplaceDao) FindSubIndustryBySubindustryCode(subIndustryCode string) string {
	var industryCode string
	dao.db.Table("workplace").Where("sub_industry_code = ?", subIndustryCode).Select("industry_code").First(&industryCode)

	return industryCode
}

func CreateWorkplaceDao(client *gorm.DB) *WorkplaceDao {
	return &WorkplaceDao{client}
}
