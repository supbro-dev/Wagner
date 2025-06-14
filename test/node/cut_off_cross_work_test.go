/*
* @Author: supbro
* @Date:   2025/6/12 11:58
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 11:58
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

// 测试交叉截断
func TestCutOffCrossWork(t *testing.T) {
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
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:00:00")
					return &t
				}(),
			},
		},

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:30:00")
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
					t, _ := datetime_util.ParseDatetime("2025-06-13 06:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 20:00:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.CutOffCrossWork(&ctx)
	assert.Equal(t, time.Date(2025, 6, 13, 5, 30, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, "B2", ctxRes.TodayWorkList[0].(*domain.DirectWork).CutOffWorkCode)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 30, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)
	assert.Equal(t, "B3", ctxRes.TodayWorkList[1].(*domain.IndirectWork).CutOffWorkCode)

}
