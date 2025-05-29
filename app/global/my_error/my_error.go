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
)
