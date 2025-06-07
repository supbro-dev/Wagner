/*
* @Author: supbro
* @Date:   2025/6/2 10:44
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:44
 */
package entity

type CalcDynamicParamEntity struct {
	BaseEntity
	Type            string `gorm:"column:type" json:"type"`
	IndustryCode    string `gorm:"column:industry_code" json:"industry_code"`
	SubIndustryCode string `gorm:"column:sub_industry_code" json:"sub_industry_code"`
	Content         string `gorm:"column:content" json:"content"`
}

func (u *CalcDynamicParamEntity) TableName() string {
	return "calc_dynamic_param" // 自定义表名
}

const (
	INJECT_SOURCE     = "INJECT_SOURCE"
	SINK_STORAGE      = "SINK_STORAGE"
	DYNAMIC_CALC_NODE = "DYNAMIC_CALC_NODE"
	CALC_PARAM        = "CALC_PARAM"
)

// json格式解析
var (
	// INJECT_SOURCE
	SINK_TYPE         = "sinkType"
	TABLE_NAME        = "tableName"
	FIELD_COLUMN_LIST = "fieldColumnList"
	FIELD_NAME        = "fieldName"
	COLUMN_NAME       = "columnName"
	NODE_NAMES        = "nodeNames"
	AGGREGATE_FILEDS  = "aggregateFields"
	WORK_LOAD         = "workLoad"
)
