package my_error

import (
	"errors"
	"fmt"
)

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
)

func NewError(errorCode int, errorMsg string) error {
	return errors.New(fmt.Sprintf("%v, %v", errorCode, errorMsg))
}
