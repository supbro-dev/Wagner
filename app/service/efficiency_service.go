/*
* @Author: supbro
* @Date:   2025/6/11 09:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:15
 */
package service

import (
	"fmt"
	"time"
	"wagner/app/domain"
	"wagner/app/http/vo"
	"wagner/infrastructure/persistence/olap_dao"
	"wagner/infrastructure/persistence/query"
)

type EfficiencyService struct {
	dao *olap_dao.HourSummaryResultDao
}

func CreateEfficiencyService(dao *olap_dao.HourSummaryResultDao) *EfficiencyService {
	return &EfficiencyService{dao}
}

func (service *EfficiencyService) EmployeeEfficiency(workplaceCode, employeeNumber string, dateRange []*time.Time, aggregateDimension domain.AggregateDimension, isCrossPosition domain.IsCrossPosition, units []string) *vo.EmployeeEfficiencyVO {
	resultQuery := query.HourSummaryResultQuery{WorkplaceCode: workplaceCode, EmployeeNumber: employeeNumber, DateRange: dateRange, AggregateDimension: aggregateDimension, IsCrossPosition: isCrossPosition, WorkLoadUnit: units}
	employeeSummaryEntities := service.dao.QueryEmployeeEfficiency(resultQuery)

	fmt.Println(employeeSummaryEntities)

	return &vo.EmployeeEfficiencyVO{}
}
