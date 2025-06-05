package container

import (
	"github.com/gin-gonic/gin"
	"wagner/app/global/container"
	"wagner/app/http/controller"
	"wagner/app/utils/log"
)

type ApiInvoker interface {
	Invoke(context *gin.Context)
}

var apiCache *container.GenericCache[string, ApiInvoker]

func init() {
	//创建容器
	cache, err := container.GetOrCreateCache[string, ApiInvoker](container.API)
	if err != nil {
		panic(err)
	}

	apiCache = cache
	apiCache.Set("efficiencyCompute", controller.EfficiencyComputeHandler{})
}

func GetHandler(key string) func(context *gin.Context) {
	if value := apiCache.Get(key); value != nil {
		if val, isOk := value.(ApiInvoker); isOk {
			return val.Invoke
		}
	}
	log.SystemLogger.Error("获取API handler异常")
	return nil
}
