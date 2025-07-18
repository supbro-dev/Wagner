package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wagner/app/global/variable"
	"wagner/app/http/controller"
	"wagner/app/utils/gin_release"
)

func InitRouter() *gin.Engine {
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

	//router.GET("/", func(context *gin.Context) {
	//	context.String(http.StatusOK, "Api 模块接口 hello Wagner！")
	//})

	//  创建一个门户类接口路由组
	vApi := router.Group("/api/v1/")
	{
		employee := vApi.Group("/employee")
		employeeHandler := controller.EmployeeHandler{}
		{
			employee.GET("findByInfo", employeeHandler.FindByInfo)
		}

		workplace := vApi.Group("workplace/")
		workplaceHandler := controller.WorkplaceHandler{}
		{
			workplace.GET("all", workplaceHandler.FindAll)
			workplace.GET("allIndustry", workplaceHandler.FindAllIndustry)
			workplace.GET("allSubIndustry", workplaceHandler.FindAllSubIndustry)
			workplace.GET("findWorkplaceByCode", workplaceHandler.FindWorkplaceByCode)
		}
		workGroup := vApi.Group("workGroup/")
		workGroupHandler := controller.WorkGroupHandler{}
		{
			workGroup.GET("findByCode", workGroupHandler.FindByCode)
		}

		process := vApi.Group("process/")
		processHandler := controller.ProcessHandler{}
		{
			process.GET("implementation", processHandler.Implementation)
			process.POST("saveImplementation", processHandler.SaveImplementation)
			process.GET("getImplementationById", processHandler.GetImplementationById)
			process.GET("getProcessPositionTree", processHandler.GetProcessPositionTree)
			process.GET("findProcessByParentProcessCode", processHandler.FindProcessByParentProcessCode)
			process.GET("generateProcessCode", processHandler.GenerateProcessCode)
			process.POST("saveProcessPosition", processHandler.SaveProcessPosition)
			process.POST("saveProcess", processHandler.SaveProcess)
			process.POST("deleteProcessPosition", processHandler.DeleteProcessPosition)
			process.POST("changeImplStatus", processHandler.ChangeImplStatus)
			process.GET("getWorkplaceStructureTree", processHandler.GetWorkplaceStructureTree)
		}
		position := vApi.Group("position/")
		{
			position.GET("findAll", controller.PositionHandler{}.FindAll)
		}
		efficiency := vApi.Group("efficiency/")
		efficiencyHandler := controller.EfficiencyHandler{}
		{
			efficiency.GET("computeEmployee", efficiencyHandler.ComputeEmployee)
			efficiency.GET("computeWorkplace", efficiencyHandler.ComputeWorkplace)
			efficiency.GET("employee", efficiencyHandler.EmployeeEfficiency)
			efficiency.GET("timeOnTask", efficiencyHandler.TimeOnTask)
			efficiency.GET("workplace", efficiencyHandler.WorkplaceEfficiency)
			efficiency.GET("employeeStatus", efficiencyHandler.EmployeeStatus)
			efficiency.POST("saveOtherParams", efficiencyHandler.SaveOtherParams)
			efficiency.GET("findCalcParamByImplementationId", efficiencyHandler.FindCalcParamByImplementationId)
		}
	}

	// 前端静态文件加载
	// 服务 React 静态文件
	reactStaticPath := filepath.Join("static", "web")
	router.Static("/web", reactStaticPath)

	// 配置主页路由
	router.GET("/", func(c *gin.Context) {
		// 重定向到 React 应用的入口
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	// 处理前端路由：所有未匹配的路由都返回 React 的 index.html
	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(reactStaticPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			c.String(http.StatusNotFound, "Page not found")
		}
	})

	// 处理前端路由的中间件
	router.Use(func(c *gin.Context) {
		// 排除静态文件和 API 请求
		if strings.HasPrefix(c.Request.URL.Path, "/static") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/v1") {
			c.Next()
			return
		}

		// 返回 index.html
		c.File(filepath.Join("./static", "index.html"))
		c.Abort() // 终止后续处理
	})

	return router
}
