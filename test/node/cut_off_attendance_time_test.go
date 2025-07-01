/*
* @Author: supbro
* @Date:   2025/6/12 09:46
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 09:46
 */
package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wagner/app/domain"
	"wagner/app/service/calc/calc_node/golang_node"
	"wagner/app/utils/datetime_util"
)

// 测试处理昨天/当天/第二天，3天考勤计算出来有重叠的场景
func TestCutOffAttendanceTime(t *testing.T) {
	// 验证昨天下班时间大于当天上班时间
	ctx := BuildTestCtx()

	ctx.YesterdayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 05:00:00")
		return &t
	}()

	ctx.YesterdayAttendance = &domain.Attendance{
		Action: domain.Action{
			EndTime: ctx.YesterdayAttendanceEndTime,
		},
	}

	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 04:00:00")
		return &t
	}()
	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.TodayAttendanceStartTime,
		},
	}

	ctxRes := golang_node.CutOffAttendanceTime(&ctx)

	assert.Equal(t, time.Date(2025, 6, 12, 3, 59, 59, 0, time.Local), *ctxRes.YesterdayAttendanceEndTime)
	assert.Equal(t, time.Date(2025, 6, 12, 3, 59, 59, 0, time.Local), *ctxRes.YesterdayAttendance.ComputedEndTime)

	// 验证昨天上班时间大于昨天下班时间
	ctx = BuildTestCtx()
	ctx.YesterdayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-11 19:00:00")
		return &t
	}()
	ctx.YesterdayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-11 18:00:00")
		return &t
	}()

	ctx.YesterdayAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.YesterdayAttendanceStartTime,
			EndTime:   ctx.YesterdayAttendanceEndTime,
		},
	}

	ctxRes = golang_node.CutOffAttendanceTime(&ctx)

	assert.Equal(t, time.Date(2025, 6, 11, 18, 0, 0, 0, time.Local), *ctxRes.YesterdayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 11, 18, 0, 0, 0, time.Local), *ctxRes.YesterdayAttendance.ComputedStartTime)

	// 验证当天下班时间大于第二天上班时间
	ctx = BuildTestCtx()
	ctx.TomorrowAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 04:00:00")
		return &t
	}()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 05:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			EndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctx.TomorrowAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.TomorrowAttendanceStartTime,
		},
	}

	ctxRes = golang_node.CutOffAttendanceTime(&ctx)

	assert.Equal(t, time.Date(2025, 6, 12, 3, 59, 59, 0, time.Local), *ctxRes.TodayAttendanceEndTime)
	assert.Equal(t, time.Date(2025, 6, 12, 3, 59, 59, 0, time.Local), *ctxRes.TodayAttendance.ComputedEndTime)

	// 验证当天上班时间大于当天下班时间
	ctx = BuildTestCtx()
	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 06:00:00")
		return &t
	}()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 05:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.TodayAttendanceStartTime,
			EndTime:   ctx.TodayAttendanceEndTime,
		},
	}

	ctxRes = golang_node.CutOffAttendanceTime(&ctx)

	assert.Equal(t, time.Date(2025, 6, 12, 5, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 12, 5, 0, 0, 0, time.Local), *ctxRes.TodayAttendance.ComputedStartTime)

}
