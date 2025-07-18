/*
* @Author: supbro
* @Date:   2025/7/17 20:09
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/17 20:09
 */
package vo

type WorkplaceVo struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	RegionCode      string `json:"regionCode"`
	IndustryCode    string `json:"industryCode"`
	SubIndustryCode string `json:"subIndustryCode"`
	Desc            string `json:"desc"`
}
