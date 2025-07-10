/*
* @Author: supbro
* @Date:   2025/6/12 10:04
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 10:04
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

// 测试把第二天/昨天两天的数据，在今天考勤上下班时间范围内的数据，归属到今天
func TestAddCrossDayData(T *testing.T) {
	// 测试第二天的数据
	ctx := BuildTestCtx()
	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 18:00:00")
		return &t
	}()
	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			EndTime: ctx.TodayAttendanceEndTime,
		},
	}

	ctx.TomorrowWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 17:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 20:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 18:10:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 20:20:00")
					return &t
				}(),
			},
		},
	}
	ctx.TodayWorkList = make([]domain.Actionable, 0)

	ctxRes := golang_node.AddCrossDayData(&ctx)

	assert.Equal(T, "A1", ctxRes.TodayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 1, len(ctxRes.TodayWorkList))

	// 测试处理昨天的数据
	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 08:00:00")
		return &t
	}()
	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.TodayAttendanceStartTime,
		},
	}

	ctx.YesterdayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 09:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 12:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 07:10:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 15:20:00")
					return &t
				}(),
			},
		},
	}
	ctx.TodayWorkList = make([]domain.Actionable, 0)

	ctxRes = golang_node.AddCrossDayData(&ctx)

	assert.Equal(T, "B1", ctxRes.TodayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 2, len(ctxRes.TodayWorkList))
}
