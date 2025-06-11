/*
* @Author: supbro
* @Date:   2025/6/11 09:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:15
 */
package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	"math"
	"strings"
	"time"
	"wagner/app/domain"
	"wagner/app/http/vo"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/json_util"
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

	employeeEfficiencyVO := service.convertEntity2Vo(employeeSummaryEntities, workLoadUnits, aggregateDimension)
	return employeeEfficiencyVO
}

func (service *EfficiencyService) convertEntity2Vo(entityList []*entity.WorkLoadWithSummaryEntity, workLoadUnits []calc_dynamic_param.WorkLoadUnit, aggregateDimension domain.AggregateDimension) *vo.EmployeeEfficiencyVO {
	tableDataList := make([]*vo.EmployeeSummaryVO, 0)
	for _, e := range entityList {
		employeeSummary := vo.EmployeeSummaryVO{}
		copier.Copy(&employeeSummary, &e.EmployeeSummary)
		employeeSummary.OperateDay = datetime_util.FormatDate(e.EmployeeSummary.OperateDay)
		employeeSummary.DirectWorkTime = math.Round(employeeSummary.DirectWorkTime*10/3600.0) / 10
		employeeSummary.IndirectWorkTime = math.Round(employeeSummary.IndirectWorkTime*10/3600.0) / 10
		employeeSummary.IdleTime = math.Round(employeeSummary.IdleTime*10/3600.0) / 10
		employeeSummary.RestTime = math.Round(employeeSummary.RestTime*10/3600.0) / 10
		employeeSummary.AttendanceTime = math.Round(employeeSummary.AttendanceTime*10/3600.0) / 10
		if e.EmployeeSummary.ProcessProperty != "" {
			if json, err := json_util.Parse2Map(e.EmployeeSummary.ProcessProperty); err == nil {
				employeeSummary.ProcessName = json["name"].(string)
				employeeSummary.PositionName = json["positionName"].(string)
				employeeSummary.DeptName = service.parseDeptName(json)
			}
		}

		if e.WorkLoad != nil && len(e.WorkLoad) > 0 {
			employeeSummary.WorkLoad = e.WorkLoad
		}
		tableDataList = append(tableDataList, &employeeSummary)
	}

	columns := service.generateColumns(workLoadUnits, aggregateDimension)

	v := vo.EmployeeEfficiencyVO{
		tableDataList, columns,
	}

	return &v
}

func (service *EfficiencyService) parseDeptName(json map[string]interface{}) string {
	deptNameList := make([]string, 0)
	deptNameFmt := "deptName%v"
	for i := 1; i < 10; i++ {
		if deptNameX, exists := json[fmt.Sprintf(deptNameFmt, i)]; exists {
			deptNameList = append(deptNameList, deptNameX.(string))
		} else {
			break
		}
	}

	return strings.Join(deptNameList, "-")
}

func (service *EfficiencyService) generateColumns(workLoadUnits []calc_dynamic_param.WorkLoadUnit, dimension domain.AggregateDimension) []*vo.TableColumnVO {
	columns := []*vo.TableColumnVO{
		{"日期", "operateDay", "operateDay"},
		{"工号", "employeeNumber", "employeeNumber"},
		{"姓名", "employeeName", "employeeName"},
		{"工作点", "workplaceName", "workplaceName"},
	}

	if dimension == domain.Process {
		columns = append(columns, &vo.TableColumnVO{"作业环节", "processName", "processName"})
	}

	columns = append(columns, []*vo.TableColumnVO{
		{"作业岗位", "positionName", "positionName"},
		{"部门", "deptName", "deptName"},
	}...)

	for _, unit := range workLoadUnits {
		columns = append(columns, &vo.TableColumnVO{
			unit.Name, []string{"workLoad", unit.Code}, unit.Code,
		})
	}

	columns = append(columns, []*vo.TableColumnVO{
		{"直接作业工时(h)", "directWorkTime", "directWorkTime"},
		{"间接作业工时(h)", "indirectWorkTime", "indirectWorkTime"},
		{"闲置工时(h)", "idleTime", "idleTime"},
		{"休息时长(h)", "restTime", "restTime"},
		{"出勤工时(h)", "attendanceTime", "attendanceTime"},
	}...)

	return columns
}
