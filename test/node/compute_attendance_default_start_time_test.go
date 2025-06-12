/*
* @Author: supbro
* @Date:   2025/6/12 10:31
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 10:31
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

// 计算处理当天默认考勤开始时间
func TestComputeAttendanceDefaultStartTime(t *testing.T) {
	// 测试排班时间
	ctx := BuildTestCtx()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			ComputedEndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctx.TodayScheduling = &domain.Scheduling{
		Action: domain.Action{
			StartTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 09:00:00")
				return &t
			}(),
		},
	}

	ctxRes := golang_node.ComputeAttendanceDefaultStartTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 12, 9, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 12, 9, 0, 0, 0, time.Local), *ctxRes.TodayAttendance.ComputedStartTime)

	// 测试使用默认时间
	ctx = BuildTestCtx()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			ComputedEndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctxRes = golang_node.ComputeAttendanceDefaultStartTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local), *ctxRes.TodayAttendance.ComputedStartTime)

	ctx = BuildTestCtx()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 10:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			ComputedEndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctxRes = golang_node.ComputeAttendanceDefaultStartTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 11, 12, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 11, 12, 0, 0, 0, time.Local), *ctxRes.TodayAttendance.ComputedStartTime)

	// 测试惩罚时间
	ctx = BuildTestCtx()
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 10
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 10:00:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			ComputedEndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctxRes = golang_node.ComputeAttendanceDefaultStartTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceStartTime)
	assert.Equal(t, time.Date(2025, 6, 12, 0, 0, 0, 0, time.Local), *ctxRes.TodayAttendance.ComputedStartTime)

}
