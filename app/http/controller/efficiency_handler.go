/*
* @Author: supbro
* @Date:   2025/6/11 09:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:10
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"wagner/app/domain"
	"wagner/app/global/my_error"
	"wagner/app/service"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/response"
)

type EfficiencyHandler struct {
}

func (p EfficiencyHandler) EmployeeEfficiency(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	workplaceCode := c.Query("workplaceCode")
	aggregateDimension := c.Query("aggregateDimension")
	isCrossPosition := c.Query("isCrossPosition")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

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

	service.Holder.EfficiencyService.EmployeeEfficiency(workplaceCode, employeeNumber, []*time.Time{&startDate, &endDate},
		domain.AggregateDimension(aggregateDimension), domain.IsCrossPosition(isCrossPosition), workLoadUnits)

}
