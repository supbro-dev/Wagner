/*
* @Author: supbro
* @Date:   2025/6/11 09:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:15
 */
package service

import (
	"github.com/jinzhu/copier"
	"time"
	"wagner/app/domain"
	"wagner/app/http/vo"
	"wagner/app/service/calc_dynamic_param"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/olap_dao"
	"wagner/infrastructure/persistence/query"
)

type EfficiencyService struct {
	dao *olap_dao.HourSummaryResultDao
}

func CreateEfficiencyService(dao *olap_dao.HourSummaryResultDao) *EfficiencyService {
	return &EfficiencyService{dao}
}

func (service *EfficiencyService) EmployeeEfficiency(workplaceCode, employeeNumber string, dateRange []*time.Time, aggregateDimension domain.AggregateDimension, isCrossPosition domain.IsCrossPosition, workLoadUnits []calc_dynamic_param.WorkLoadUnit) *vo.EmployeeEfficiencyVO {
	resultQuery := query.HourSummaryResultQuery{WorkplaceCode: workplaceCode, EmployeeNumber: employeeNumber, DateRange: dateRange, AggregateDimension: aggregateDimension, IsCrossPosition: isCrossPosition, WorkLoadUnit: workLoadUnits}
	employeeSummaryEntities := service.dao.QueryEmployeeEfficiency(resultQuery)

	employeeEfficiencyVO := service.convertEntity2Vo(employeeSummaryEntities, workLoadUnits)
	return employeeEfficiencyVO
}

func (service *EfficiencyService) convertEntity2Vo(entityList []*entity.WorkLoadWithSummaryEntity, workLoadUnits []calc_dynamic_param.WorkLoadUnit) *vo.EmployeeEfficiencyVO {
	tableDataList := make([]*vo.EmployeeSummaryVO, 0)
	for _, e := range entityList {
		employeeSummary := vo.EmployeeSummaryVO{}
		copier.Copy(&employeeSummary, &e.EmployeeSummary)
		if e.WorkLoad != nil && len(e.WorkLoad) > 0 {
			employeeSummary.WorkLoad = e.WorkLoad
		}
		tableDataList = append(tableDataList, &employeeSummary)
	}

	columns := service.generateColumns(workLoadUnits)

	v := vo.EmployeeEfficiencyVO{
		tableDataList, columns,
	}

	return &v
}

func (service *EfficiencyService) generateColumns(workLoadUnits []calc_dynamic_param.WorkLoadUnit) []*vo.TableColumnVO {
	columns := []*vo.TableColumnVO{
		{"日期", "operateDay", "operateDay"},
		{"工号", "employeeNumber", "employeeNumber"},
		{"姓名", "employeeName", "employeeName"},
		{"工作点", "workplaceName", "workplaceName"},
		{"作业环节", "processName", "processName"},
		{"作业岗位", "positionName", "positionName"},
		{"部门", "deptName", "deptName"},
	}

	for _, unit := range workLoadUnits {
		columns = append(columns, &vo.TableColumnVO{
			unit.Name, []string{"workLoad", unit.Code}, unit.Code,
		})
	}

	return columns
}
