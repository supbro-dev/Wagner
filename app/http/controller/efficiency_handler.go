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
	"wagner/app/global/my_error"
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
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "workplaceCode")
		return
	}

	if operateDayStr == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "operateDay")
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, operateDayStr)
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
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "workplaceCode")
		return
	}

	if aggregateDimension == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "aggregateDimension")
		return
	}

	if startDateStr == "" || endDateStr == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "startDate or endDate")
		return
	}

	if currentPage == "" || pageSize == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "page")
		return
	}

	startDate, err := datetime_util.ParseDate(startDateStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, startDateStr)
		return
	}

	endDate, err := datetime_util.ParseDate(endDateStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, endDateStr)
		return
	}

	calcParam := service.DomainHolder.CalcDynamicParamService.FindParamsByWorkplace(workplaceCode)
	if calcParam == nil {
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
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, operateDayStr)
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, operateDayStr)
		return
	}

	service.Holder.EfficiencyComputeService.ComputeEmployee(employeeNumber, operateDay)

	response.ReturnSuccessEmptyJson(c)
}

func (p EfficiencyHandler) TimeOnTask(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	operateDayStr := c.Query("operateDay")

	if employeeNumber == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "employeeNumber")
		return
	}

	if operateDayStr == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "operateDay")
		return
	}

	operateDay, err := datetime_util.ParseDate(operateDayStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, operateDayStr)
		return
	}

	timeOnTaskVO := service.Holder.EfficiencyComputeService.TimeOnTask(employeeNumber, operateDay)

	response.ReturnSuccessJson(c, timeOnTaskVO)
}

func (p EfficiencyHandler) WorkplaceEfficiency(c *gin.Context) {
	workplaceCode := c.Query("workplaceCode")
	isCrossPosition := c.Query("isCrossPosition")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if workplaceCode == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "workplaceCode")
		return
	}

	if startDateStr == "" || endDateStr == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, "startDate or endDate")
		return
	}

	startDate, err := datetime_util.ParseDate(startDateStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, startDateStr)
		return
	}

	endDate, err := datetime_util.ParseDate(endDateStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, endDateStr)
		return
	}

	calcParam := service.DomainHolder.CalcDynamicParamService.FindParamsByWorkplace(workplaceCode)
	if calcParam == nil {
		return
	}

	workLoadUnits := calcParam.CalcOtherParam.Work.WorkLoadUnits

	workplace := service.DomainHolder.WorkplaceService.FindByCode(workplaceCode)

	standardPositions := service.DomainHolder.StandardPositionService.FindStandardPositionListByIndustry(workplace.IndustryCode, workplace.SubIndustryCode)

	workplaceEfficiencyVO := service.Holder.EfficiencyService.WorkplaceEfficiency(workplace, []*time.Time{&startDate, &endDate}, domain.IsCrossPosition(isCrossPosition), workLoadUnits, standardPositions)

	response.ReturnSuccessJson(c, workplaceEfficiencyVO)
}
