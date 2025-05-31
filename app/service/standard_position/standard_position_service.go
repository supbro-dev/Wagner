package standard_position

import (
	"wagner/app/domain"
	"wagner/infrastructure/persistence/dao"
)

type StandardPositionService struct {
	StandardPositionDao *dao.StandardPositionDao
}

func CreateStandardPositionService(standardPositionDao *dao.StandardPositionDao) *StandardPositionService {
	return &StandardPositionService{}
}

//	根据工作点编码获取工序标准模型
//
// Parameters:
// Returns:
func (service *StandardPositionService) FindStandardPositionByWorkplace(workplaceCode string) []domain.StandardPosition {
	// todo 这里应该查找工序实施配置，在没有实施流程时，直接根据工作点查找行业的标准模型
}
