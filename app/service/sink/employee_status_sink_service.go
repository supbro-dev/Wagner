/*
* @Author: supbro
* @Date:   2025/6/18 11:44
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/18 11:44
 */
package sink

import (
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type EmployeeStatusSinkService struct {
	employeeStatusDao *dao.EmployeeStatusDao
}

func CreateEmployeeStatusSinkService(employeeStatusDao *dao.EmployeeStatusDao) *EmployeeStatusSinkService {
	return &EmployeeStatusSinkService{employeeStatusDao: employeeStatusDao}
}

func (s EmployeeStatusSinkService) InsertOrUpdateEmployeeStatus(employeeStatusEntity *entity.EmployeeStatusEntity) {
	s.employeeStatusDao.InsertOrUpdate(employeeStatusEntity)
}
