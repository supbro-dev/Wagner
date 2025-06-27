/*
* @Author: supbro
* @Date:   2025/6/25 16:50
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/25 16:50
 */
package error_handler

import (
	"wagner/app/global/business_error"
	"wagner/app/utils/log"
)

// Log并Panic,如果业务异常封装了异常来源，直接panic异常来源
func LogAndPanic(err *business_error.BusinessError) {
	log.LogBusinessError(err)
	if err.CausedError != nil {
		panic(err.CausedError)
	} else {
		panic(err)
	}
}
