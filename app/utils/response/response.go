package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wagner/app/global/my_error"
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
func ErrorSystem(c *gin.Context, errorCode int, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, errorCode, my_error.ServerOccurredErrorWithMsg+msg, data)
	c.Abort()
}
