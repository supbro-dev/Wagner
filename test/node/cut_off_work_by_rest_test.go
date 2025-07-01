/*
* @Author: supbro
* @Date:   2025/6/12 12:27
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 12:27
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

func TestCutOffWorkByRest(t *testing.T) {
	// 测试一段休息
	ctx := BuildTestCtx()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
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

	ctxRes := golang_node.CutOffWorkByRest(&ctx)

	assert.Equal(t, time.Date(2025, 6, 13, 7, 10, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 20, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)

	// 测试两段休息
	ctx = BuildTestCtx()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
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
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:20:00")
					return &t
				}(),
			},
		},
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:40:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 07:50:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.CutOffWorkByRest(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 10, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 20, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 40, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 7, 50, 0, 0, time.Local), *ctxRes.TodayWorkList[2].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[2].GetAction().ComputedEndTime)

}
