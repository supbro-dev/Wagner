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
	DYNAMIC_DIMENSION_STORAGE_FIELDS   = "DYNAMIC_DIMENSION_STORAGE_FIELDS"
	DYNAMIC_DIMENSION_ORIGINAL_FIELDS  = "DYNAMIC_DIMENSION_ORIGINAL_FIELDS"
	DYNAMIC_DIMENSION_AGGREGATE_FIELDS = "DYNAMIC_DIMENSION_AGGREGATE_FIELDS"
	DYNAMIC_CALC_NODES                 = "DYNAMIC_CALC_NODES"
	DYNAMIC_CALC_PARAMS                = "DYNAMIC_CALC_PARAMS"
)

// json格式解析
var (
	// DYNAMIC_DIMENSION_STORAGE_FIELDS
	SINK_TYPE         = "SinkType"
	FIELD_COLUMN_LIST = "fieldColumnList"
	FIELD_NAME        = "fieldName"
	COLUMN_NAME       = "columnName"
)
