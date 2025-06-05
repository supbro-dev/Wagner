package my_error

const (
	// todo 检查serverError使用是否准确
	ServerOccurredErrorCode    int    = 9000
	ServerOccurredErrorMsg     string = "系统异常，请联系管理员"
	ServerOccurredErrorWithMsg string = "系统异常，请联系管理员，具体信息为："

	ParamNilCode int    = 9001
	ParamNilMsg  string = "参数必填"

	ParamErrorCode int    = 9002
	ParamErrorMsg  string = "参数格式有误"

	ScriptWrongTypeCode int    = 9003
	ScriptWrongTypeMsg  string = "脚本类型错误"

	ScriptNotExistCode int    = 9004
	ScriptNotExistMsg  string = "脚本不存在："

	JsonFormatWrongCode int    = 9005
	JsonFormatWrongMsg  string = "JSON格式有误:"

	ActionTypeError int    = 9006
	ActionTypeMsg   string = "动作类型错误"
)
