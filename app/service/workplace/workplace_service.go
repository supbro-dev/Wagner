/*
* @Author: supbro
* @Date:   2025/6/6 13:19
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:19
 */
package workplace

import (
	"github.com/jinzhu/copier"
	"wagner/app/domain"
	"wagner/infrastructure/persistence/dao"
)

type WorkplaceService struct {
	workplaceDao *dao.WorkplaceDao
}

func CreateWorkplaceService(workplaceDao *dao.WorkplaceDao) *WorkplaceService {
	return &WorkplaceService{workplaceDao: workplaceDao}
}

func (service *WorkplaceService) FindAll() []*domain.Workplace {
	workplaceList := service.workplaceDao.FindAll()

	workplaces := make([]*domain.Workplace, 0)
	for _, workplaceEntity := range workplaceList {
		workplace := domain.Workplace{}
		copier.Copy(&workplace, &workplaceEntity)

		workplaces = append(workplaces, &workplace)
	}

	return workplaces
}

func (service *WorkplaceService) FindByCode(code string) *domain.Workplace {
	workplaceEntity := service.workplaceDao.FindByCode(code)

	if workplaceEntity == nil {
		return nil
	}

	workplace := domain.Workplace{}

	copier.Copy(&workplace, &workplaceEntity)

	return &workplace
}
