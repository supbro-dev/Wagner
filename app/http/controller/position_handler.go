/*
* @Author: supbro
* @Date:   2025/7/9 11:26
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/9 11:26
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type PositionHandler struct {
}

func (p PositionHandler) FindAll(c *gin.Context) {
	positions := service.DomainHolder.PositionService.FindAll()

	selectList := make([]vo.SelectVO, 0)
	for _, position := range positions {
		selectList = append(selectList, vo.SelectVO{
			Value: position.Code,
			Label: position.Name,
		})
	}
	response.ReturnSuccessJson(c, selectList)
}
