/*
* @Author: supbro
* @Date:   2025/6/2 10:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:47
 */
package dao

import "gorm.io/gorm"

type CalcDynamicParamDao struct {
	db *gorm.DB
}

func CreateCalcDynamicParamDao(client *gorm.DB) *CalcDynamicParamDao {
	return &CalcDynamicParamDao{client}
}
