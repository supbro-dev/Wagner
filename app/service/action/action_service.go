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
	"wagner/app/global/business_error"
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/json_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type ActionService struct {
	actionDao *dao.ActionDao
}

type Day2ActionData struct {
	Day2WorkList   map[time.Time][]domain.Actionable
	Day2Attendance map[time.Time]*domain.Attendance
	Day2Scheduling map[time.Time]*domain.Scheduling
	Day2RestList   map[time.Time][]*domain.Rest
}

func CreateActionService(actionDao *dao.ActionDao) *ActionService {
	return &ActionService{actionDao: actionDao}
}

// 根据工号和日期列表查找动作，并转换成动作对应子类型
// Parameters: employeeNumber，operateDayList 最近3天列表，originalFieldParam 属性映射关系
// Returns: 天2动作列表
func (service *ActionService) FindEmployeeActions(employeeNumber string, operateDayList []time.Time, originalFieldParam calc_dynamic_param.InjectSource) (*Day2ActionData, *business_error.BusinessError) {
	actionList := service.actionDao.FindBy(employeeNumber, operateDayList)

	return service.convertAction(actionList, originalFieldParam)
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func (service *ActionService) convertAction(actionEntities []*entity.ActionEntity, param calc_dynamic_param.InjectSource) (*Day2ActionData, *business_error.BusinessError) {
	day2WorkList := make(map[time.Time][]domain.Actionable)
	day2Attendance := make(map[time.Time]*domain.Attendance)
	day2Scheduling := make(map[time.Time]*domain.Scheduling)
	day2RestList := make(map[time.Time][]*domain.Rest)

	for _, e := range actionEntities {
		actionType := e.ActionType
		properties, err := handleExtraProperty(e.Properties, param)
		if err != nil {
			return nil, err
		}
		operateDay := e.OperateDay

		switch domain.ActionType(actionType) {
		case domain.DIRECT_WORK:
			raw, err := json_util.Parse2Map(e.WorkLoad)
			if err != nil {
				return nil, business_error.InjectDataError(err)
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

			// 这里先设置计算后的时间为原始时间，看之后是否需要去掉
			scheduling.ComputedStartTime = scheduling.StartTime
			scheduling.ComputedEndTime = scheduling.EndTime

			restList, err := service.convertRestListFromScheduling(&scheduling, e.Properties)
			if err != nil {
				return nil, err
			}
			if len(restList) > 0 {
				day2RestList[operateDay] = restList
			}
		case domain.ATTENDANCE:
			attendance := domain.Attendance{Action: domain.Action{Properties: properties}}

			copier.Copy(&attendance, &e)
			attendance.ComputedStartTime = attendance.StartTime
			attendance.ComputedEndTime = attendance.EndTime
			day2Attendance[operateDay] = &attendance
		default:

		}
	}
	return &Day2ActionData{day2WorkList, day2Attendance, day2Scheduling, day2RestList}, nil
}

func (service *ActionService) convertRestListFromScheduling(scheduling *domain.Scheduling, properties string) ([]*domain.Rest, *business_error.BusinessError) {
	if properties == "" {
		return make([]*domain.Rest, 0), nil
	}

	json, err := json_util.Parse2Json(properties)
	if err != nil {
		return nil, business_error.InjectDataError(err)
	}
	array, exists := json.CheckGet("restList")
	if !exists {
		return nil, nil
	}

	restList := make([]*domain.Rest, 0)
	for i := 0; i < len(array.MustArray()); i++ {
		r := array.GetIndex(i)
		startTime, err := datetime_util.ParseDatetime(r.Get("startTime").MustString())
		if err != nil {
			return nil, business_error.InjectDataError(err)
		}

		endTime, err := datetime_util.ParseDatetime(r.Get("endTime").MustString())
		if err != nil {
			return nil, business_error.InjectDataError(err)
		}
		rest := domain.Rest{
			Action: domain.Action{
				EmployeeNumber:    scheduling.EmployeeNumber,
				WorkplaceCode:     scheduling.WorkplaceCode,
				OperateDay:        scheduling.OperateDay,
				ActionType:        domain.REST,
				StartTime:         &startTime,
				EndTime:           &endTime,
				ComputedStartTime: &startTime,
				ComputedEndTime:   &endTime,
			},
		}
		restList = append(restList, &rest)
	}
	return restList, nil
}

// 如果配置了数据来源有额外属性，在这个方法设置
// Parameters: properties原始属性, param配置参数
// return: 过滤后的属性
func handleExtraProperty(properties string, param calc_dynamic_param.InjectSource) (map[string]interface{}, *business_error.BusinessError) {
	if param.FieldSet == nil || param.FieldSet.IsEmpty() || properties == "" {
		return nil, nil
	}

	domainProperties := make(map[string]interface{})

	propertyMap, err := json_util.Parse2Map(properties)
	if err != nil {
		return nil, business_error.InjectDataError(err)
	}
	for key, value := range propertyMap {
		if param.FieldSet.Contains(key) {
			domainProperties[key] = value
		}
	}

	return domainProperties, nil
}
