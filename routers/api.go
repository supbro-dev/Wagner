package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"wagner/app/global/variable"
	"wagner/app/http/controller"
	"wagner/app/utils/gin_release"
)

func InitApiRouter() *gin.Engine {

	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.Config.GetBool("AppDebug") == false {
		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	_ = router.SetTrustedProxies(nil)

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Api 模块接口 hello Wagner！")
	})

	//处理静态资源（不建议gin框架处理静态资源，参见 Public/readme.md 说明 ）
	router.Static("/public", "./public") //  定义静态资源路由与实际目录映射关系
	//router.StaticFile("/abcd", "./public/readme.md") // 可以根据文件名绑定需要返回的文件名

	//  创建一个门户类接口路由组
	vApi := router.Group("/api/v1/")
	{
		workplace := vApi.Group("workplace/")
		{
			workplace.GET("all", controller.WorkplaceHandler{}.FindAll)
		}
		efficiency := vApi.Group("efficiency/")
		efficiencyHandler := controller.EfficiencyHandler{}
		{
			efficiency.GET("employee", efficiencyHandler.EmployeeEfficiency)
			efficiency.GET("compute", efficiencyHandler.ComputeEmployee)
			efficiency.GET("timeOnTask", efficiencyHandler.TimeOnTask)
			efficiency.GET("workplace", efficiencyHandler.WorkplaceEfficiency)
			efficiency.GET("employeeStatus", efficiencyHandler.EmployeeStatus)
		}
	}
	return router
}
