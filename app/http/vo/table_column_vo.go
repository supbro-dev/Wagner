/*
* @Author: supbro
* @Date:   2025/6/11 14:17
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 14:17
 */
package vo

type TableColumnVO struct {
	Title     string      `json:"title"`
	DataIndex interface{} `json:"dataIndex"`
	Key       string      `json:"key"`
}
