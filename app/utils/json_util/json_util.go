/*
* @Author: supbro
* @Date:   2025/6/3 10:49
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 10:49
 */
package json_util

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
)

func Parse2JsonArray(data string) (*simplejson.Json, error) {
	array, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return nil, err
	}
	return array, nil
}

func Parse2Json(data string) (*simplejson.Json, error) {
	json, err := simplejson.NewJson([]byte(data))

	if err != nil {
		return nil, err
	}
	return json, nil
}

func Parse2Map(data string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Parse2Object[T any](data string, obj *T) error {
	// 将字符串转为字节切片
	jsonBytes := []byte(data)

	// 反序列化
	err := json.Unmarshal(jsonBytes, obj)
	if err != nil {
		return err
	}
	return nil
}
