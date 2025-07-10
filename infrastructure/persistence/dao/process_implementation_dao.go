/*
* @Author: supbro
* @Date:   2025/7/1 13:03
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 13:03
 */
package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/query"
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
		d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'subIndustry' and target_code=? and status = 'online'", subIndustryCode).First(processImplementation)
		if processImplementation != nil {
			return processImplementation
		}
	}
	d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'industry' and target_code=? and status = 'online'", industryCode).First(processImplementation)
	return processImplementation
}

func (d *ProcessImplementationDao) FindByWorkplaceCode(workplaceCode string) *entity.ProcessImplementationEntity {
	processImplementation := &entity.ProcessImplementationEntity{}
	d.db.Model(entity.ProcessImplementationEntity{}).Where("status = 'online' and target_type = 'workplace' and target_code=? and status = 'online'", workplaceCode).First(processImplementation)
	return processImplementation
}

func (d *ProcessImplementationDao) QueryProcessImplementation(query query.ProcessImplementationQuery) []*entity.ProcessImplementationEntity {
	processImplementationList := make([]*entity.ProcessImplementationEntity, 0)
	tx := d.db.Model(entity.ProcessImplementationEntity{}).
		Limit(query.PageSize).
		Offset((query.CurrentPage - 1) * query.PageSize)
	if query.TargetCode != "" {
		tx.Where("target_code = ? and target_type = ?", query.TargetCode, query.TargetType)
	} else {
		tx.Where("target_type = ?", query.TargetType)
	}
	tx.Find(&processImplementationList)

	return processImplementationList
}

func (d *ProcessImplementationDao) CountProcessImplementation(query query.ProcessImplementationQuery) int {
	var total int
	tx := d.db.Model(entity.ProcessImplementationEntity{}).
		Select("count(1)").
		Limit(query.PageSize).
		Offset((query.CurrentPage - 1) * query.PageSize)
	if query.TargetCode != "" {
		tx.Where("target_code = ? and target_type = ?", query.TargetCode, query.TargetType)
	} else {
		tx.Where("target_type = ?", query.TargetType)
	}
	tx.Find(&total)

	return total
}
func (d *ProcessImplementationDao) FindByTarget(targetType entity.TargetType, targetCode string) []*entity.ProcessImplementationEntity {
	processImplementation := make([]*entity.ProcessImplementationEntity, 0)
	d.db.Model(entity.ProcessImplementationEntity{}).Where("target_type = ? and target_code=?", targetType, targetCode).Find(&processImplementation)
	return processImplementation
}

// 根据targetType,targetCode,code查找相同的环节实施
func (d *ProcessImplementationDao) FindOne(q *query.ProcessImplementationQuery) *entity.ProcessImplementationEntity {
	var processImplementation entity.ProcessImplementationEntity
	tx := d.db.Model(entity.ProcessImplementationEntity{}).Where("target_type = ? and target_code = ? and code = ?", q.TargetType, q.TargetCode, q.Code).First(&processImplementation)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil // 记录不存在时返回 nil
	}

	return &processImplementation
}

// 保存或更新基础信息
func (d *ProcessImplementationDao) Save(impl *entity.ProcessImplementationEntity) int64 {
	d.db.Model(entity.ProcessImplementationEntity{}).Omit("gmt_create", "gmt_modified").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "target_code, target_type, code"}}, // 冲突检测列（唯一索引或主键）
		DoUpdates: clause.AssignmentColumns([]string{"name", "status"}),      // 更新字段
	}).Create(impl)

	return impl.Id
}

// 根据id查找环节实施信息
func (d *ProcessImplementationDao) FindById(id int64) *entity.ProcessImplementationEntity {
	processImplementation := entity.ProcessImplementationEntity{}
	d.db.Model(entity.ProcessImplementationEntity{}).Where("id = ?", id).First(&processImplementation)
	return &processImplementation
}

func (d *ProcessImplementationDao) ChangeStatusById(id int64, status entity.ImplementationStatus) {
	d.db.Model(entity.ProcessImplementationEntity{}).Where("id = ?", id).Update("status", status)
}
