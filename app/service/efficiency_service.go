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
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/json_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/olap_dao"
	"wagner/infrastructure/persistence/query"
)

type EfficiencyService struct {
	hourSummaryResultDao *olap_dao.HourSummaryResultDao
	employeeStatusDao    *dao.EmployeeStatusDao
}

func CreateEfficiencyService(hourSummaryResultDao *olap_dao.HourSummaryResultDao, employeeStatusDao *dao.EmployeeStatusDao) *EfficiencyService {
	return &EfficiencyService{hourSummaryResultDao, employeeStatusDao}
}

func (service *EfficiencyService) EmployeeEfficiency(workplaceCode, employeeNumber string, dateRange []*time.Time, aggregateDimension domain.AggregateDimension, isCrossPosition domain.IsCrossPosition, workLoadUnits []calc_dynamic_param.WorkLoadUnit, currentPage, pageSize int) *vo.EmployeeEfficiencyVO {
	resultQuery := query.HourSummaryResultQuery{
		WorkplaceCode:      workplaceCode,
		EmployeeNumber:     employeeNumber,
		DateRange:          dateRange,
		AggregateDimension: aggregateDimension,
		IsCrossPosition:    isCrossPosition,
		WorkLoadUnit:       workLoadUnits,
		CurrentPage:        currentPage,
		PageSize:           pageSize}
	employeeSummaryEntities := service.hourSummaryResultDao.QueryEmployeeEfficiency(resultQuery)

	total := service.hourSummaryResultDao.TotalEmployeeEfficiency(resultQuery)

	employeeEfficiencyVO := service.convertEntity2Vo(employeeSummaryEntities, workLoadUnits, aggregateDimension)
	employeeEfficiencyVO.Page = &vo.Page{CurrentPage: currentPage, PageSize: pageSize, Total: total}

	return employeeEfficiencyVO
}

func (service *EfficiencyService) convertEntity2Vo(entityList []*entity.WorkLoadWithEmployeeSummary, workLoadUnits []calc_dynamic_param.WorkLoadUnit, aggregateDimension domain.AggregateDimension) *vo.EmployeeEfficiencyVO {
	tableDataList := make([]*vo.EmployeeSummaryVO, 0)
	for _, e := range entityList {
		employeeSummary := vo.EmployeeSummaryVO{}
		copier.Copy(&employeeSummary, &e.EmployeeSummary)
		employeeSummary.Key = e.EmployeeSummary.UniqueKey
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

		if employeeSummary.AttendanceTime != float64(0) {
			employeeSummary.DirectWorkTimeRate = math.Round(employeeSummary.DirectWorkTime / employeeSummary.AttendanceTime * 100)
			employeeSummary.IndirectWorkTimeRate = math.Round(employeeSummary.IndirectWorkTime / employeeSummary.AttendanceTime * 100)
			employeeSummary.IdleTimeRate = math.Round(employeeSummary.IdleTime / employeeSummary.AttendanceTime * 100)
		}

		if e.WorkLoad != nil && len(e.WorkLoad) > 0 {
			employeeSummary.WorkLoad = e.WorkLoad
		}
		tableDataList = append(tableDataList, &employeeSummary)
	}

	columns := service.generateEmployeeColumns(workLoadUnits, aggregateDimension)

	v := vo.EmployeeEfficiencyVO{
		TableDataList: tableDataList,
		Columns:       columns,
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

func (service *EfficiencyService) generateEmployeeColumns(workLoadUnits []calc_dynamic_param.WorkLoadUnit, dimension domain.AggregateDimension) []*vo.TableColumnVO {
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
		{"工时占比", "timeRate", "timeRate"},
	}...)

	return columns
}

func (service *EfficiencyService) WorkplaceEfficiency(workplace *domain.Workplace, dateRange []*time.Time, isCrossPosition domain.IsCrossPosition, workLoadUnits []calc_dynamic_param.WorkLoadUnit, standardPositions []*domain.StandardPosition) *vo.WorkplaceEfficiencyVO {
	resultQuery := query.HourSummaryResultQuery{WorkplaceCode: workplace.Code, DateRange: dateRange, IsCrossPosition: isCrossPosition, WorkLoadUnit: workLoadUnits}
	processSummaries := service.hourSummaryResultDao.QueryWorkplaceEfficiency(resultQuery)

	treeRoot := service.buildWorkplaceStructureTree(workplace, standardPositions, processSummaries, workLoadUnits)
	columns := service.generateWorkplaceColumns(workLoadUnits)
	return &vo.WorkplaceEfficiencyVO{treeRoot, columns}
}

