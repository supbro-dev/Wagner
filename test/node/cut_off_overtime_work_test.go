/*
* @Author: supbro
* @Date:   2025/6/12 11:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 11:48
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

func TestCutOffOvertimeWork(t *testing.T) {
	// 超长截断取环节自身配置
	ctx := BuildTestCtx()
	ctx.CalcOtherParam.Work.DefaultMaxTimeInMinute = 120

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
				Process: &domain.ProcessPosition{
					Properties: map[string]interface{}{
						"maxTimeInMinute": 60,
					},
				},
			},
		},
	}

	ctxRes := golang_node.CutOffOvertimeWork(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)

	// 取默认配置
	ctx = BuildTestCtx()
	ctx.CalcOtherParam.Work.DefaultMaxTimeInMinute = 120

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.CutOffOvertimeWork(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)

	// 考勤截断
	ctx = BuildTestCtx()
	ctx.CalcOtherParam.Work.DefaultMaxTimeInMinute = 120

	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 05:30:00")
		return &t
	}()

	ctx.TodayAttendance = &domain.Attendance{
		Action: domain.Action{
			ComputedEndTime: ctx.TodayAttendanceEndTime,
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
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.CutOffOvertimeWork(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 5, 30, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
}
