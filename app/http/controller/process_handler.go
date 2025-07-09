/*
* @Author: supbro
* @Date:   2025/7/2 10:31
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 10:31
 */
package controller

import (
	"fmt"
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
		return
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(int64(id))

	response.ReturnSuccessJson(c, p.convertDomain2Vo(impl))
}

func (p ProcessHandler) GetProcessPositionTree(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("id"))
		return
	}

	tree := service.DomainHolder.ProcessService.GetProcessPositionTree(int64(id))

	treeNodeVo := p.iterateConvert2Vo(tree)

	response.ReturnSuccessJson(c, treeNodeVo)
}

func (p ProcessHandler) FindProcessByParentProcessCode(c *gin.Context) {
	processCode := c.Query("processCode")
	processImplId := c.Query("processImplId")

	id, err := strconv.Atoi(processImplId)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("processImplId"))
		return
	}

	if processCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("processCode"))
		return
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(int64(id))
	version := impl.ProcessPositionRootId

	processPositionList := service.DomainHolder.ProcessService.FindProcessByParentCode(processCode, version)

	detailList := p.convertProcessDomainList2Detail(processPositionList)

	response.ReturnSuccessJson(c, detailList)
}
func (p ProcessHandler) GenerateProcessCode(c *gin.Context) {
	processName := c.Query("processName")
	processImplId := c.Query("processImplId")

	id, err := strconv.Atoi(processImplId)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("processImplId"))
		return
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(int64(id))
	version := impl.ProcessPositionRootId

	code := service.DomainHolder.ProcessService.GenerateProcessCode(processName, version)

	response.ReturnSuccessJson(c, code)
}

func (p ProcessHandler) SaveProcessPosition(c *gin.Context) {
	var req qo.ProcessPositionSaveQo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, business_error.SubmitDataIsWrong(err))
		return
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(req.ProcessImplId)

	var parentCode string
	if qo.AddLevelType(req.AddLevelType) == qo.NextLevel {
		parentCode = req.ParentProcessCode
	} else {
		sameLevelPosition := service.DomainHolder.ProcessService.FindProcessByCode(req.ParentProcessCode, impl.ProcessPositionRootId)
		parentCode = sameLevelPosition.ParentCode
	}

	d := domain.ProcessPosition{
		Name:       req.ProcessName,
		Code:       req.ProcessCode,
		ParentCode: parentCode,
		Type:       entity.ProcessPositionType(req.Type),
		Version:    int(impl.ProcessPositionRootId),
		SortIndex:  req.SortIndex,
	}
	if req.Id != 0 {
		d.Id = req.Id
	}

	if req.WorkLoadRollUp != "" {
		if workLoadRollUpBool, err := strconv.ParseBool(req.WorkLoadRollUp); err == nil {
			d.Properties = map[string]interface{}{"workLoadRollUp": workLoadRollUpBool}
		}
	}

	service.DomainHolder.ProcessService.SaveProcessPosition(&d)
	response.ReturnSuccessEmptyJson(c)
}

func (p ProcessHandler) convertProcessDomainList2Detail(processPositionList []*domain.ProcessPosition) []*vo.ProcessDetailVo {
	detailList := make([]*vo.ProcessDetailVo, 0)
	for _, process := range processPositionList {
		detail := vo.ProcessDetailVo{
			Id:          process.Id,
			ProcessName: process.Name,
			ProcessCode: process.Code,
			TypeDesc:    entity.ProcessPositionType2Desc(process.Type),
			Script:      process.Script,
		}

		if maxTimeInMinute, exists := process.Properties[entity.MaxTimeInMinuteKey]; exists {
			detail.MaxTimeInMinute = fmt.Sprintf("%v", maxTimeInMinute)
		} else {
			detail.MaxTimeInMinute = "默认"
		}
		if minIdleTimeInMinute, exists := process.Properties[entity.MinIdleTimeKey]; exists {
			detail.MinIdleTimeInMinute = fmt.Sprintf("%v", minIdleTimeInMinute)
		} else {
			detail.MinIdleTimeInMinute = "默认"
		}
		if workLoadRollUp, exists := process.Properties[entity.WorkLoadRollUpKey]; exists {
			if parseBool, err := strconv.ParseBool(fmt.Sprintf("%v", workLoadRollUp)); err == nil && parseBool {
				detail.WorkLoadRollUpDesc = "是"
			}
		}
		detailList = append(detailList, &detail)
	}

	return detailList
}

func (p ProcessHandler) iterateConvert2Vo(node *domain.ProcessPositionTreeNode) *vo.ProcessPositionTreeNodeVo {
	v := vo.ProcessPositionTreeNodeVo{
		Id:             node.Id,
		Title:          node.Name,
		Key:            node.Code,
		Type:           string(node.Type),
		ParentName:     node.ParentName,
		ParentCode:     node.ParentCode,
		SortIndex:      node.SortIndex,
		WorkLoadRollUp: node.WorkLoadRollUp,
		Children:       make([]*vo.ProcessPositionTreeNodeVo, 0),
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
