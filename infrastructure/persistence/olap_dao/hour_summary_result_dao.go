/*
* @Author: supbro
* @Date:   2025/6/6 13:02
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:02
 */
package olap_dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wagner/infrastructure/persistence/entity"
)

type HourSummaryResultDao struct {
	olapDb *gorm.DB
}

func CreateHourSummaryResultDao(olapClient *gorm.DB) *HourSummaryResultDao {
	return &HourSummaryResultDao{olapClient}
}

const batchSize = 500

func (dao *HourSummaryResultDao) BatchInsertOrUpdateByUnqKey(resultList []*entity.HourSummaryResultEntity) {
	// todo 如果没有任何字段更新，gmt_modified即便设置了CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP，也不会自动更新，看这里是否需要手动更新该字段
	dao.olapDb.Omit("gmt_create", "gmt_modified").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},                                                                                                                                                                                                                                                                                                                                                                          // 冲突检测列（唯一索引或主键）
		DoUpdates: clause.AssignmentColumns([]string{"operate_time", "operate_day", "process_code", "position_code", "workplace_code", "employee_number", "employee_name", "employee_position_code", "work_group_code", "region_code", "industry_code", "sub_industry_code", "work_load", "direct_work_time", "indirect_work_time", "idle_time", "rest_time", "attendance_time", "process_property", "properties", "is_deleted"}), // 更新字段
	}).CreateInBatches(resultList, batchSize)
}
