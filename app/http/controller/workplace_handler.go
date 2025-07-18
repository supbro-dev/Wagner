/*
* @Author: supbro
* @Date:   2025/6/10 20:26
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/10 20:26
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

func (p WorkplaceHandler) FindAllIndustry(c *gin.Context) {
	industries := service.DomainHolder.WorkplaceService.FindAllIndustry()

	selectList := make([]vo.SelectVO, 0)
	for _, s := range industries {
		selectList = append(selectList, vo.SelectVO{
			Value: s,
			Label: s,
		})
	}
	response.ReturnSuccessJson(c, selectList)
}

func (p WorkplaceHandler) FindAllSubIndustry(c *gin.Context) {
	subIndustries := service.DomainHolder.WorkplaceService.FindAllSubIndustry()

	selectList := make([]vo.SelectVO, 0)
	for _, s := range subIndustries {
		selectList = append(selectList, vo.SelectVO{
			Value: s,
			Label: s,
		})
	}
	response.ReturnSuccessJson(c, selectList)
}

func (p WorkplaceHandler) FindWorkplaceByCode(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	workplace := service.DomainHolder.WorkplaceService.FindByCode(workplaceCode)

	workplaceVo := vo.WorkplaceVo{}
	copier.Copy(&workplaceVo, &workplace)

	response.ReturnSuccessJson(c, workplaceVo)
}
