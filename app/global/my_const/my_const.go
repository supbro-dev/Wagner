package my_const

// SinkType标识数据加工完之后如何进行存储
type SinkType string

const (
	// 数据汇总
	SUMMARY SinkType = "SUMMARY"
	// TIME ON TASK
	TIME_ON_TASK = "TIME_ON_TASK"
	// 个人人效
	INDIVIDUAL_EFFICIENCY = "INDIVIDUAL_EFFICIENCY"
	// 个人当日状态
	INDIVIDUAL_DATE_STATUS = "INDIVIDUAL_DATE_STATUS"
)