func (service *EfficiencyService) buildWorkplaceStructureTree(workplace *domain.Workplace, standardPositions []*domain.StandardPosition, summaries []*entity.WorkLoadWithProcessSummary, workLoadUnits []calc_dynamic_param.WorkLoadUnit) *vo.WorkplaceStructureVO {
	if standardPositions == nil || len(standardPositions) == 0 {
		return nil
	}

	processCode2Summary := make(map[string]*entity.WorkLoadWithProcessSummary, 0)
	for _, summary := range summaries {
		processCode2Summary[summary.ProcessSummary.ProcessCode] = summary
	}

	// 1.构建树
	root := service.convert2Structure(&domain.StandardPosition{
		Name: workplace.Name,
		Code: "-1",
	})
	code2Node := make(map[string]*vo.WorkplaceStructureVO, 0)
	code2Node[root.Code] = root

	for _, position := range standardPositions {
		parentCode := position.ParentCode
		parentNode := code2Node[parentCode]
		node := service.convert2Structure(position)

		if position.Type == entity.DIRECT_PROCESS {
			if summary, exists := processCode2Summary[node.Code]; exists {
				service.mergeTime(node, summary)
				service.mergeWorkLoad(node, summary, workLoadUnits)
			}
		} else if position.Type == entity.INDIRECT_PROCESS {
			if summary, exists := processCode2Summary[node.Code]; exists {
				service.mergeTime(node, summary)
			}
		}

		parentNode.Children = append(parentNode.Children, node)

		code2Node[node.Code] = node
	}

	service.iterateTreeRollUp(root, workLoadUnits)

	return root
}

func (service *EfficiencyService) iterateTreeRollUp(node *vo.WorkplaceStructureVO, workLoadUnits []calc_dynamic_param.WorkLoadUnit) {
	if node == nil {
		return
	}
	if node.Children != nil {
		for _, child := range node.Children {
			service.iterateTreeRollUp(child, workLoadUnits)

			node.DirectWorkTime += child.DirectWorkTime
			node.IndirectWorkTime += child.IndirectWorkTime
			node.IdleTime += child.IdleTime
			node.RestTime += child.RestTime
			node.AttendanceTime += child.AttendanceTime

			// 工作量是否向上汇总
			if child.WorkLoadRollUp {
				for _, workLoadUnit := range workLoadUnits {
					nodeValue := float64(0)
					if n, exists := node.WorkLoad[workLoadUnit.Code]; exists {
						nodeValue = n
					}
					childValue := float64(0)
					if c, exists := child.WorkLoad[workLoadUnit.Code]; exists {
						childValue = c
					}
					node.WorkLoad[workLoadUnit.Code] = nodeValue + childValue
				}
			}
		}
	}
}

var NeedRollUp = "workLoadRollUp"

func (service *EfficiencyService) convert2Structure(position *domain.StandardPosition) *vo.WorkplaceStructureVO {
	v := vo.WorkplaceStructureVO{}
	v.Name = position.Name
	v.Code = position.Code
	v.Level = position.Level
	v.Children = make([]*vo.WorkplaceStructureVO, 0)
	if position.Properties != nil {
		if needRollUp, exists := position.Properties[NeedRollUp]; exists {
			v.WorkLoadRollUp = needRollUp.(bool)
		}
	}
	v.WorkLoad = make(map[string]float64)
	return &v
}

func (service *EfficiencyService) mergeTime(structure *vo.WorkplaceStructureVO, summary *entity.WorkLoadWithProcessSummary) {
	structure.DirectWorkTime += math.Round(float64(summary.ProcessSummary.DirectWorkTime)*10/3600.0) / 10
	structure.IndirectWorkTime += math.Round(float64(summary.ProcessSummary.IndirectWorkTime)*10/3600.0) / 10
	structure.IdleTime += math.Round(float64(summary.ProcessSummary.IdleTime)*10/3600.0) / 10
	structure.RestTime += math.Round(float64(summary.ProcessSummary.RestTime)*10/3600.0) / 10
	structure.AttendanceTime += math.Round(float64(summary.ProcessSummary.AttendanceTime)*10/3600.0) / 10
}

func (service *EfficiencyService) mergeWorkLoad(structure *vo.WorkplaceStructureVO, summary *entity.WorkLoadWithProcessSummary, workLoadUnits []calc_dynamic_param.WorkLoadUnit) {
	if structure.WorkLoad == nil {
		structure.WorkLoad = make(map[string]float64, 0)
	}

	for _, workLoadUnit := range workLoadUnits {
		structureValue := float64(0)
		if sv, exists := structure.WorkLoad[workLoadUnit.Code]; exists {
			structureValue = sv
		}
		summaryValue := float64(0)
		if v, exists := summary.WorkLoad[workLoadUnit.Code]; exists {
			summaryValue = v
		}

		structure.WorkLoad[workLoadUnit.Code] = summaryValue + structureValue
	}
}

