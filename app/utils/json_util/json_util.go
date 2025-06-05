/*
* @Author: supbro
* @Date:   2025/6/3 10:49
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 10:49
 */
package json_util

import (
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
