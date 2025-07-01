/*
* @Author: supbro
* @Date:   2025/7/1 13:03
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 13:03
 */
package dao

import (
	"gorm.io/gorm"
)

type ProcessImplementDao struct {
	db *gorm.DB
}

func CreateProcessImplementDao(client *gorm.DB) *ProcessImplementDao {
	return &ProcessImplementDao{client}
}
