/*
* @Author: supbro
* @Date:   2025/6/12 13:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 13:48
 */
package node

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"wagner/app/domain"
	"wagner/app/service"
	"wagner/app/service/calc_node/golang_node"
	"wagner/app/utils/datetime_util"
)

// 匹配休息的归属环节（休息归属于他上一个环节， 上一个环节不存在，则归属于休息的下一个环节）
func TestMatchRestProcess(t *testing.T) {
	// 测试休息取上一个环节的环节编码
	ctx := BuildTestCtx()

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
				Process: &domain.StandardPosition{
					Code: "P1",
				},
			},
		},

		&domain.Rest{
			Action: domain.Action{
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

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 09:00:00")
					return &t
				}(),
				Process: &domain.StandardPosition{
					Code: "P2",
				},
			},
		},
	}

	ctxRes := golang_node.MatchRestProcess(&ctx)
	assert.Equal(t, "P1", ctxRes.TodayWorkList[1].GetAction().Process.Code)

	// 测试休息取下一个环节的环节编码
	ctx = BuildTestCtx()

	ctx.TodayWorkList = []domain.Actionable{
		&domain.Rest{
			Action: domain.Action{
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

		&domain.IndirectWork{
			Action: domain.Action{
				ActionCode: "B2",
				ComputedStartTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 08:30:00")
					return &t
				}(),
				ComputedEndTime: func() *time.Time {
					t, _ := datetime_util.ParseDatetime("2025-06-13 09:00:00")
					return &t
				}(),
				Process: &domain.StandardPosition{
					Code: "P2",
				},
			},
		},
	}

	ctxRes = golang_node.MatchRestProcess(&ctx)
	assert.Equal(t, "P2", ctxRes.TodayWorkList[0].GetAction().Process.Code)

	// 测试全天休息员工岗位下的第一个环节
	ctx = BuildTestCtx()
	ctx.Employee.PositionCode = "checker"

	ctx.TodayWorkList = []domain.Actionable{
		&domain.Rest{
			Action: domain.Action{
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

	standardPositionMock := new(StandardPositionMock)
	service.DomainHolder = service.DomainServiceHolder{
		StandardPositionService: standardPositionMock,
	}
	standardPositionMock.On("FindPositionFirstProcess", "checker", "FOOD", "ConvenientFood").Return(&domain.StandardPosition{
		Code: "C1",
	})

	ctxRes = golang_node.MatchRestProcess(&ctx)
	assert.Equal(t, "C1", ctxRes.TodayWorkList[0].GetAction().Process.Code)
}

type StandardPositionInterface interface {
	FindPositionFirstProcess(positionCode string, industryCode, subIndustryCode string) *domain.StandardPosition
}

type StandardPositionMock struct {
	mock.Mock
}

func (service *StandardPositionMock) FindStandardPositionListByIndustry(industryCode, subIndustryCode string) []*domain.StandardPosition {
	return nil
}

func (service *StandardPositionMock) FindStandardPositionByWorkplace(workplaceCode string) []*domain.StandardPosition {
	return nil
}

func (service *StandardPositionMock) FindStandardPositionByIndustry(industryCode, subIndustryCode string) []*domain.StandardPosition {
	return nil
}

func (service *StandardPositionMock) FindPositionFirstProcess(positionCode string, industryCode, subIndustryCode string) *domain.StandardPosition {
	args := service.Called(positionCode, industryCode, subIndustryCode) // 捕获调用参数
	return args.Get(0).(*domain.StandardPosition)
}
