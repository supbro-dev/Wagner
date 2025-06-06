package controller

import (
	"github.com/gin-gonic/gin"
	"wagner/app/global/my_error"
	"wagner/app/service"
	"wagner/app/utils/datetime_util"
	"wagner/app/utils/response"
)

type EfficiencyComputeHandler struct {
}

func (p EfficiencyComputeHandler) Invoke(c *gin.Context) {
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
