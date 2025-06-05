/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package action

import (
	"github.com/jinzhu/copier"
	"time"
	"wagner/app/domain"
	"wagner/app/service/calc_dynamic_param"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type ActionService struct {
	actionDao *dao.ActionDao
}

func CreateActionService(actionDao *dao.ActionDao) *ActionService {
	return &ActionService{actionDao: actionDao}
}

// 根据工号和日期列表查找动作，并转换成动作对应子类型
// Parameters: employeeNumber，operateDayList 最近3天列表，originalFieldParam 属性映射关系
// Returns: 天2动作列表
func (service *ActionService) FindEmployeeActions(employeeNumber string, operateDayList []time.Time, originalFieldParam *calc_dynamic_param.OriginalField) (day2WorkList map[time.Time][]domain.Work,
	day2Attendance map[time.Time]domain.Attendance,
	day2Scheduling map[time.Time]domain.Scheduling) {
	actionList := service.actionDao.FindBy(employeeNumber, operateDayList)

	return convertAction(&actionList, originalFieldParam)
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func convertAction(actionEntities *[]entity.ActionEntity, param *calc_dynamic_param.OriginalField) (
	day2WorkList map[time.Time][]domain.Work,
	day2Attendance map[time.Time]domain.Attendance,
	day2Scheduling map[time.Time]domain.Scheduling) {

	day2WorkList = make(map[time.Time][]domain.Work)
	day2Attendance = make(map[time.Time]domain.Attendance)
	day2Scheduling = make(map[time.Time]domain.Scheduling)

	for _, e := range *actionEntities {
		actionType := e.ActionType
		properties := handleExtraProperty(e.Properties, param)
		operateDay := e.OperateDay

		switch domain.ActionType(actionType) {
		case domain.DIRECT_WORK:
			work := domain.DirectWork{WorkLoad: e.WorkLoad, Action: domain.Action{Properties: *properties}}

			copier.Copy(&work, &e)
			if day2WorkList[work.OperateDay] == nil {
				day2WorkList[work.OperateDay] = make([]domain.Work, 0)
			}
			day2WorkList[work.OperateDay] = append(day2WorkList[work.OperateDay], work)
		case domain.INDIRECT_WORK:
			work := domain.IndirectWork{Action: domain.Action{Properties: *properties}}

			copier.Copy(&work, &e)
			if day2WorkList[work.OperateDay] == nil {
				day2WorkList[work.OperateDay] = make([]domain.Work, 0)
			}
			day2WorkList[work.OperateDay] = append(day2WorkList[work.OperateDay], work)
		case domain.SCHEDULING:
			scheduling := domain.Scheduling{Action: domain.Action{Properties: *properties}}

			copier.Copy(&scheduling, &e)
			day2Scheduling[operateDay] = scheduling
		case domain.ATTENDANCE:
			attendance := domain.Attendance{Action: domain.Action{Properties: *properties}}

			copier.Copy(&attendance, &e)
			day2Attendance[operateDay] = attendance
		default:

		}
	}
	return
}

// 如果配置了数据来源有额外属性，在这个方法设置
// Parameters: properties原始属性, param配置参数
// return: 过滤后的属性
func handleExtraProperty(properties map[string]interface{}, param *calc_dynamic_param.OriginalField) *map[string]interface{} {
	if param.FieldSet.IsEmpty() {
		return nil
	}

	domainProperties := make(map[string]interface{})

	for key, value := range properties {
		if param.FieldSet.Contains(key) {
			domainProperties[key] = value
		}
	}

	return &domainProperties
}
