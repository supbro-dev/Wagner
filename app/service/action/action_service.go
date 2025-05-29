package action

import (
	"time"
	"wagner/app/domain"
	"wagner/app/global/my_const"
	"wagner/app/global/my_error"
	"wagner/app/utils/log"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type ActionService struct {
	actionRepo *dao.ActionDao
}

// 通过构造函数注入 DAO
func CreateActionService(actionRepo *dao.ActionDao) *ActionService {
	return &ActionService{actionRepo: actionRepo}
}

func (service *ActionService) FindEmployeeActions(employeeNumber, operateDayList []string) []domain.Action {
	actionList := service.actionRepo.FindBy(employeeNumber, operateDayList)
	return convertAction(actionList)
}

func (service *ActionService) FindWorkplaceActions(workplaceCode, operateDay string) []domain.Action {
	return nil
}

func convertAction(actionEntities []entity.ActionEntity) []domain.Action {
	var actions []domain.Action

	for _, e := range actionEntities {
		action := domain.Action{}
		action.EmployeeNumber = e.EmployeeNumber
		action.OperateDay = e.OperateDay
		action.WorkplaceCode = e.WorkplaceCode
		action.ActionCode = e.ActionCode
		action.Properties = e.Properties
		action.ActionType = e.ActionType
		startTime, err := time.Parse(my_const.DateTimeLayout, e.StartTime)
		if err == nil {
			action.StartTime = startTime
		} else {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}
		endTime, err := time.Parse(my_const.DateTimeLayout, e.EndTime)
		if err == nil {
			action.EndTime = endTime
		} else {
			log.SystemLogger.Error(my_error.ServerOccurredErrorMsg)
		}

		actions = append(actions, action)
	}

	return actions
}
