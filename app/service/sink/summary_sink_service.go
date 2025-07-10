/*
* @Author: supbro
* @Date:   2025/6/6 13:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:10
 */
package sink

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"wagner/app/domain"
	"wagner/app/global/cache"
	"wagner/app/utils/json_util"
	"wagner/app/utils/md5_util"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/olap_dao"
	"wagner/infrastructure/persistence/query"
)

type SummarySinkService struct {
	hourSummaryResultDao *olap_dao.HourSummaryResultDao
	cache                *cache.HourSummaryCheckCache
}

// 通过构造函数注入 DAO
func CreateSummarySinkService(hourSummaryResultDao *olap_dao.HourSummaryResultDao, cache *cache.HourSummaryCheckCache) *SummarySinkService {
	return &SummarySinkService{hourSummaryResultDao: hourSummaryResultDao, cache: cache}
}

func (service *SummarySinkService) BatchInsertSummaryResult(resultList []*domain.HourSummaryResult, employee *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time) {
	entityList := service.convertDomain2Entity(resultList, employee, workplace, operateDay)

	// 排序
	sort.Slice(entityList, func(i, j int) bool {
		this := entityList[i]
		that := entityList[j]
		if this.OperateTime != that.OperateTime {
			return this.OperateTime.Before(that.OperateTime)
		} else {
			return strings.Compare(this.UniqueKey, that.UniqueKey) > 0
		}
	})

	md5 := md5_util.Md5(json_util.ToJsonString(entityList))

	if md5Value, exists := service.cache.GutResultMd5(employee.Number, workplace.Code, operateDay); exists && md5 == md5Value {
		// 计算结果相等，无需重复写库
		return
	} else {
		service.updateDeleted(&query.HourSummaryResultDelete{EmployeeNumber: employee.Number, WorkplaceCode: workplace.Code, OperateDay: operateDay}, resultList)

		service.hourSummaryResultDao.BatchInsertOrUpdateByUnqKey(entityList)

		service.cache.PutResultMd5(employee.Number, workplace.Code, operateDay, md5)
	}
}

// 尝试把没有被更新的数据进行逻辑删除
func (service *SummarySinkService) updateDeleted(delete *query.HourSummaryResultDelete, resultList []*domain.HourSummaryResult) {
	uniqueKeyList := make([]string, 0)
	for _, result := range resultList {
		uniqueKey := service.generateUniqueKey(&result.AggregateKey)
		uniqueKeyList = append(uniqueKeyList, uniqueKey)
	}

	delete.UniqueKeyList = uniqueKeyList

	service.hourSummaryResultDao.UpdateDeletedByUniqueKeyList(delete)
}

func (service *SummarySinkService) convertDomain2Entity(resultList []*domain.HourSummaryResult, employee *domain.EmployeeSnapshot, workplace *domain.Workplace, operateDay time.Time) []*entity.HourSummaryResultEntity {
	list := make([]*entity.HourSummaryResultEntity, 0)
	for _, d := range resultList {
		e := entity.HourSummaryResultEntity{
			OperateTime:          d.AggregateKey.OperateTime,
			OperateDay:           operateDay,
			ProcessCode:          d.AggregateKey.ProcessCode,
			PositionCode:         service.getPositionCode(d.Process),
			ProcessProperty:      service.convert2ProcessProperty(d.Process),
			WorkplaceCode:        workplace.Code,
			WorkplaceName:        workplace.Name,
			EmployeeNumber:       employee.Number,
			EmployeeName:         employee.Name,
			EmployeePositionCode: employee.PositionCode,
			WorkGroupCode:        employee.WorkGroupCode,
			RegionCode:           workplace.RegionCode,
			IndustryCode:         workplace.IndustryCode,
			SubIndustryCode:      workplace.SubIndustryCode,
			WorkLoad:             json_util.ToJsonString(d.WorkLoad),
			DirectWorkTime:       d.DirectWorkTime,
			IndirectWorkTime:     d.IndirectWorkTime,
			IdleTime:             d.IdleTime,
			RestTime:             d.RestTime,
			AttendanceTime:       d.AttendanceTime,
			Properties:           json_util.ToJsonString(d.Properties),
			UniqueKey:            service.generateUniqueKey(&d.AggregateKey),
		}

		list = append(list, &e)
	}

	return list
}

func (service *SummarySinkService) getPositionCode(process *domain.ProcessPosition) string {
	position := process.Path[0]
	return position.Code
}

func (service *SummarySinkService) convert2ProcessProperty(process *domain.ProcessPosition) string {
	json := json_util.NewJson()
	json.Set("name", process.Name)

	for i, p := range process.Path {
		if i == 0 {
			json.Set("positionCode", p.Code)
			json.Set("positionName", p.Name)
		} else {
			json.Set(fmt.Sprint("deptCode", process.MaxDeptLevel-(i-1)), p.Code)
			json.Set(fmt.Sprint("deptName", process.MaxDeptLevel-(i-1)), p.Name)
		}
	}

	return json_util.ToString(json)
}

func (service *SummarySinkService) generateUniqueKey(key *domain.HourSummaryAggregateKey) string {
	return strings.Join([]string{key.EmployeeNumber, key.OperateTime.String(), key.ProcessCode, key.WorkplaceCode, key.PropertyValues}, "")
}
