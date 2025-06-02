/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package calc_dynamic_param

import (
	"wagner/app/global/container"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type CalcParam struct {
}

type CalcDynamicParamService struct {
	calcDynamicParamDao *dao.CalcDynamicParamDao
	workplaceDao        *dao.WorkplaceDao
	cache               *container.GenericCache[string, CalcParam]
}

func CreateCalcDynamicParamService(calcDynamicParamDao *dao.CalcDynamicParamDao, workplaceDao *dao.WorkplaceDao) *CalcDynamicParamService {
	cache, err := container.GetOrCreateCache[string, CalcParam](container.DYNAMIC_PARAM)
	if err != nil {
		panic(err)
	}
	return &CalcDynamicParamService{calcDynamicParamDao: calcDynamicParamDao, workplaceDao: workplaceDao, cache: cache}
}

// 根据工作点信息获取全量计算参数配置
// Parameters: 工作点编码
// Returns: 计算参数配置列表
func (service CalcDynamicParamService) FindParamsByWorkplace(workplaceCode string) []entity.CalcDynamicParamEntity {
	return nil
}
