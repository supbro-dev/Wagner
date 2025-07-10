/*
* @Author: supbro
* @Date:   2025/6/12 12:03
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 12:03
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

// 测试在作业之间添加合理休息
func TestAddReasonableBreakTime(t *testing.T) {
	// 不生效场景
	ctx := BuildTestCtx()

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 10

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:00:00")
					return &t
				}(),
			},
		},

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:00:00")
					return &t
				}(),
			},
		},

		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B3",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.AddReasonableBreakTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)

	// 环节上的配置生效
	ctx = BuildTestCtx()

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 10

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:00:00")
					return &t
				}(),
				Process: &domain.ProcessPosition{
					Properties: map[string]interface{}{
						"minIdleTimeInMinute": 30,
					},
				},
			},
		},

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:00:00")
					return &t
				}(),
			},
		},

		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B3",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.AddReasonableBreakTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 30, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)

	// 考虑休息情况
	ctx = BuildTestCtx()

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 10

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:00:00")
					return &t
				}(),
				Process: &domain.ProcessPosition{
					Properties: map[string]interface{}{
						"minIdleTimeInMinute": 30,
					},
				},
			},
		},

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:00:00")
					return &t
				}(),
			},
		},

		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B3",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctx.TodayRestList = []*domain.Rest{{
		Action: domain.Action{
			ActionType: domain.REST,
			ComputedStartTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-13 07:10:00")
				return &t
			}(),
			ComputedEndTime: func() *time.Time {
				t, _ := datetime_util.ParseDatetime("2025-06-13 07:20:00")
				return &t
			}(),
		},
	},
	}

	ctxRes = golang_node.AddReasonableBreakTime(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 30, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 10, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 20, 0, 0, time.Local), *ctxRes.TodayRestList[0].GetAction().ComputedEndTime)
}
