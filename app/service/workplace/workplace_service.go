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

// 暂时从工作点表里获取
func (service *WorkplaceService) FindAllIndustry() []string {
	workplaceList := service.workplaceDao.FindAll()

	industries := make([]string, 0)
	for _, workplaceEntity := range workplaceList {
		industries = append(industries, workplaceEntity.IndustryCode)
	}

	return industries
}

// 暂时从工作点表里获取
func (service *WorkplaceService) FindAllSubIndustry() []string {
	workplaceList := service.workplaceDao.FindAll()

	subIndustries := make([]string, 0)
	for _, workplaceEntity := range workplaceList {
		subIndustries = append(subIndustries, workplaceEntity.SubIndustryCode)
	}

	return subIndustries
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

// 暂时通过工作点获取，实际应该有元数据表
func (service *WorkplaceService) FindIndustryBySubIndustry(subIndustry string) string {
	return service.workplaceDao.FindSubIndustryBySubindustryCode(subIndustry)
}
