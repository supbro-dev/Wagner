/*
* @Author: supbro
* @Date:   2025/7/2 10:31
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 10:31
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wagner/app/global/business_error"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
	"wagner/infrastructure/persistence/entity"
)

type ProcessHandler struct {
}

func (p ProcessHandler) Implementation(c *gin.Context) {
	targetType := c.Query("targetType")
	workplaceCode := c.Query("workplaceCode")
	industryCode := c.Query("industryCode")
	subIndustryCode := c.Query("subIndustryCode")
	currentPage := c.Query("currentPage")
	pageSize := c.Query("pageSize")

	if targetType == "" {
		response.ReturnError(c, business_error.ProcessTargetTypeError())
		return
	}

	currentPageInt, _ := strconv.Atoi(currentPage)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	list, total := service.DomainHolder.ProcessService.FindProcessImplementationListByPage(targetType, workplaceCode, industryCode, subIndustryCode, currentPageInt, pageSizeInt)

	implVoList := make([]*vo.ProcessImplementationVO, 0)
	for _, impl := range list {
		var targetTypeDesc string
		switch impl.TargetType {
		case entity.Workplace:
			targetTypeDesc = "工作点"
		case entity.Industry:
			targetTypeDesc = "行业"
		case entity.SubIndustry:
			targetTypeDesc = "子行业"
		}

		var statusDesc string
		switch impl.Status {
		case entity.Offline:
			statusDesc = "下线"
		case entity.Online:
			statusDesc = "上线"
		}
		implVoList = append(implVoList, &vo.ProcessImplementationVO{
			strconv.FormatInt(impl.Id, 10), impl.Id, impl.Name, targetTypeDesc, impl.TargetName, string(impl.Status), statusDesc,
		})
	}

	response.ReturnSuccessJson(c, vo.ProcessImplementationPageVO{
		TableDataList: implVoList,
		Page: &vo.Page{
			CurrentPage: currentPageInt,
			PageSize:    pageSizeInt,
			Total:       total,
		},
	})
}
