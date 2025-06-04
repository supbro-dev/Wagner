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
	"wagner/app/global/my_const"
	"wagner/app/global/my_error"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/utils/json_util"
	"wagner/app/utils/log"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type ActionService struct {
	actionDao *dao.ActionDao
}

func CreateActionService(actionDao *dao.ActionDao) *ActionService {
	return &ActionService{actionDao: actionDao}
}

func (service *ActionService) FindEmployeeActions(employeeNumber string, operateDayList []time.Time, originalFieldParam *calc_dynamic_param.OriginalField) *[]domain.Action {
	actionList := service.actionDao.FindBy(employeeNumber, operateDayList)
	return convertAction(&actionList, originalFieldParam)
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func convertAction(actionEntities *[]entity.ActionEntity, param *calc_dynamic_param.OriginalField) *[]domain.Action {
	var actions []domain.Action

	for _, e := range *actionEntities {
		action := domain.Action{}
		copyErr := copier.Copy(&action, &e)
		if copyErr != nil {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}
		operateDay, err := time.Parse(my_const.DateLayout, e.OperateDay)
		if err == nil {
			action.OperateDay = operateDay
		} else {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}
		startTime, err := time.Parse(my_const.DateTimeMsLayout, e.StartTime)
		if err == nil {
			action.StartTime = startTime
		} else {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}
		endTime, err := time.Parse(my_const.DateTimeMsLayout, e.EndTime)
		if err == nil {
			action.EndTime = endTime
		} else {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}

		handleExtraProperty(&action, &e, param)
		actions = append(actions, action)
	}

	return &actions
}

// 如果配置了数据来源有额外属性，在这个方法设置
// Parameters: domainAction, entityAction, param配置参数
func handleExtraProperty(domain *domain.Action, entity *entity.ActionEntity, param *calc_dynamic_param.OriginalField) {
	if param.FieldSet.IsEmpty() {
		return
	}

	json, err := json_util.Parse2Json(entity.Properties)

	if err != nil {
		panic(err)
	}

	// 获取整个JSON对象为map
	dataMap, err := json.Map()
	if err != nil {
		panic(err)
	}

	domain.Properties = make(map[string]interface{})

	for key := range dataMap {
		if param.FieldSet.Contains(key) {
			domain.Properties[key] = dataMap[key]
		}
	}
}
