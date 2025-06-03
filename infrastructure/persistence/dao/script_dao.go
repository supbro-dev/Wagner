/*
* @Author: supbro
* @Date:   2025/6/3 11:21
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 11:21
 */
package dao

import (
	"gorm.io/gorm"
	"wagner/infrastructure/persistence/entity"
)

type ScriptDao struct {
	db *gorm.DB
}

func (d ScriptDao) FindByNameWithMaxVersion(names []string) *[]entity.ScriptEntity {
	scripts := make([]entity.ScriptEntity, 0)
	d.db.Exec("SELECT name, type, desc, content FROM"+
		" (SELECT  name, type, desc, content ,ROW_NUMBER() OVER (PARTITION BY name ORDER BY version DESC) AS rn  FROM script  WHERE name IN ? ) ranked WHERE rn = 1", names).Find(&scripts)

	return &scripts
}

func CreateScriptDao(client *gorm.DB) *ScriptDao {
	return &ScriptDao{client}
}
