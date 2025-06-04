/*
* @Author: supbro
* @Date:   2025/6/4 10:20
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/4 10:20
 */
package common

import "time"

type BaseDao struct {
}

func (dao *BaseDao) Time2Date(time time.Time) Date {
	return Date(time)
}

func (dao *BaseDao) TimeList2DateList(times []time.Time) []Date {
	dateList := make([]Date, 0)
	for _, time := range times {
		dateList = append(dateList, Date(time))
	}
	return dateList
}
