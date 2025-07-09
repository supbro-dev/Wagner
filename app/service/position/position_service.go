/*
* @Author: supbro
* @Date:   2025/7/9 11:28
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/9 11:28
 */
package position

import (
	"github.com/jinzhu/copier"
	"wagner/app/domain"
	"wagner/infrastructure/persistence/dao"
)

type PositionService struct {
	positionDao *dao.PositionDao
}

func (s *PositionService) FindAll() []*domain.Position {
	positionEntityList := s.positionDao.FindAll()

	positionList := make([]*domain.Position, 0)
	for _, e := range positionEntityList {
		position := &domain.Position{}
		copier.Copy(&position, e)
		positionList = append(positionList, position)
	}

	return positionList
}

func CreatePositionService(positionDao *dao.PositionDao) *PositionService {
	return &PositionService{positionDao: positionDao}
}
