package standard_position

import "wagner/infrastructure/persistence/dao"

type StandardPositionService struct {
	StandardPositionDao *dao.StandardPositionDao
}

// 通过构造函数注入 DAO
func CreateStandardPositionService(standardPositionDao *dao.StandardPositionDao) *StandardPositionService {
	return &StandardPositionService{}
}
