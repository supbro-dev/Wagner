/*
* @Author: supbro
* @Date:   2025/6/6 11:58
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 11:58
 */
package datetime_util

import (
	"time"
)

var (
	DateLayout       = "2006-01-02"
	DateTimeLayout   = "2006-01-02 15:04:05"
	DateTimeMsLayout = "2006-01-02 15:04:05.000"
)

func ParseDatetime(datetime string) (time.Time, error) {
	parse, err := time.ParseInLocation(DateTimeLayout, datetime, time.Local)

	return parse, err
}

func ParseDate(date string) (time.Time, error) {
	parse, err := time.ParseInLocation(DateLayout, date, time.Local)

	return parse, err
}

func FormatDatetime(time time.Time) string {
	return time.Format(DateTimeLayout)
}
