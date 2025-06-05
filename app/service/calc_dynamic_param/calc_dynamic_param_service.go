/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package calc_dynamic_param

import (
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
	"wagner/app/global/container"
	"wagner/app/global/my_const"
	"wagner/app/utils/json_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

// 人效计算参数
type CalcParam struct {
	DimensionStorageFields *[]DimensionStorageField
	OriginalField          *OriginalField
	AggregateField         *AggregateField
	CalcNodeList           *CalcNodeList
	CalcOtherParam         *CalcOtherParam
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
	FieldSet mapset.Set[string]
	// 把原始字段名转化成Action的属性名（如需要）
	OriginalFieldName2FieldName map[string]string
}

// 聚合存储时的聚合维度（除employeeNumber + operateDay之外）
type AggregateField struct {
	// 需要聚合的属性名
	FieldSet mapset.Set[string]
}

type CalcNodeList struct {
	List *[]CalcNode
}

// 计算节点
type CalcNode struct {
	NodeName string
	NodeType entity.ScriptType
	// 计算脚本
	Script string
}

// 其他各类参数
type CalcOtherParam struct {
	params interface{}
}

type CalcDynamicParamService struct {
	calcDynamicParamDao *dao.CalcDynamicParamDao
	workplaceDao        *dao.WorkplaceDao
	scriptDao           *dao.ScriptDao
	cache               *container.GenericCache[string, CalcParam]
}

func CreateCalcDynamicParamService(calcDynamicParamDao *dao.CalcDynamicParamDao, workplaceDao *dao.WorkplaceDao, scriptDao *dao.ScriptDao) *CalcDynamicParamService {
	cache, err := container.GetOrCreateCache[string, CalcParam](container.DYNAMIC_PARAM)
	if err != nil {
		panic(err)
	}
	return &CalcDynamicParamService{calcDynamicParamDao: calcDynamicParamDao, workplaceDao: workplaceDao, scriptDao: scriptDao, cache: cache}
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

	calcParam := CalcParam{}

	for _, param := range paramList {
		switch param.Type {
		case entity.DYNAMIC_DIMENSION_STORAGE_FIELDS:
			calcParam.DimensionStorageFields = service.buildDimensionStorageField(param)
		case entity.DYNAMIC_DIMENSION_ORIGINAL_FIELDS:
			calcParam.OriginalField = service.buildOriginalField(param)
		case entity.DYNAMIC_DIMENSION_AGGREGATE_FIELDS:
			calcParam.AggregateField = service.buildAggregateField(param)
		case entity.DYNAMIC_CALC_NODES:
			calcParam.CalcNodeList = service.buildCalcNodeList(param)
		case entity.DYNAMIC_CALC_PARAMS:
			calcParam.CalcOtherParam = service.buildCalcOtherParam(param)
		}
	}
	return &calcParam
}

func (service CalcDynamicParamService) buildDimensionStorageField(param entity.CalcDynamicParamEntity) *[]DimensionStorageField {
	array, err := json_util.Parse2JsonArray(param.Content)
	if err != nil {
		// todo 所有panic检查是否可以做处理
		panic(err)
	}

	fields := make([]DimensionStorageField, 0)
	for i := 0; i < len(array.MustArray()); i++ {
		data := array.GetIndex(i)
		dimensionStorageField := DimensionStorageField{
			SinkType:  my_const.SinkType(data.Get(entity.SINK_TYPE).MustString()),
			tableName: data.Get(entity.TABLE_NAME).MustString(),
		}

		fieldColumnArray, hasValue := data.CheckGet(entity.FIELD_COLUMN_LIST)
		// 解析字段映射
		if hasValue {
			fieldName2ColumnName := make(map[string]string)

			for j := 0; j < len(fieldColumnArray.MustArray()); j++ {
				fieldMapping := fieldColumnArray.GetIndex(j)
				fieldName2ColumnName[fieldMapping.Get(entity.FIELD_NAME).MustString()] = fieldMapping.Get(entity.COLUMN_NAME).MustString()
			}
			dimensionStorageField.FieldName2ColumnName = fieldName2ColumnName
		}

		fields = append(fields, dimensionStorageField)
	}

	return &fields
}

func (service CalcDynamicParamService) buildOriginalField(param entity.CalcDynamicParamEntity) *OriginalField {
	array, err := json_util.Parse2JsonArray(param.Content)
	if err != nil {
		panic(err)
	}

	originalField := OriginalField{
		FieldSet:                    mapset.NewSet[string](),
		OriginalFieldName2FieldName: make(map[string]string),
	}

	for i := 0; i < len(array.MustArray()); i++ {
		field := array.GetIndex(i)
		originalField.FieldSet.Add(field.Get(entity.FIELD_NAME).MustString())
		columnNameField, hasValue := field.CheckGet(entity.COLUMN_NAME)
		if hasValue {
			originalField.OriginalFieldName2FieldName[columnNameField.MustString()] = field.Get(entity.FIELD_NAME).MustString()
		}
	}

	return &originalField
}

func (service CalcDynamicParamService) buildAggregateField(param entity.CalcDynamicParamEntity) *AggregateField {
	array, err := json_util.Parse2JsonArray(param.Content)
	if err != nil {
		panic(err)
	}

	aggregateField := AggregateField{
		FieldSet: mapset.NewSet[string](),
	}
	for i := 0; i < len(array.MustArray()); i++ {
		field := array.GetIndex(i)
		aggregateField.FieldSet.Add(field.Get(entity.FIELD_NAME).MustString())
	}
	return &aggregateField
}

func (service CalcDynamicParamService) buildCalcNodeList(param entity.CalcDynamicParamEntity) *CalcNodeList {
	json, err := json_util.Parse2Json(param.Content)
	if err != nil {
		panic(err)
	}

	nodeNames := strings.Split(json.Get(entity.NODE_NAMES).MustString(), ",")

	scripts := service.scriptDao.FindByNameWithMaxVersion(nodeNames)

	scriptName2Entity := make(map[string]entity.ScriptEntity)
	for _, scriptEntity := range scripts {
		scriptName2Entity[scriptEntity.Name] = scriptEntity
	}

	calcNodes := make([]CalcNode, 0)
	for _, nodeName := range nodeNames {
		scriptEntity := scriptName2Entity[nodeName]

		scriptType := entity.ScriptType(scriptEntity.Type)

		node := CalcNode{
			NodeName: nodeName,
			NodeType: scriptType,
			Script:   scriptEntity.Content,
		}
		calcNodes = append(calcNodes, node)
	}

	return &CalcNodeList{&calcNodes}
}

func (service CalcDynamicParamService) buildCalcOtherParam(param entity.CalcDynamicParamEntity) *CalcOtherParam {
	paramMap, err := json_util.Parse2Json(param.Content)
	if err != nil {
		panic(err)
	}
	return &CalcOtherParam{
		params: paramMap,
	}
}
