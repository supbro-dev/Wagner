/*
* @Author: supbro
* @Date:   2025/6/12 08:04
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 08:04
 */
package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wagner/app/domain"
	"wagner/app/service/calc_node/golang_node"
	"wagner/app/utils/datetime_util"
)

// 测试前一天/当天/第二天的上下班时间是否都已设置完成
func TestSetCrossDayAttendance(t *testing.T) {
	ctx := BuildTestCtx()

	todayAttendance := &domain.Attendance{
		Action: domain.Action{
			StartTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 08:00:00")
				return &t
			}(),
			EndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")
				return &t
			}(),
		},
	}

	yesterdayAttendance := &domain.Attendance{
		Action: domain.Action{
			StartTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-11 10:00:00")
				return &t
			}(),
			EndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-11 22:00:00")
				return &t
			}(),
		},
	}

	tomorrowAttendance := &domain.Attendance{
		Action: domain.Action{
			StartTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 07:00:00")
				return &t
			}(),
			EndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 21:00:00")
				return &t
			}(),
		},
	}

	ctx.TodayAttendance = todayAttendance
	ctx.TomorrowAttendance = tomorrowAttendance
	ctx.YesterdayAttendance = yesterdayAttendance

	ctxRes := golang_node.SetCrossDayAttendance(&ctx)

	assert.Equal(t, time.Date(2025, 6, 12, 8, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 12, 18, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceEndTime)
	assert.Equal(t, time.Date(2025, 6, 11, 10, 0, 0, 0, time.Local), *ctxRes.YesterdayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 11, 22, 0, 0, 0, time.Local), *ctxRes.YesterdayAttendanceEndTime)
	assert.Equal(t, time.Date(2025, 6, 12, 7, 0, 0, 0, time.Local), *ctxRes.TomorrowAttendanceStartTime)
}
