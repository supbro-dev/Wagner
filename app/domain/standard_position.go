/*
* @Author: supbro
* @Date:   2025/6/5 10:30
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 10:50
 */
package domain

type StandardPosition struct {
	Name string
	Code string
	// 路径
	Path []StandardPosition
}
