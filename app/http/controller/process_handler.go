/*
* @Author: supbro
* @Date:   2025/7/2 10:31
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 10:31
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"strconv"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/http/qo"
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
		implVoList = append(implVoList, p.convertDomain2Vo(impl))
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

func (p ProcessHandler) convertDomain2Vo(d *domain.ProcessImplementation) *vo.ProcessImplementationVO {
	var targetTypeDesc string
	switch d.TargetType {
	case entity.Workplace:
		targetTypeDesc = "工作点"
	case entity.Industry:
		targetTypeDesc = "行业"
	case entity.SubIndustry:
		targetTypeDesc = "子行业"
	}

	var statusDesc string
	switch d.Status {
	case entity.Offline:
		statusDesc = "下线"
	case entity.Online:
		statusDesc = "上线"
	}
	v := vo.ProcessImplementationVO{}
	copier.Copy(&v, &d)
	v.Key = strconv.FormatInt(d.Id, 10)
	v.TargetTypeDesc = targetTypeDesc
	v.StatusDesc = statusDesc

	return &v
}

func (p ProcessHandler) Save(c *gin.Context) {
	var req qo.ProcessImplementationSaveQo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, business_error.SubmitDataIsWrong(err))
		return
	}

	d := domain.ProcessImplementation{}
	copier.Copy(&d, &req)
	if i, err := strconv.Atoi(req.Id); err == nil {
		d.Id = int64(i)
	}
	id, businessError := service.DomainHolder.ProcessService.Save(&d)

	if businessError != nil {
		response.ReturnError(c, businessError)
	} else {
		response.ReturnSuccessJson(c, id)
	}
}

func (p ProcessHandler) GetImplementationById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("id"))
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(int64(id))

	response.ReturnSuccessJson(c, p.convertDomain2Vo(impl))
}

func (p ProcessHandler) GetProcessPositionTree(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("id"))
	}

	tree := service.DomainHolder.ProcessService.GetProcessPositionTree(int64(id))

	treeNodeVo := p.iterateConvert2Vo(tree)

	response.ReturnSuccessJson(c, treeNodeVo)
}

func (p ProcessHandler) iterateConvert2Vo(node *domain.ProcessPositionTreeNode) *vo.ProcessPositionTreeNodeVo {
	v := vo.ProcessPositionTreeNodeVo{
		Title:    node.Name,
		Key:      node.Code,
		Type:     string(node.Type),
		Children: make([]*vo.ProcessPositionTreeNodeVo, 0),
	}

	if node.Children != nil && len(node.Children) > 0 {
		children := make([]*vo.ProcessPositionTreeNodeVo, 0)
		for _, c := range node.Children {
			childVo := p.iterateConvert2Vo(c)
			children = append(children, childVo)
		}
		v.Children = children
	}

	return &v
}
