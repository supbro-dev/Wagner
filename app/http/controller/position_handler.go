/*
* @Author: supbro
* @Date:   2025/7/9 11:26
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/9 11:26
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/global/business_error"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type PositionHandler struct {
}

func (p PositionHandler) FindAll(c *gin.Context) {
	industryCode := c.Query("industryCode")
	if industryCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("industryCode"))
		return
	}

	subIndustryCode := c.Query("subIndustryCode")

	positions := service.DomainHolder.PositionService.FindAll(industryCode, subIndustryCode)

	selectList := make([]vo.SelectVO, 0)
	for _, position := range positions {
		selectList = append(selectList, vo.SelectVO{
			Value: position.Code,
			Label: position.Name,
		})
	}
	response.ReturnSuccessJson(c, selectList)
}
