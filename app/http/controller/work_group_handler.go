/*
* @Author: supbro
* @Date:   2025/7/17 20:19
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/17 20:19
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"wagner/app/global/business_error"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type WorkGroupHandler struct {
}

func (p WorkGroupHandler) FindByCode(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsWrong("workplaceCode"))
		return
	}

	workGroupCode := c.Query("workGroupCode")
	if workGroupCode == "" {
		response.ReturnError(c, business_error.ParamIsWrong("workGroupCode"))
		return
	}
	workGroup := service.DomainHolder.WorkGroupService.FindByCode(workGroupCode, workplaceCode)

	workGroupVO := vo.WorkGroupVO{}
	copier.Copy(&workGroupVO, &workGroup)

	response.ReturnSuccessJson(c, workGroupVO)
}
