package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wagner/app/global/business_error"
)

var successCode = 0
var successMsg = "success"

func ReturnSuccessEmptyJson(context *gin.Context) {
	ReturnJson(context, http.StatusOK, successCode, successMsg, nil)
}

func ReturnSuccessJson(context *gin.Context, data interface{}) {
	ReturnJson(context, http.StatusOK, successCode, successMsg, data)
}

func ReturnJson(context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

// ErrorSystem 系统执行代码错误
func ReturnError(c *gin.Context, error *business_error.BusinessError) {
	var msg string
	if error.Args != nil && len(error.Args) > 0 {
		msg = fmt.Sprintf(error.Message, error.Args)
	} else {
		msg = error.Message
	}
	ReturnJson(c, http.StatusInternalServerError, error.Code, msg, nil)
	c.Abort()
}
