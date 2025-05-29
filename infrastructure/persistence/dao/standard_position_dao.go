package dao

import "gorm.io/gorm"

type StandardPositionDao struct {
	db *gorm.DB
}

func CreateStandardPositionDao(client *gorm.DB) *StandardPositionDao {
	return &StandardPositionDao{client}
}
