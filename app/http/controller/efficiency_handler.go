/*
* @Author: supbro
* @Date:   2025/6/11 09:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:10
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/service"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/response"
)

type EfficiencyHandler struct {
}

func (p EfficiencyHandler) EmployeeStatus(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	operateDayStr := c.Query("operateDay")

	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	if operateDayStr == "" {
		response.ReturnError(c, business_error.ParamIsNil("operateDay"))
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("operateDay"))
		return
	}

	employeeStatusVO := service.Holder.EfficiencyService.QueryEmployeeStatus(workplaceCode, operateDay)

	response.ReturnSuccessJson(c, employeeStatusVO)
}

func (p EfficiencyHandler) EmployeeEfficiency(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	workplaceCode := c.Query("workplaceCode")
	aggregateDimension := c.Query("aggregateDimension")
	isCrossPosition := c.Query("isCrossPosition")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	currentPage := c.Query("currentPage")
	pageSize := c.Query("pageSize")

	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	if aggregateDimension == "" {
		response.ReturnError(c, business_error.ParamIsNil("aggregateDimension"))
		return
	}

	if startDateStr == "" || endDateStr == "" {
		response.ReturnError(c, business_error.ParamIsNil("startDate", "endDate"))
		return
	}

	if currentPage == "" || pageSize == "" {
		response.ReturnError(c, business_error.ParamIsNil("currentPage", "pageSize"))
		return
	}

	startDate, err := datetime_util.ParseDate(startDateStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("startDate"))
		return
	}

	endDate, err := datetime_util.ParseDate(endDateStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("endDate"))
		return
	}

	calcParam, e := service.DomainHolder.CalcDynamicParamService.FindParamsByWorkplace(workplaceCode)
	if e != nil {
		response.ReturnError(c, e)
		return
	}
	if calcParam == nil {
		response.ReturnError(c, business_error.CannotFindCalcParamByWorkplace(workplaceCode))
		return
	}

	workLoadUnits := calcParam.CalcOtherParam.Work.WorkLoadUnits

	currentPageInt, _ := strconv.Atoi(currentPage)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	efficiencyVO := service.Holder.EfficiencyService.EmployeeEfficiency(workplaceCode, employeeNumber, []*time.Time{&startDate, &endDate},
		domain.AggregateDimension(aggregateDimension), domain.IsCrossPosition(isCrossPosition), workLoadUnits, currentPageInt, pageSizeInt)

	response.ReturnSuccessJson(c, efficiencyVO)
}

func (p EfficiencyHandler) ComputeEmployee(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	operateDayStr := c.Query("operateDay")

	if employeeNumber == "" {
		response.ReturnError(c, business_error.ParamIsNil("employeeNumber"))
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("operateDay"))
		return
	}

	if isSuccess, err := service.Holder.EfficiencyComputeService.ComputeEmployee(employeeNumber, operateDay); err != nil {
		response.ReturnError(c, err)
	} else {
		response.ReturnSuccessJson(c, isSuccess)
	}
}

func (p EfficiencyHandler) ComputeWorkplace(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	operateDayStr := c.Query("operateDay")

	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("operateDay"))
		return
	}

	if isSuccess, err := service.Holder.EfficiencyComputeService.ComputeWorkplace(workplaceCode, operateDay); err != nil {
		response.ReturnError(c, err)
	} else {
		response.ReturnSuccessJson(c, isSuccess)
	}
}

func (p EfficiencyHandler) TimeOnTask(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	operateDayStr := c.Query("operateDay")

	if employeeNumber == "" {
		response.ReturnError(c, business_error.ParamIsNil("employeeNumber"))
		return
	}

	if operateDayStr == "" {
		response.ReturnError(c, business_error.ParamIsNil("operateDay"))
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("operateDay"))
		return
	}

	if timeOnTaskVO, err := service.Holder.EfficiencyComputeService.TimeOnTask(employeeNumber, operateDay); err != nil {
		response.ReturnError(c, err)
	} else {
		response.ReturnSuccessJson(c, timeOnTaskVO)
	}

}

func (p EfficiencyHandler) WorkplaceEfficiency(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	isCrossPosition := c.Query("isCrossPosition")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if workplaceCode == "" {
		response.ReturnError(c, business_error.ParamIsNil("workplaceCode"))
		return
	}

	if startDateStr == "" || endDateStr == "" {
		response.ReturnError(c, business_error.ParamIsNil("startDate", "endDate"))
		return
	}

	startDate, err := datetime_util.ParseDate(startDateStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("startDate"))
		return
	}

	endDate, err := datetime_util.ParseDate(endDateStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("endDate"))
		return
	}

	calcParam, e := service.DomainHolder.CalcDynamicParamService.FindParamsByWorkplace(workplaceCode)
	if e != nil {
		response.ReturnError(c, e)
		return
	}
	if calcParam == nil {
		response.ReturnError(c, business_error.CannotFindCalcParamByWorkplace(workplaceCode))
		return
	}

	workLoadUnits := calcParam.CalcOtherParam.Work.WorkLoadUnits

	workplace := service.DomainHolder.WorkplaceService.FindByCode(workplaceCode)

	processPositions := service.DomainHolder.ProcessService.FindProcessPositionList(workplace)

	workplaceEfficiencyVO := service.Holder.EfficiencyService.WorkplaceEfficiency(workplace, []*time.Time{&startDate, &endDate}, domain.IsCrossPosition(isCrossPosition), workLoadUnits, processPositions)

	response.ReturnSuccessJson(c, workplaceEfficiencyVO)
}
