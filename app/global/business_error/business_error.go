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
	OsError                    string = "操作系统异常"
	ParamError                 string = "传参异常"
	HttpInvokeError            string = "HTTP调用异常"
	SystemError                string = "系统异常"
	LockHandleError            string = "锁异常"
	ComputeError               string = "计算异常"
	CacheError                 string = "缓存异常"
	BasicDataError             string = "基础数据异常"
	MysqlError                 string = "Mysql异常"
	DaoError                   string = "DAO异常"
	ProcessImplementationError string = "环节实施异常"
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

func SubmitDataIsWrong(err error) *BusinessError {
	return &BusinessError{ParamError, 9105, "参数提交有误", nil, err}
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

func ParseCalcParamError(err error) *BusinessError {
	return &BusinessError{ParamError, 9303, "解析计算参数异常", nil, err}
}

func CannotFindCalcParamByWorkplace(workplaceCode interface{}) *BusinessError {
	return &BusinessError{ParamError, 9304, "根据工作点%v查找不到计算参数", []interface{}{workplaceCode}, nil}
}

// CACHE
func CreateCacheError(err error) *BusinessError {
	return &BusinessError{CacheError, 9401, "创建缓存失败", nil, err}
}

func SetToRedisError(err error) *BusinessError {
	return &BusinessError{CacheError, 9402, "Redis缓存Set失败", nil, err}
}

// BASIC_DATA
func WorkplaceDoseNotExist(code interface{}) *BusinessError {
	return &BusinessError{BasicDataError, 9501, "工作点%v不存在", []interface{}{code}, nil}
}

// MYSQL
func CreateMysqlClientError(err error) *BusinessError {
	return &BusinessError{MysqlError, 9601, "创建Mysql客户端异常", nil, err}
}

func CreateOlapClientError(err error) *BusinessError {
	return &BusinessError{MysqlError, 9602, "创建Olap客户端异常", nil, err}
}

// DAO
func ReflectSetDataError(err error) *BusinessError {
	return &BusinessError{DaoError, 9701, "通过反射设置数据异常", nil, err}
}

func UnsupportedFieldTypeError() *BusinessError {
	return &BusinessError{DaoError, 9702, "不支持的属性类型", nil, nil}
}

// 环节
func ProcessTargetTypeError() *BusinessError {
	return &BusinessError{ProcessImplementationError, 9801, "环节实施类型异常", nil, nil}
}

func ExistSameCodeProcessImpl(code ...interface{}) *BusinessError {
	return &BusinessError{ProcessImplementationError, 9802, "环节实施已存在:%v", code, nil}
}

func ExistSameCodeProcessPosition(code ...interface{}) *BusinessError {
	return &BusinessError{ProcessImplementationError, 9803, "不能添加相同的环节:%v", code, nil}
}
