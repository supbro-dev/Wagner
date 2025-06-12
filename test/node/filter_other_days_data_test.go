/*
* @Author: supbro
* @Date:   2025/6/12 10:13
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 10:13
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

// 测试把今天数据中开始时间小于昨天考勤下班时间，或大于第二天考勤上班时间的数据过滤掉丢弃
func TestFilterOtherDaysData(T *testing.T) {
	// 测试删除昨天的数据
	ctx := BuildTestCtx()
	ctx.YesterdayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-12 04:00:00")
		return &t
	}()
	ctx.YesterdayAttendance = &domain.Attendance{
		Action: domain.Action{
			EndTime: ctx.YesterdayAttendanceEndTime,
		},
	}

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 03:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 10:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 05:10:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 12:20:00")
					return &t
				}(),
			},
		},
	}
	ctxRes := golang_node.FilterOtherDaysData(&ctx)

	assert.Equal(T, "A2", ctxRes.TodayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 1, len(ctxRes.TodayWorkList))

	// 测试删除第二天的数据
	ctx = BuildTestCtx()
	ctx.TomorrowAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 04:00:00")
		return &t
	}()
	ctx.TomorrowAttendance = &domain.Attendance{
		Action: domain.Action{
			StartTime: ctx.TomorrowAttendanceStartTime,
		},
	}

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 10:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 20:10:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 02:20:00")
					return &t
				}(),
			},
		},
	}
	ctxRes = golang_node.FilterOtherDaysData(&ctx)

	assert.Equal(T, "B2", ctxRes.TodayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 1, len(ctxRes.TodayWorkList))
}
