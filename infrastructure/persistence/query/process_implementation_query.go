/*
* @Author: supbro
* @Date:   2025/7/2 16:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 16:12
 */
package query

type ProcessImplementationQuery struct {
	TargetType  string
	TargetCode  string
	Code        string
	CurrentPage int
	PageSize    int
}
