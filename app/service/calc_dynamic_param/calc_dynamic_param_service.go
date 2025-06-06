/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package calc_dynamic_param

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/jinzhu/copier"
	"strings"
	"wagner/app/global/container"
	"wagner/app/global/my_const"
	"wagner/app/utils/json_util"
	"wagner/app/utils/script_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

// 人效计算参数
type CalcParam struct {
	InjectSource   InjectSource
	SinkStorages   []SinkStorage
	CalcNodeList   CalcNodeList
	CalcOtherParam CalcOtherParam
}

// 不同维度的存储属性
type SinkStorage struct {
	// 存储类型
	SinkType my_const.SinkType
	// 表名
	tableName string
	// 如果是聚合场景，聚合字段
	AggregateFields []string
	// 属性名转化成字段名
	FieldName2ColumnName map[string]string
}

// action动态字段配置,以及原始字段名转化
type InjectSource struct {
	// 需要关注的原始属性名，其他的直接丢弃
	FieldSet mapset.Set[string]
	// 把原始字段名转化成Action的属性名（如需要）
	OriginalFieldName2FieldName map[string]string
}

type CalcNodeList struct {
	List *[]CalcNode
}

// 计算节点
type CalcNode struct {
	NodeName string
	NodeType script_util.ScriptType
	// 计算脚本
	Script string
}

// 其他各类参数
type CalcOtherParam struct {
	Attendance  AttendanceParam
	HourSummary HourSummaryParam
	Work        WorkParam
}

type AttendanceParam struct {
	// 考勤缺卡惩罚时长（H）
	AttendanceAbsencePenaltyHour int
}

type WorkLoadAggregateType string

var (
	AggregateEndHour    WorkLoadAggregateType = "end"        // 物品数量记录到结束小时
	AggregateProportion WorkLoadAggregateType = "proportion" // 物品数量按比例分摊
)

type HourSummaryParam struct {
	WorkLoadAggregateType WorkLoadAggregateType
}
type CalcDynamicParamService struct {
	calcDynamicParamDao *dao.CalcDynamicParamDao
	workplaceDao        *dao.WorkplaceDao
	scriptDao           *dao.ScriptDao
	cache               *container.GenericCache[string, CalcParam]
}

type WorkParam struct {
	WorkLoadUnits mapset.Set[string]
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
		case entity.INJECT_SOURCE:
			calcParam.InjectSource = service.buildInjectSources(param)
		case entity.SINK_STORAGE:
			calcParam.SinkStorages = service.buildSinkStorages(param)
		case entity.DYNAMIC_CALC_NODE:
			calcParam.CalcNodeList = service.buildCalcNodeList(param)
		case entity.CALC_PARAM:
			calcParam.CalcOtherParam = service.buildCalcOtherParam(param)
		}
	}
	return &calcParam
}

func (service CalcDynamicParamService) buildSinkStorages(param entity.CalcDynamicParamEntity) []SinkStorage {
	array, err := json_util.Parse2JsonArray(param.Content)
	if err != nil {
		// todo 所有panic检查是否可以做处理
		panic(err)
	}

	fields := make([]SinkStorage, 0)
	for i := 0; i < len(array.MustArray()); i++ {
		data := array.GetIndex(i)
		sinkStorage := SinkStorage{
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
			sinkStorage.FieldName2ColumnName = fieldName2ColumnName
		}

		// 如果是SUMMARY类型，添加聚合字段
		aggregateFieldArray, hasValue := data.CheckGet(entity.AGGREGATE_FILEDS)
		if hasValue {
			aggregateFields := make([]string, 0)
			for i := 0; i < len(aggregateFieldArray.MustArray()); i++ {
				aggregateFields = append(aggregateFields, aggregateFieldArray.GetIndex(i).MustString())
			}
			sinkStorage.AggregateFields = aggregateFields
		}

		fields = append(fields, sinkStorage)
	}

	return fields
}

func (service CalcDynamicParamService) buildInjectSources(param entity.CalcDynamicParamEntity) InjectSource {
	array, err := json_util.Parse2JsonArray(param.Content)
	if err != nil {
		panic(err)
	}

	originalField := InjectSource{
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

	return originalField
}

func (service CalcDynamicParamService) buildCalcNodeList(param entity.CalcDynamicParamEntity) CalcNodeList {
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

		scriptType := script_util.ScriptType(scriptEntity.Type)

		node := CalcNode{
			NodeName: nodeName,
			NodeType: scriptType,
			Script:   scriptEntity.Content,
		}
		calcNodes = append(calcNodes, node)
	}

	return CalcNodeList{&calcNodes}
}

var defaultCalcOtherParam = CalcOtherParam{
	Attendance: AttendanceParam{
		// 默认惩罚8小时
		AttendanceAbsencePenaltyHour: 8,
	},
	HourSummary: HourSummaryParam{
		// 默认聚合到结束的那个小时里
		WorkLoadAggregateType: AggregateEndHour,
	},
	Work: WorkParam{
		WorkLoadUnits: mapset.NewSet[string]("itemNum", "skuNum", "packageNum"),
	},
}

func (service CalcDynamicParamService) buildCalcOtherParam(param entity.CalcDynamicParamEntity) CalcOtherParam {
	otherParam := CalcOtherParam{}
	copyError := copier.Copy(&otherParam, &defaultCalcOtherParam)
	if copyError != nil {
		panic(copyError)
	}
	err := json_util.Parse2Object[CalcOtherParam](param.Content, &otherParam)
	if err != nil {
		panic(err)
	}
	return otherParam
}
