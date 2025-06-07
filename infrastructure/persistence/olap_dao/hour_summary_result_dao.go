/*
* @Author: supbro
* @Date:   2025/6/6 13:02
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:02
 */
package olap_dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type HourSummaryResultDao struct {
	olapDb *gorm.DB
}

func CreateHourSummaryResultDao(olapClient *gorm.DB) *HourSummaryResultDao {
	return &HourSummaryResultDao{olapClient}
}

const batchSize = 500

func (dao *HourSummaryResultDao) BatchInsert(resultList []*entity.HourSummaryResultEntity) {
	dao.olapDb.Omit("gmt_create", "gmt_modified").CreateInBatches(resultList, batchSize)
}
