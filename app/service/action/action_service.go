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
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/json_util"
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
func (service *ActionService) FindEmployeeActions(employeeNumber string, operateDayList []time.Time, originalFieldParam calc_dynamic_param.InjectSource) (day2WorkList map[time.Time][]domain.Actionable,
	day2Attendance map[time.Time]*domain.Attendance,
	day2Scheduling map[time.Time]*domain.Scheduling,
	day2RestList map[time.Time][]*domain.Rest) {
	actionList := service.actionDao.FindBy(employeeNumber, operateDayList)

	list, attendance, scheduling, restList := service.convertAction(actionList, originalFieldParam)

	return list, attendance, scheduling, restList
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func (service *ActionService) convertAction(actionEntities []*entity.ActionEntity, param calc_dynamic_param.InjectSource) (
	day2WorkList map[time.Time][]domain.Actionable,
	day2Attendance map[time.Time]*domain.Attendance,
	day2Scheduling map[time.Time]*domain.Scheduling,
	day2RestList map[time.Time][]*domain.Rest) {

	day2WorkList = make(map[time.Time][]domain.Actionable)
	day2Attendance = make(map[time.Time]*domain.Attendance)
	day2Scheduling = make(map[time.Time]*domain.Scheduling)
	day2RestList = make(map[time.Time][]*domain.Rest)

	for _, e := range actionEntities {
		actionType := e.ActionType
		properties := handleExtraProperty(e.Properties, param)
		operateDay := e.OperateDay

		switch domain.ActionType(actionType) {
		case domain.DIRECT_WORK:
			raw, err := json_util.Parse2Map(e.WorkLoad)
			if err != nil {
				panic(err)
			}

			workLoad := make(map[string]float64)
			for key, value := range raw {
				workLoad[key] = value.(float64)
			}

			work := domain.DirectWork{WorkLoad: workLoad, Action: domain.Action{Properties: properties}}

			copier.Copy(&work, &e)
			// 这里先设置计算后的时间为原始时间，看之后是否需要去掉
			work.ComputedStartTime = work.StartTime
			work.ComputedEndTime = work.EndTime

			if day2WorkList[work.OperateDay] == nil {
				day2WorkList[work.OperateDay] = make([]domain.Actionable, 0)
			}
			day2WorkList[work.OperateDay] = append(day2WorkList[work.OperateDay], &work)
		case domain.INDIRECT_WORK:
			work := domain.IndirectWork{Action: domain.Action{Properties: properties}}

			copier.Copy(&work, &e)
			// 这里先设置计算后的时间为原始时间，看之后是否需要去掉
			work.ComputedStartTime = work.StartTime
			work.ComputedEndTime = work.EndTime

			if day2WorkList[work.OperateDay] == nil {
				day2WorkList[work.OperateDay] = make([]domain.Actionable, 0)
			}
			day2WorkList[work.OperateDay] = append(day2WorkList[work.OperateDay], &work)
		case domain.SCHEDULING:
			scheduling := domain.Scheduling{Action: domain.Action{Properties: properties}}

			copier.Copy(&scheduling, &e)
			day2Scheduling[operateDay] = &scheduling

			if restList := service.convertRestListFromScheduling(&scheduling); restList != nil {
				day2RestList[operateDay] = restList
			}
		case domain.ATTENDANCE:
			attendance := domain.Attendance{Action: domain.Action{Properties: properties}}

			copier.Copy(&attendance, &e)
			day2Attendance[operateDay] = &attendance
		default:

		}
	}
	return
}

func (service *ActionService) convertRestListFromScheduling(scheduling *domain.Scheduling) []*domain.Rest {
	if restListValue, exist := scheduling.Properties["restList"]; exist {
		array, err := json_util.Parse2JsonArray(restListValue.(string))
		if err != nil {
			panic(err)
		}

		restList := make([]*domain.Rest, 0)
		for i := 0; i < len(array.MustArray()); i++ {
			json := array.GetIndex(i)
			startTime, err := datetime_util.ParseDatetime(json.Get("startTime").MustString())
			if err != nil {
				panic(err)
			}

			endTime, err := datetime_util.ParseDatetime(json.Get("endTime").MustString())
			if err != nil {
				panic(err)
			}
			rest := domain.Rest{
				Action: domain.Action{
					StartTime: &startTime,
					EndTime:   &endTime,
				},
			}
			restList = append(restList, &rest)
		}
		return restList
	} else {
		return nil
	}
}

// 如果配置了数据来源有额外属性，在这个方法设置
// Parameters: properties原始属性, param配置参数
// return: 过滤后的属性
func handleExtraProperty(properties string, param calc_dynamic_param.InjectSource) map[string]interface{} {
	if param.FieldSet.IsEmpty() || properties == "" {
		return nil
	}

	domainProperties := make(map[string]interface{})

	propertyMap, err := json_util.Parse2Map(properties)
	if err != nil {
		panic(err)
	}
	for key, value := range propertyMap {
		if param.FieldSet.Contains(key) {
			domainProperties[key] = value
		}
	}

	return domainProperties
}
