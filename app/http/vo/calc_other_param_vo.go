/*
* @Author: supbro
* @Date:   2025/7/10 10:39
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/10 10:39
 */
package vo

type CalcOtherParamVo struct {
	AttendanceAbsencePenaltyHour string `json:"attendanceAbsencePenaltyHour"`
	MaxRunUpTimeInMinute         string `json:"maxRunUpTimeInMinute"`
	WorkLoadUnits                string `json:"workLoadUnits"`
	LookBackDays                 string `json:"lookBackDays"`
	DefaultMaxTimeInMinute       string `json:"defaultMaxTimeInMinute"`
	DefaultMinIdleTimeInMinute   string `json:"defaultMinIdleTimeInMinute"`
	WorkLoadAggregateType        string `json:"workLoadAggregateType"`
}
