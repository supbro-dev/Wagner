/*
* @Author: supbro
* @Date:   2025/6/4 10:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/4 10:15
 */
package common

import (
	"database/sql/driver"
	"errors"
	"sync"
	"time"
)

// 用于保存gorm schema的全局缓存
var SchemaCache = &sync.Map{}

// Date 自定义日期类型
type Date time.Time

// 实现 Scanner 接口 - 从数据库读取
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return errors.New("无法转换值为时间类型")
	}

	*d = Date(t)
	return nil
}

// 实现 Valuer 接口 - 写入数据库
func (d Date) Value() (driver.Value, error) {
	// 返回日期部分 (YYYY-MM-DD)
	return time.Time(d).Format("2006-01-02"), nil
}

// 实现 GormDataType 接口
func (Date) GormDataType() string {
	return "date"
}

// 实现 MarshalJSON 接口
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}
