package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"wagner/app/global/my_const"
	"wagner/app/global/my_error"
	"wagner/app/service"
	"wagner/app/utils/response"
)

// todo 不要有ppr关键字
type PprComputeHandler struct {
}

func (p PprComputeHandler) Invoke(c *gin.Context) {
	employeeNumber := c.Query("employeeNumber")
	operateDayStr := c.Query("operateDay")

	if employeeNumber == "" {
		response.ErrorSystem(c, my_error.ParamNilCode, my_error.ParamNilMsg, operateDayStr)
		return
	}

	operateDay, err := time.Parse(my_const.DateLayout, operateDayStr)
	if err != nil {
		response.ErrorSystem(c, my_error.ParamErrorCode, my_error.ParamErrorMsg, operateDayStr)
		return
	}

	service.Holder.PprComputeService.ComputeEmployee(employeeNumber, operateDay)

	response.ReturnSuccessEmptyJson(c)
}
