/*
* @Author: supbro
* @Date:   2025/7/18 15:08
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/18 15:08
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/global/business_error"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/utils/response"
)

type EmployeeHandler struct {
}

func (e EmployeeHandler) FindByInfo(c *gin.Context) {
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

	employeeVo := vo.EmployeeVO{
		Name:          employee.Name,
		Number:        employee.Number,
		WorkplaceCode: employee.WorkplaceCode,
		WorkGroupCode: employee.WorkGroupCode,
	}

	response.ReturnSuccessJson(c, employeeVo)
}
