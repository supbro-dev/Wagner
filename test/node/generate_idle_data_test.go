/*
* @Author: supbro
* @Date:   2025/6/12 13:21
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 13:21
 */
package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wagner/app/domain"
	"wagner/app/service"
	"wagner/app/service/calc_node/golang_node"
	"wagner/app/utils/datetime_util"
)

// 测试生成闲置工时
func TestGenerateIdleData(t *testing.T) {
	ctx := BuildTestCtx()

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 1

	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
		return &t
	}()

	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 19:00:00")
		return &t
	}()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.DirectWork{
			Action: domain.Action{
				ActionCode: "B1",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 09:00:00")
					return &t
				}(),
			},
		},

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 15:00:00")
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
					t, _ := datetime_util.ParseDatetime("2025-06-13 12:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 12:30:00")
					return &t
				}(),
			},
		},
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:50:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 18:30:00")
					return &t
				}(),
			},
		},
	}

	ctxRes := golang_node.GenerateIdleData(&ctx)
	assert.Equal(t, 5, len(ctxRes.TodayIdleList))
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 8, 30, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 9, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[1].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 12, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[1].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 12, 30, 0, 0, time.Local), *ctxRes.TodayIdleList[2].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 15, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[2].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 17, 45, 0, 0, time.Local), *ctxRes.TodayIdleList[3].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 17, 50, 0, 0, time.Local), *ctxRes.TodayIdleList[3].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 18, 30, 0, 0, time.Local), *ctxRes.TodayIdleList[4].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 19, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[4].GetAction().ComputedEndTime)

	// 测试全天闲置，全天只有上下班
	ctx = BuildTestCtx()
	ctx.Employee.PositionCode = "checker"

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 1

	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
		return &t
	}()

	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 19:00:00")
		return &t
	}()

	standardPositionMock := new(StandardPositionMock)
	service.DomainHolder = service.DomainServiceHolder{
		StandardPositionService: standardPositionMock,
	}
	standardPositionMock.On("FindPositionFirstProcess", "checker", "FOOD", "ConvenientFood").Return(&domain.StandardPosition{
		Code: "C1",
	})

	ctxRes = golang_node.GenerateIdleData(&ctx)
	assert.Equal(t, 1, len(ctxRes.TodayIdleList))
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 19, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedEndTime)

	// 测试全天闲置，全天有上下班和休息的情况
	ctx = BuildTestCtx()
	ctx.Employee.PositionCode = "checker"

	ctx.CalcOtherParam.Work.DefaultMinIdleTimeInMinute = 1

	ctx.TodayAttendanceStartTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 08:00:00")
		return &t
	}()

	ctx.TodayAttendanceEndTime = func() *time.Time {
		t, _ := datetime_util.ParseDatetime("2025-06-13 19:00:00")
		return &t
	}()

	ctx.TodayRestList = []*domain.Rest{
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 12:00:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 12:30:00")
					return &t
				}(),
			},
		},
		{
			Action: domain.Action{
				ActionType: domain.REST,
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 17:50:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 18:30:00")
					return &t
				}(),
			},
		},
	}

	ctxRes = golang_node.GenerateIdleData(&ctx)
	assert.Equal(t, 3, len(ctxRes.TodayIdleList))
	assert.Equal(t, time.Date(2025, 6, 13, 8, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 12, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[0].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 12, 30, 0, 0, time.Local), *ctxRes.TodayIdleList[1].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 17, 50, 0, 0, time.Local), *ctxRes.TodayIdleList[1].GetAction().ComputedEndTime)

	assert.Equal(t, time.Date(2025, 6, 13, 18, 30, 0, 0, time.Local), *ctxRes.TodayIdleList[2].GetAction().ComputedStartTime)
	assert.Equal(t, time.Date(2025, 6, 13, 19, 0, 0, 0, time.Local), *ctxRes.TodayIdleList[2].GetAction().ComputedEndTime)

}
