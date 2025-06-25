package business_error

import (
	"fmt"
)

type BusinessError struct {
	Type        string        // 异常类型
	Code        int           // 业务错误码
	Message     string        // 用户友好消息
	Args        []interface{} // 异常信息，填入到Message占位符中
	CausedError error         // 触发的异常
}

var (
	OsError         string = "操作系统异常"
	ParamError      string = "传参异常"
	HttpInvokeError string = "HTTP调用异常"
	SystemError     string = "系统异常"
	LockHandleError string = "锁异常"
	ComputeError    string = "计算异常"
)

// 实现 error 接口
func (e *BusinessError) Error() string {
	if e.CausedError != nil {
		return fmt.Sprintf("code=%d, msg=%s, invokedError:%v", e.Code, e.Message, e.CausedError)
	} else if e.Args == nil {
		return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Message)
	} else {
		return fmt.Sprintf("code=%d, msg=%s", e.Code, fmt.Sprintf(e.Message, e.Args))
	}
}

func ServerErrorCausedBy(err error) *BusinessError {
	return &BusinessError{SystemError, 9000, "系统异常，请联系管理员", nil, err}
}

func ServerOccurredError(errorType string, args ...interface{}) *BusinessError {
	if len(args) == 0 {
		return &BusinessError{errorType, 9000, "系统异常，请联系管理员", nil, nil}
	} else {
		return &BusinessError{errorType, 9000, "系统异常，请联系管理员，具体信息为:%v", args, nil}
	}
}

// PARAM_ERROR
func ElScriptMustUseMap() *BusinessError {
	return &BusinessError{ParamError, 9101, "El表达式传参必须使用map类型", nil, nil}
}

func ScriptWrongType() *BusinessError {
	return &BusinessError{ParamError, 9102, "脚本类型错误", nil, nil}
}

func ParamIsNil(paramNames ...interface{}) *BusinessError {
	return &BusinessError{ParamError, 9103, "参数必填:%v", paramNames, nil}
}

func ParamIsWrong(paramNames ...interface{}) *BusinessError {
	return &BusinessError{ParamError, 9104, "参数错误:%v", paramNames, nil}
}

// LOCK_ERROR
func LockFailureBySystemError(err error) *BusinessError {
	return &BusinessError{LockHandleError, 9201, "加锁异常", nil, err}
}

func UnlockFailureBySystemError(err error) *BusinessError {
	return &BusinessError{LockHandleError, 9202, "解锁异常", nil, err}
}

func LockFailure() *BusinessError {
	return &BusinessError{LockHandleError, 9203, "加锁失败", nil, nil}
}

func UnlockFailure() *BusinessError {
	return &BusinessError{LockHandleError, 9204, "解锁失败", nil, nil}
}

// COMPUTE
func InjectDataError(err error) *BusinessError {
	return &BusinessError{ComputeError, 9301, "动作数据注入异常", nil, err}
}

func NoCalcNodeError() *BusinessError {
	return &BusinessError{ComputeError, 9302, "没有找到对应的计算节点", nil, nil}
}
