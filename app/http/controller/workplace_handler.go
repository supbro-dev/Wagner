/*
* @Author: supbro
* @Date:   2025/6/10 20:26
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 20:26
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type WorkplaceHandler struct {
}

func (p WorkplaceHandler) FindAll(c *gin.Context) {
	workplaces := service.DomainHolder.WorkplaceService.FindAll()

	selectList := make([]vo.SelectVO, 0)
	for _, workplace := range workplaces {
		selectList = append(selectList, vo.SelectVO{
			Value: workplace.Code,
			Label: workplace.Name,
		})
	}
	response.ReturnSuccessJson(c, selectList)
}
