/*
* @Author: supbro
* @Date:   2025/6/5 18:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 18:12
 */
package domain

import (
	"time"
)

type HourSummaryResult struct {
	AggregateKey HourSummaryAggregateKey

	WorkLoad map[string]float64 // 小时内工作量

	// 工时
	DirectWorkTime   int
	IndirectWorkTime int
	IdleTime         int
	RestTime         int
	AttendanceTime   int
	// 额外属性
	Properties map[string]interface{}

	// 环节信息
	Process StandardPosition
}

// 聚合key
type HourSummaryAggregateKey struct {
	EmployeeNumber string
	OperateTime    time.Time

	ProcessCode    string
	WorkplaceCode  string
	PropertyValues string
}

type AggregateDimension string

var (
	AggregateDimension_Process  AggregateDimension = "process"  //员工+作业环节聚合
	AggregateDimension_Position AggregateDimension = "position" //员工+作业岗位
)

type IsCrossPosition string

var (
	IsCrossPosition_Cross   IsCrossPosition = "cross"
	IsCrossPosition_NoCross IsCrossPosition = "noCross"
	IsCrossPosition_All     IsCrossPosition = "all"
)

// 根据聚合属性构建一个用来聚合的汇总结果
func MakeHourSummaryResult(aggregateKey HourSummaryAggregateKey, work Actionable, field2Column map[string]string) HourSummaryResult {
	result := HourSummaryResult{
		AggregateKey: aggregateKey,
		WorkLoad:     make(map[string]float64),
		Properties:   make(map[string]interface{}),
		Process:      work.GetAction().Process,
	}

	for fieldName, columnName := range field2Column {
		if value, exist := work.GetAction().Properties[fieldName]; exist {
			result.Properties[columnName] = value
		}
	}
	return result
}

func (r *HourSummaryResult) MergeTime(work Actionable, duration float64) {
	durationTime := int(duration)
	switch work.GetAction().ActionType {
	case DIRECT_WORK:
		r.DirectWorkTime += durationTime
		r.AttendanceTime += durationTime
	case INDIRECT_WORK:
		r.IndirectWorkTime += durationTime
		r.AttendanceTime += durationTime
	case IDLE:
		r.IdleTime += durationTime
		r.AttendanceTime += durationTime
	case REST:
		r.RestTime += durationTime
	}
}

func (r *HourSummaryResult) MergeWorkLoad(workLoad map[string]float64, workLoadUnits []string, proportion float64) {
	// 遍历所有工作负载单位
	for _, unit := range workLoadUnits {
		// 获取当前对象的负载值（如果不存在则为0）
		thisValue := float64(0)
		if val, exists := r.WorkLoad[unit]; exists {
			thisValue = val
		}

		// 获取传入的负载值（如果不存在则为0）
		thatValue := float64(0)
		if val, exists := workLoad[unit]; exists {
			thatValue = val * proportion
		}

		// 合并值
		r.WorkLoad[unit] = thisValue + thatValue
	}

}
