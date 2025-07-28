/*
* @Author: supbro
* @Date:   2025/7/23 22:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/23 22:12
 */
package vo

type AiDataVo struct {
	TextKey      string      `json:"text_key"`
	Data         interface{} `json:"data"`
	DefaultValue string      `json:"default_value"`
}
