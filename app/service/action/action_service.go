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

func (service *ActionService) FindEmployeeActions(employeeNumber string, operateDayList []time.Time) []domain.Action {
	actionList := service.actionDao.FindBy(employeeNumber, operateDayList)
	return convertAction(actionList)
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func convertAction(actionEntities []entity.ActionEntity) []domain.Action {
	var actions []domain.Action

	for _, e := range actionEntities {
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
		actions = append(actions, action)
	}

	return actions
}
