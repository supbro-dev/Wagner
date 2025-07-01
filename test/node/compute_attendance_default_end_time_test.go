/*
* @Author: supbro
* @Date:   2025/6/12 08:56
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 08:56
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

// 缺卡情况设置默认下班时间
func TestComputeAttendanceDefaultEndTime(t *testing.T) {
	ctx := BuildTestCtx()

	ctx.TodayAttendance = &domain.Attendance{}

	// 有排班使用排班下班时间
	ctx.CalcStartTime = time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local)
	todayAttendanceStartTime := time.Date(2025, 6, 12, 8, 0, 0, 0, time.Local)
	ctx.TodayAttendanceStartTime = &todayAttendanceStartTime
	ctx.TodayScheduling = &domain.Scheduling{
		Action: domain.Action{
			EndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 18:30:00")
				return &t
			}(),
		},
	}
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 12
	ctxRes := golang_node.ComputeAttendanceDefaultEndTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 12, 18, 30, 0, 0, time.Local), *ctxRes.TodayAttendanceEndTime)

	// 取惩罚下班时间
	ctx = BuildTestCtx()
	ctx.CalcStartTime = time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local)
	todayAttendanceStartTime = time.Date(2025, 6, 12, 8, 0, 0, 0, time.Local)
	ctx.TodayAttendanceStartTime = &todayAttendanceStartTime
	ctx.TodayScheduling = &domain.Scheduling{
		Action: domain.Action{
			EndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-12 21:30:00")
				return &t
			}(),
		},
	}
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 12
	ctxRes = golang_node.ComputeAttendanceDefaultEndTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 12, 20, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceEndTime)

	// 取默认下班时间(第二天0点)
	ctx = BuildTestCtx()
	ctx.CalcStartTime = time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local)
	todayAttendanceStartTime = time.Date(2025, 6, 12, 8, 0, 0, 0, time.Local)
	ctx.TodayAttendanceStartTime = &todayAttendanceStartTime
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 20
	ctxRes = golang_node.ComputeAttendanceDefaultEndTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 0, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceEndTime)

	// 取默认下班时间(第二天12点)
	ctx = BuildTestCtx()
	ctx.CalcStartTime = time.Date(2025, 6, 13, 23, 0, 0, 0, time.Local)
	todayAttendanceStartTime = time.Date(2025, 6, 12, 20, 0, 0, 0, time.Local)
	ctx.TodayAttendanceStartTime = &todayAttendanceStartTime
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 20
	ctxRes = golang_node.ComputeAttendanceDefaultEndTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 12, 0, 0, 0, time.Local), *ctxRes.TodayAttendanceEndTime)

	// 测试前一天
	// 取默认下班时间(第二天0点)
	ctx = BuildTestCtx()
	ctx.CalcStartTime = time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local)
	yesterdayAttendanceStartTime := time.Date(2025, 6, 12, 8, 0, 0, 0, time.Local)
	ctx.YesterdayAttendanceStartTime = &yesterdayAttendanceStartTime
	ctx.CalcOtherParam.Attendance.AttendanceAbsencePenaltyHour = 20
	ctxRes = golang_node.ComputeAttendanceDefaultEndTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 0, 0, 0, 0, time.Local), *ctxRes.YesterdayAttendanceEndTime)
}
