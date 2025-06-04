/*
* @Author: supbro
* @Date:   2025/6/3 10:49
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 10:49
 */
package json_util

import (
	"encoding/json"
	"fmt"
)

// 泛型方法：将JSON解析为Map对象
func parse2Map[T any](data string) (map[string]T, error) {
	// 创建目标map的指针
	var result map[string]T
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return result, nil
}

// 泛型方法：将JSON解析为Map对象数组
func parse2MapSlice[T any](data string) ([]map[string]T, error) {
	// 创建目标slice的指针
	var result []map[string]T
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return result, nil
}
