/*
* @Author: supbro
* @Date:   2025/6/6 13:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:10
 */
package sink

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/utils/json_util"
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
			OperateTime:          d.AggregateKey.OperateTime,
			OperateDay:           operateDay,
			ProcessCode:          d.AggregateKey.ProcessCode,
			PositionCode:         service.getPositionCode(&d.Process),
			ProcessProperty:      service.convert2ProcessProperty(&d.Process),
			WorkplaceCode:        workplace.Code,
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
			AttendanceTime:       d.AttendanceTime,
			Properties:           json_util.ToJsonString(d.Properties),
		}

		list = append(list, e)
	}

	return &list
}

func (service *SummarySinkService) getPositionCode(process *domain.StandardPosition) string {
	position := process.Path[0]
	return position.Code
}

func (service *SummarySinkService) convert2ProcessProperty(process *domain.StandardPosition) string {
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

	s, err := json.String()
	if err != nil {
		panic(err)
	}

	return s
}
