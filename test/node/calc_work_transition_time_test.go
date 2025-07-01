/*
* @Author: supbro
* @Date:   2025/6/12 12:50
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 12:50
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

// 测试处理上班打卡到第一次作业开始的时间，处理最后一次作业完成到下班打卡的时间
func TestCalcWorkTransitionTime(t *testing.T) {
	// 测试上下班时间
	ctx := BuildTestCtx()
	ctx.CalcOtherParam.Attendance.MaxRunUpTimeInMinute = 30
	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 20

	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
		return &t
	}()

	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 18:00:00")
		return &t
	}()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:20:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 09:00:00")
					return &t
				}(),
			},
		},

		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:20:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:45:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.CalcWorkTransitionTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 18, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)

	// 测试休息后的开班时间
	ctx = BuildTestCtx()
	ctx.CalcOtherParam.Attendance.MaxRunUpTimeInMinute = 30

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:20:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 09:00:00")
					return &t
				}(),
			},
		},

		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:20:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:45:00")
					return &t
				}(),
			},
		},
	}

	ctx.TodayRestList = []*domain.Rest{
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:10:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:10:00")
					return &t
				}(),
			},
		},
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:05:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.CalcWorkTransitionTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 8, 10, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 17, 5, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedStartTime)

}
