/*
* @Author: supbro
* @Date:   2025/7/5 11:41
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/5 11:41
 */
package qo

type ProcessImplementationSaveQo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	TargetType string `json:"targetType"`
	TargetCode string `json:"targetCode"`
}
