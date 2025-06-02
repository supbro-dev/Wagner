/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package calc_dynamic_param

import (
	"wagner/app/domain"
	"wagner/app/global/container"
	"wagner/app/global/my_const"
	"wagner/infrastructure/persistence/dao"
)

// 人效计算参数
type CalcParam struct {
	DimensionStorageField DimensionStorageField
	OriginalField         OriginalField
	AggregateField        AggregateField
	CalcNodeList          CalcNodeList
	CalcOtherParam        CalcOtherParam
}

// 不同维度的存储属性
type DimensionStorageField struct {
	// 存储类型
	SinkType my_const.SinkType
	// 表名
	tableName string
	// 属性名转化成字段名
	FieldName2ColumnName map[string]string
}

// action动态字段配置,以及原始字段名转化
type OriginalField struct {
	// 需要关注的原始属性名，其他的直接丢弃
	FieldSet map[string]int
	// 把原始字段名转化成Action的属性名（如需要）
	OriginalFieldName2FieldName map[string]string
}

// 聚合存储时的聚合维度（除employeeNumber + operateDay之外）
type AggregateField struct {
	// 需要聚合的属性名
	FieldSet map[string]int
}

type CalcNodeList struct {
	List []CalcNode
}

// 计算节点
type CalcNode struct {
	NodeName string
	NodeType my_const.NodeType
	// 计算加工节点标准输入输出
	invoker func(domain.ComputeContext, domain.ComputeContext)
}

// 其他各类参数
type CalcOtherParam struct {
	params map[string]interface{}
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
func (service CalcDynamicParamService) FindParamsByWorkplace(workplaceCode string) *CalcParam {
	workplaceEntity := service.workplaceDao.FindByCode(workplaceCode)
	if &workplaceEntity == nil {
		return nil
	}

	paramList := service.calcDynamicParamDao.FindByIndustry(workplaceEntity.IndustryCode, workplaceEntity.SubIndustryCode)

	for _, param := range paramList {

	}

	return nil
}
