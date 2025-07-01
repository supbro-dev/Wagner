/*
* @Author: supbro
* @Date:   2025/6/12 11:44
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 11:44
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

func TestPaddingUnfinishedWorkEndTime(t *testing.T) {
	ctx := BuildTestCtx()
	ctx.CalcStartTime = func() time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 06:00:00")
		return t
	}()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				StartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 05:00:00")
					return &t
				}(),
			},
		},
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				StartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-12 20:10:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.PaddingUnfinishedWorkEndTime(&ctx)

	assert.Equal(t, time.Date(2025, 6, 13, 6, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[0].GetAction().ComputedEndTime)
	assert.Equal(t, time.Date(2025, 6, 13, 6, 0, 0, 0, time.Local), *ctxRes.TodayWorkList[1].GetAction().ComputedEndTime)
}