func (service *EfficiencyService) generateWorkplaceColumns(workLoadUnits []calc_dynamic_param.WorkLoadUnit) []*vo.TableColumnVO {
	columns := []*vo.TableColumnVO{
		{"名称", "name", "name"},
	}

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

func (service *EfficiencyService) QueryEmployeeStatus(workplaceCode string, operateDay time.Time) *vo.EmployeeStatusVO {
	employeeStatusEntities := service.employeeStatusDao.FindByWorkplaceAndDate(workplaceCode, operateDay)
	// 如果有工作组基本信息，在这里需要进行转化
	employeeStatusVO := service.convertEmployeeStatusEntity2Vo(employeeStatusEntities, nil)
	return employeeStatusVO
}

func (service *EfficiencyService) convertEmployeeStatusEntity2Vo(entities []*entity.EmployeeStatusEntity, groupCode2Name map[string]string) *vo.EmployeeStatusVO {
	groupCodeList := make([]string, 0)
	group2EachEmployeeStatusList := make(map[string][]*vo.EachEmployeeStatusVO)
	group2Count := make(map[string]*vo.GroupStatusNumVO, 0)
	status2Desc := map[string]string{
		"DIRECT_WORKING":            "直接工作中",
		"INDIRECT_WORKING":          "间接工作中",
		"IDLE":                      "空闲中",
		"REST":                      "休息中",
		"OFF_DUTY":                  "已下班",
		"OFF_DUTY_WITHOUT_END_TIME": "下班未打卡",
	}

	defaultWorkGroupName := "默认分组"

	for _, statusEntity := range entities {
		workGroupName := statusEntity.WorkGroupCode
		if statusEntity.WorkGroupCode == "" {
			workGroupName = defaultWorkGroupName
		}

		if groupCode2Name != nil {
			workGroupName = groupCode2Name[statusEntity.WorkGroupCode]
		}

		if _, exists := group2EachEmployeeStatusList[workGroupName]; !exists {
			group2EachEmployeeStatusList[workGroupName] = make([]*vo.EachEmployeeStatusVO, 0)
			groupCodeList = append(groupCodeList, statusEntity.WorkGroupCode)
			group2Count[workGroupName] = &vo.GroupStatusNumVO{}
		}

		eachEmployeeStatusList := group2EachEmployeeStatusList[workGroupName]
		groupStatusNumVO := group2Count[workGroupName]

		var lastActionDesc string
		if statusEntity.Status == entity.DirectWorking {
			lastActionDesc = fmt.Sprintf("于%v开始进行%v直接作业", datetime_util.FormatDatetime(*statusEntity.LastActionTime), statusEntity.LastActionCode)
		} else if statusEntity.Status == entity.IndirectWorking {
			lastActionDesc = fmt.Sprintf("于%v开始进行%v间接作业", datetime_util.FormatDatetime(*statusEntity.LastActionTime), statusEntity.LastActionCode)
		} else if statusEntity.Status == entity.Idle {
			lastActionDesc = fmt.Sprintf("从%v起处于闲置中", datetime_util.FormatDatetime(*statusEntity.LastActionTime))
		} else if statusEntity.Status == entity.Idle {
			lastActionDesc = fmt.Sprintf("从%v起处于休息中", datetime_util.FormatDatetime(*statusEntity.LastActionTime))
		} else if statusEntity.Status == entity.OffDutyWithoutEndTime {
			lastActionDesc = fmt.Sprintf("从%v完成最后一次作业，下班未打卡", datetime_util.FormatDatetime(*statusEntity.LastActionTime))
		} else {
			lastActionDesc = fmt.Sprintf("于%v下班", datetime_util.FormatDatetime(*statusEntity.LastActionTime))
		}
		eachEmployeeStatusList = append(eachEmployeeStatusList, &vo.EachEmployeeStatusVO{
			EmployeeNumber: statusEntity.EmployeeNumber,
			EmployeeName:   statusEntity.EmployeeName,
			StatusDesc:     status2Desc[string(statusEntity.Status)],
			LastActionDesc: lastActionDesc,
		})
		group2EachEmployeeStatusList[workGroupName] = eachEmployeeStatusList

		switch statusEntity.Status {
		case entity.DirectWorking:
			groupStatusNumVO.DirectWorkingNum++
		case entity.IndirectWorking:
			groupStatusNumVO.IndirectWorkingNum++
		case entity.Idle:
			groupStatusNumVO.IdleNum++
		case entity.Rest:
			groupStatusNumVO.RestNum++
		case entity.OffDuty:
			groupStatusNumVO.OffDutyNum++
		case entity.OffDutyWithoutEndTime:
			groupStatusNumVO.OffDutyWithoutEndTimeNum++
		}
	}

	// 按顺序生成小组结果
	groupStatus := make([]*vo.GroupStatusVO, 0)
	for _, groupCode := range groupCodeList {
		groupStatusVO := vo.GroupStatusVO{
			GroupCode: groupCode,
		}
		if groupCode2Name != nil {
			groupStatusVO.GroupName = groupCode2Name[groupCode]
		} else if groupCode == "" {
			groupStatusVO.GroupName = defaultWorkGroupName
		} else {
			groupStatusVO.GroupName = groupCode
		}
		groupStatusVO.GroupStatusNum = group2Count[groupStatusVO.GroupName]
		groupStatusVO.EmployeeStatusList = group2EachEmployeeStatusList[groupStatusVO.GroupName]

		groupStatus = append(groupStatus, &groupStatusVO)
	}

	return &vo.EmployeeStatusVO{
		groupStatus,
	}
}
