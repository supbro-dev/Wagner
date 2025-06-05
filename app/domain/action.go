package domain

import "time"

type Action struct {
	// 动作操作人
	EmployeeNumber string
	// 工作点编码
	WorkplaceCode string
	// 动作发生日期
	OperateDay time.Time
	// 动作编码
	ActionCode string
	// 动作类型
	ActionType ActionType
	// 原始开始时间
	StartTime time.Time
	// 原始结束时间
	EndTime time.Time
	// 计算后开始时间
	ComputedStartTime time.Time
	// 计算后技术时间
	ComputedEndTime time.Time
	// 额外属性
	Properties map[string]interface{} `copier:"-"` // 字段排除
	// 加工处理过程
	ProcessList []string
}

// 直接作业
type DirectWork struct {
	Action
	// 工作量
	WorkLoad map[string]float32
	// 任务的发起人
	Starter string
}

// 间接作业
type IndirectWork struct {
	Action
}

// 作业（直接&间接）
type Work interface {
	GetWorkType() ActionType
	GetComputedStartTime() time.Time
	GetComputedEndTime() time.Time
	GetWorkLoad() map[string]float64
}

func (a DirectWork) GetWorkType() ActionType {
	return DIRECT_WORK
}

func (a DirectWork) GetComputedStartTime() time.Time {
	return a.ComputedStartTime
}

func (a DirectWork) GetComputedEndTime() time.Time {
	return a.ComputedEndTime
}

func (a DirectWork) GetWorkLoad() map[string]float64 {
	return a.WorkLoad
}

func (a IndirectWork) GetWorkType() ActionType {
	return INDIRECT_WORK
}

func (a IndirectWork) GetComputedStartTime() time.Time {
	return a.ComputedStartTime
}

func (a IndirectWork) GetComputedEndTime() time.Time {
	return a.ComputedEndTime
}

func (a IndirectWork) GetWorkLoad() map[string]float64 {
	return make(map[string]float32)
}

// 考勤
type Attendance struct {
	Action
}

// 排班
type Scheduling struct {
	Action
}

type ActionType string

var (
	// 直接作业
	DIRECT_WORK ActionType = "DirectWork"
	// 间接作业
	INDIRECT_WORK ActionType = "IndirectWork"
	// 考勤
	ATTENDANCE ActionType = "Attendance"
	// 排班
	SCHEDULING ActionType = "Scheduling"
	// 离岗休息
	SHORT_BREAK ActionType = "ShortBreak"
	// 外出就餐
	MEAL_TIME ActionType = "MealTime"
)

// 任务的开始人放在扩展属性里
const STARTER = "starter"

// 判断动作的执行人(Starter)和完成人(employeeNumber)是否为同一个工号
// Parameters:
// Returns: 操作人是否改变
func (a *Action) IsChangeOperator() bool {
	return a.Properties[STARTER].(string) == a.EmployeeNumber
}
