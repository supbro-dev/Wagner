/*
* @Author: supbro
* @Date:   2025/6/5 18:12
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 18:12
 */
package domain

import (
	mapset "github.com/deckarep/golang-set/v2"
	"time"
	"wagner/app/utils/reflect_util"
)

type HourSummaryResult struct {
	OperateTime   time.Time          // 作业时间（小时）
	TotalWorkLoad map[string]float64 // 小时内工作量

	// 额外属性
	Properties map[string]interface{}
}

// 根据聚合属性构建一个用来聚合的汇总结果
func MakeHourSummaryResult(aggregateValueObj interface{}, aggregateFields []string) HourSummaryResult {
	result := HourSummaryResult{}

	for _, field := range aggregateFields {
		value, err := reflect_util.GetField(aggregateValueObj, field)
		if err != nil {
			panic(err)
		}

		err = reflect_util.SetField(result, field, value)
		if err != nil {
			hasField, err := reflect_util.HasField(result, field)
			if !hasField && err != nil {
				if result.Properties != nil {
					result.Properties = make(map[string]interface{})
				} else {
					result.Properties[field] = value
				}
			}
		}
	}

	return result
}

func (r HourSummaryResult) MergeTimeAndWorkLoad(that HourSummaryResult, workLoadUnits mapset.Set[string], proportion float64) *HourSummaryResult {

}

func (r HourSummaryResult) MergeOtherProperties(that HourSummaryResult) *HourSummaryResult {

}

func (r HourSummaryResult) MergeWorkLoad(workLoad map[string]float64, workLoadUnits mapset.Set[string], proportion float64) map[string]float64 {
	total := make(map[string]float64)

	// 遍历所有工作负载单位
	for unit := range workLoadUnits.Iter() {
		// 获取当前对象的负载值（如果不存在则为0）
		thisValue := float64(0)
		if val, exists := r.TotalWorkLoad[unit]; exists {
			thisValue = val
		}

		// 获取传入的负载值（如果不存在则为0）
		thatValue := float64(0)
		if val, exists := workLoad[unit]; exists {
			thatValue = val * proportion
		}

		// 合并值
		total[unit] = thisValue + thatValue
	}

	return total
}
