/*
* @Author: supbro
* @Date:   2025/7/17 14:33
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/17 14:33
 */
package work_group

import (
	"github.com/jinzhu/copier"
	"wagner/app/domain"
	"wagner/infrastructure/persistence/dao"
)

type WorkGroupService struct {
	workGroupDao *dao.WorkGroupDao
}

func CreateWorkGroupService(workGroupDao *dao.WorkGroupDao) *WorkGroupService {
	return &WorkGroupService{workGroupDao}
}

func (service WorkGroupService) FindGroupListByWorkplace(workplaceCode string) []*domain.WorkGroup {
	groupList := service.workGroupDao.FindByWorkplaceCode(workplaceCode)

	groupDomainList := make([]*domain.WorkGroup, 0)
	for _, group := range groupList {
		groupDomain := domain.WorkGroup{}
		copier.Copy(&groupDomain, group)

		groupDomainList = append(groupDomainList, &groupDomain)
	}

	return groupDomainList
}
