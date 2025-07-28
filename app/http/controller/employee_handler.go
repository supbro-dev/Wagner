/*
* @Author: supbro
* @Date:   2025/7/18 15:08
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/18 15:08
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type EmployeeHandler struct {
}

func (e EmployeeHandler) FindByName(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	workGroupCode := c.Query("workGroupCode")
	if workGroupCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workGroupCode"))
		return
	}

	employeeName := c.Query("employeeName")
	if employeeName == "" {
		response.ReturnError(c, business_error.ParamIsNil("employeeName"))
		return
	}

	employee := service.DomainHolder.EmployeeSnapshotService.FindByInfo(employeeName, workGroupCode, workplaceCode)

	employeeVo := e.convertDomain2Vo(employee)
	response.ReturnSuccessJson(c, employeeVo)
}

func (e EmployeeHandler) FindByWorkGroupCode(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	workGroupCode := c.Query("workGroupCode")
	if workGroupCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workGroupCode"))
		return
	}

	domainList := service.DomainHolder.EmployeeSnapshotService.FindByWorkGroupCode(workGroupCode, workplaceCode)

	voList := make([]*vo.EmployeeVO, 0)
	for _, emploiyee := range domainList {
		voList = append(voList, e.convertDomain2Vo(emploiyee))
	}

	response.ReturnSuccessJson(c, voList)
}

func (e EmployeeHandler) convertDomain2Vo(employee *domain.EmployeeSnapshot) *vo.EmployeeVO {
	employeeVo := vo.EmployeeVO{
		Name:          employee.Name,
		Number:        employee.Number,
		WorkplaceCode: employee.WorkplaceCode,
		WorkGroupCode: employee.WorkGroupCode,
	}

	return &employeeVo
}
