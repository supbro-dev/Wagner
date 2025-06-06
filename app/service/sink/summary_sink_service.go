/*
* @Author: supbro
* @Date:   2025/6/6 13:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:10
 */
package sink

import (
	"time"
	"wagner/app/domain"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/olap_dao"
)

type SummarySinkService struct {
	hourSummaryResultDao *olap_dao.HourSummaryResultDao
}

// 通过构造函数注入 DAO
func CreateSummarySinkService(hourSummaryResultDao *olap_dao.HourSummaryResultDao) *SummarySinkService {
	return &SummarySinkService{hourSummaryResultDao: hourSummaryResultDao}
}

func (service *SummarySinkService) BatchInsertSummaryResult(resultList *[]domain.HourSummaryResult, employee *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time) {
	entityList := service.convertDomain2Entity(resultList, employee, workplace, operateDay)

	service.hourSummaryResultDao.BatchInsert(entityList)
}

func (service *SummarySinkService) convertDomain2Entity(resultList *[]domain.HourSummaryResult, employee *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time) *[]entity.HourSummaryResultEntity {
	list := make([]entity.HourSummaryResultEntity, 0)
	for _, d := range *resultList {
		e := entity.HourSummaryResultEntity{
			OperateTime: d.AggregateKey.OperateTime,
			OperateDay:  operateDay,
			ProcessCode: d.AggregateKey.ProcessCode,
			//todo 处理环节额外属性
			ProcessProperty: d.Process.Name,
		}

		list = append(list, e)
	}

	return &list
}
