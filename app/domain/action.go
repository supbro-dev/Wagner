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
	StartTime *time.Time
	// 原始结束时间
	EndTime *time.Time
	// 计算后开始时间
	ComputedStartTime *time.Time
	// 计算后技术时间
	ComputedEndTime *time.Time
	// 额外属性
	Properties map[string]interface{} `copier:"-"` // 字段排除
	// 加工处理过程
	OperationMsgList []string
	// 环节编码
	ProcessCode string
	// 环节实例
	Process *StandardPosition
}

type Actionable interface {
	GetAction() *Action
}

// 直接作业
type DirectWork struct {
	Action
	// 工作量
	WorkLoad map[string]float64
	// 任务的发起人
	Starter string
	// 截断当前作业的作业编码
	CutOffWorkCode string
}

// 间接作业
type IndirectWork struct {
	Action
	// 截断当前作业的作业编码
	CutOffWorkCode string
}

// 使用Work声明对象或切片时不需要使用&Work,因为是DirectWork/IndirectWork的指针实现的Work接口
// 作业（直接&间接）
type Work interface {
	GetAction() *Action
	GetWorkLoad() map[string]float64
	SetCutOffWorkCode(actionCode string)
}

func (d *DirectWork) GetAction() *Action {
	return &d.Action
}

func (d *DirectWork) GetWorkLoad() map[string]float64 {
	return d.WorkLoad
}

func (d *DirectWork) SetCutOffWorkCode(actionCode string) {
	d.CutOffWorkCode = actionCode
}

func (in *IndirectWork) GetAction() *Action {
	return &in.Action
}

func (in *IndirectWork) GetWorkLoad() map[string]float64 {
	return nil
}

func (in *IndirectWork) SetCutOffWorkCode(actionCode string) {
	in.CutOffWorkCode = actionCode
}

// 考勤
type Attendance struct {
	Action
}

// 排班
type Scheduling struct {
	Action
}

// 休息
type Rest struct {
	Action
}

func (rest *Rest) GetAction() *Action {
	return &rest.Action
}

// 闲置时长
type Idle struct {
	Action
}

func (idle *Idle) GetAction() *Action {
	return &idle.Action
}

// 追加Action操作日志
func (a *Action) AppendOperationMsg(msg string) {
	if &a.OperationMsgList == nil {
		a.OperationMsgList = make([]string, 0)
	}
	a.OperationMsgList = append(a.OperationMsgList, msg)
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
	// 休息（从排班中解析得来）
	REST ActionType = "Rest"
	// 闲置工时
	IDLE ActionType = "Idle"
)

// 任务的开始人放在扩展属性里
const STARTER = "starter"

// 判断动作的执行人(Starter)和完成人(employeeNumber)是否为同一个工号
// Parameters:
// Returns: 操作人是否改变
func (a *Action) IsChangeOperator() bool {
	return a.Properties[STARTER].(string) == a.EmployeeNumber
}
