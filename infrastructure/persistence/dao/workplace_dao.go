/*
* @Author: supbro
* @Date:   2025/6/2 11:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 11:29
 */
package dao

import "gorm.io/gorm"

type WorkplaceDao struct {
	db *gorm.DB
}

func CreateWorkplaceDao(client *gorm.DB) *WorkplaceDao {
	return &WorkplaceDao{client}
}
