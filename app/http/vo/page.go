/*
* @Author: supbro
* @Date:   2025/6/23 22:44
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/23 22:44
 */
package vo

type Page struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	Total       int `json:"total"`
}
