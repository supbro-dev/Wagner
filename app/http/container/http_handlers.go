package container

import (
	"github.com/gin-gonic/gin"
	"wagner/app/global/container"
	"wagner/app/global/my_const"
	"wagner/app/http/controller"
	"wagner/app/utils/log"
)

type ApiInvoker interface {
	Invoke(context *gin.Context)
}

func init() {
	//创建容器
	containers := container.GetOrCreateContainer(my_const.API)

	containers.Set("pprCompute", controller.PprComputeHandler{})
}

func GetHandler(key string) func(context *gin.Context) {
	containers := container.GetOrCreateContainer(my_const.API)
	if value := containers.Get(key); value != nil {
		if val, isOk := value.(ApiInvoker); isOk {
			return val.Invoke
		}
	}
	log.SystemLogger.Error("获取API handler异常")
	return nil
}
