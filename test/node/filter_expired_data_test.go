/*
* @Author: supbro
* @Date:   2025/6/12 10:18
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 10:18
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

// 测试根据参数配置的每个operateDay只计算x天内的数据，把早于这天的数据丢弃
func TestFilterExpiredData(T *testing.T) {
	ctx := BuildTestCtx()
	// 3天过期时间
	ctx.OperateDay = func() time.Time {
		t, _ := datetime_util.ParseDate("2025-06-13")
		return t
	}()
	ctx.CalcOtherParam.Work.LookBackDays = 3

	ctx.YesterdayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-09 00:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 05:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "A2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 01:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 05:00:00")
					return &t
				}(),
			},
		},
	}

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-09 00:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-09 05:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 01:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 05:00:00")
					return &t
				}(),
			},
		},
	}

	ctx.TomorrowWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "C1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-09 23:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 00:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "C2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 00:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-10 05:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.FilterExpiredData(&ctx)

	assert.Equal(T, "A2", ctxRes.YesterdayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 1, len(ctxRes.YesterdayWorkList))

	assert.Equal(T, "B2", ctxRes.TodayWorkList[0].GetAction().ActionCode)
	assert.Equal(T, 1, len(ctxRes.TodayWorkList))

	assert.Equal(T, 0, len(ctxRes.TomorrowWorkList))
}
