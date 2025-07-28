/*
* @Author: supbro
* @Date:   2025/6/11 09:10
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 09:10
 */
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/http/qo"
	"wagner/app/http/vo"
	"wagner/app/service"
	"wagner/app/service/calc/calc_dynamic_param"
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

	var employeeNumberList []string
	if employeeNumber != "" {
		employeeNumberList = strings.Split(employeeNumber, ",")
	} else {
		employeeNumberList = make([]string, 0)
	}

	efficiencyVO := service.Holder.EfficiencyService.EmployeeEfficiency(workplaceCode, employeeNumberList, []*time.Time{&startDate, &endDate},
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

func (p EfficiencyHandler) FindCalcParamByImplementationId(c *gin.Context) {
	idStr := c.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("id"))
		return
	}

	impl := service.DomainHolder.ProcessService.GetImplementationById(int64(id))

	calcOtherParam := service.Holder.EfficiencyComputeService.GetCalcOtherParam(impl)
	if calcOtherParam == nil {
		response.ReturnSuccessEmptyJson(c)
		return
	}

	v := p.convertCalcOtherParam2Vo(calcOtherParam)
	response.ReturnSuccessJson(c, v)
}

func (p EfficiencyHandler) SaveOtherParams(c *gin.Context) {
	var req qo.ProcessOtherParamSaveQo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, business_error.SubmitDataIsWrong(err))
		return
	}

	calcOtherParam := calc_dynamic_param.CalcOtherParam{
		Attendance:  calc_dynamic_param.AttendanceParam{},
		HourSummary: calc_dynamic_param.HourSummaryParam{},
		Work:        calc_dynamic_param.WorkParam{},
	}

	if req.AttendanceAbsencePenaltyHour == "" {
		response.ReturnError(c, business_error.ParamIsNil("attendanceAbsencePenaltyHour"))
		return
	}
	if attendanceAbsencePenaltyHour, err := strconv.ParseFloat(req.AttendanceAbsencePenaltyHour, 32); err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("attendanceAbsencePenaltyHour"))
		return
	} else {
		calcOtherParam.Attendance.AttendanceAbsencePenaltyHour = float32(attendanceAbsencePenaltyHour)
	}

	if req.MaxRunUpTimeInMinute == "" {
		response.ReturnError(c, business_error.ParamIsNil("maxRunUpTimeInMinute"))
		return
	}
	if maxRunUpTimeInMinute, err := strconv.Atoi(req.MaxRunUpTimeInMinute); err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("maxRunUpTimeInMinute"))
		return
	} else {
		calcOtherParam.Attendance.MaxRunUpTimeInMinute = maxRunUpTimeInMinute
	}

	if req.WorkLoadUnits == "" {
		response.ReturnError(c, business_error.ParamIsNil("workLoadUnits"))
		return
	}
	workLoadUnits := make([]calc_dynamic_param.WorkLoadUnit, 0)
	for _, kv := range strings.Split(req.WorkLoadUnits, ",") {
		names := strings.Split(kv, ":")
		workLoadUnits = append(workLoadUnits, calc_dynamic_param.WorkLoadUnit{
			names[0], names[1],
		})
	}
	calcOtherParam.Work.WorkLoadUnits = workLoadUnits

	if req.LookBackDays == "" {
		response.ReturnError(c, business_error.ParamIsNil("lookBackDays"))
		return
	}
	if lookBackDays, err := strconv.Atoi(req.LookBackDays); err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("lookBackDays"))
		return
	} else {
		calcOtherParam.Work.LookBackDays = lookBackDays
	}

	if req.DefaultMaxTimeInMinute == "" {
		response.ReturnError(c, business_error.ParamIsNil("defaultMaxTimeInMinute"))
		return
	}
	if defaultMaxTimeInMinute, err := strconv.Atoi(req.DefaultMaxTimeInMinute); err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("defaultMaxTimeInMinute"))
		return
	} else {
		calcOtherParam.Work.DefaultMaxTimeInMinute = defaultMaxTimeInMinute
	}

	if req.DefaultMinIdleTimeInMinute == "" {
		response.ReturnError(c, business_error.ParamIsNil("defaultMinIdleTimeInMinute"))
		return
	}
	if defaultMinIdleTimeInMinute, err := strconv.Atoi(req.DefaultMinIdleTimeInMinute); err != nil {
		response.ReturnError(c, business_error.ParamIsWrong("defaultMinIdleTimeInMinute"))
		return
	} else {
		calcOtherParam.Work.DefaultMinIdleTimeInMinute = defaultMinIdleTimeInMinute
	}

	if req.WorkLoadAggregateType == "" {
		response.ReturnError(c, business_error.ParamIsNil("workLoadAggregateType"))
		return
	}
	calcOtherParam.HourSummary.WorkLoadAggregateType = calc_dynamic_param.WorkLoadAggregateType(req.WorkLoadAggregateType)

	processImplementation := service.DomainHolder.ProcessService.GetImplementationById(int64(req.ProcessImplId))

	service.Holder.EfficiencyComputeService.SaveCalcOtherParam(processImplementation, calcOtherParam)
	response.ReturnSuccessEmptyJson(c)
}

func (p EfficiencyHandler) convertCalcOtherParam2Vo(param *calc_dynamic_param.CalcOtherParam) *vo.CalcOtherParamVo {
	v := vo.CalcOtherParamVo{
		AttendanceAbsencePenaltyHour: fmt.Sprintf("%v", param.Attendance.AttendanceAbsencePenaltyHour),
		MaxRunUpTimeInMinute:         strconv.Itoa(param.Attendance.MaxRunUpTimeInMinute),
		LookBackDays:                 strconv.Itoa(param.Work.LookBackDays),
		DefaultMaxTimeInMinute:       strconv.Itoa(param.Work.DefaultMaxTimeInMinute),
		DefaultMinIdleTimeInMinute:   strconv.Itoa(param.Work.DefaultMinIdleTimeInMinute),
		WorkLoadAggregateType:        string(param.HourSummary.WorkLoadAggregateType),
	}

	workLoads := make([]string, 0)
	for _, workLoad := range param.Work.WorkLoadUnits {
		workLoads = append(workLoads, fmt.Sprintf("%s:%s", workLoad.Name, workLoad.Code))
	}
	v.WorkLoadUnits = strings.Join(workLoads, ",")

	return &v
}
