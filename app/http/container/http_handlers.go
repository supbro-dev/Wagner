package container

import (
	"github.com/gin-gonic/gin"
	"wagner/app/http/controller"
	"wagner/app/utils/log"
)

var apiCache map[string]func(*gin.Context)

func init() {
	//创建容器
	apiCache = make(map[string]func(*gin.Context))
	apiCache["efficiencyCompute"] = controller.EfficiencyComputeHandler{}.ComputeEmployee
	apiCache["workplace"] = controller.WorkplaceHandler{}.FindAll
	efficiencyHandler := controller.EfficiencyHandler{}
	apiCache["efficiency.employee"] = efficiencyHandler.EmployeeEfficiency

}

func GetHandler(key string) func(context *gin.Context) {
	if value, exists := apiCache[key]; exists {
		return value
	}
	log.SystemLogger.Error("获取API handler异常")
	return nil
}
