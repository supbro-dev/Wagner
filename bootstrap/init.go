package bootstrap

import (
	"wagner/app/global/business_error"
	"wagner/app/global/cache"
	"wagner/app/global/error_handler"
	"wagner/app/global/variable"
	"wagner/app/service"
	"wagner/app/service/action"
	"wagner/app/service/calc/calc_dynamic_param"
	"wagner/app/service/calc/calc_node"
	golang_node2 "wagner/app/service/calc/calc_node/golang_node"
	"wagner/app/service/employee_snapshot"
	"wagner/app/service/process"
	"wagner/app/service/sink"
	"wagner/app/service/workplace"
	"wagner/app/utils/gorm"
	"wagner/app/utils/lock"
	yml_config "wagner/app/utils/yml_config/impl"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/olap_dao"
)

func init() {
	// 4.启动针对配置文件(confgi.yml、gorm_v2.yml)变化的监听， 配置文件操作指针，初始化为全局变量
	variable.Config = yml_config.CreateYamlFactory()
	variable.Config.ConfigFileChangeListen()
	// config>gorm_v2.yml 启动文件变化监听事件
	variable.OrmConfig = yml_config.CreateYamlFactory("gorm")
	variable.OrmConfig.ConfigFileChangeListen()

	lockType := variable.Config.Get("Lock.Type")

	if lockType == "Local" {
		lock.InitLocalLock()
	} else if lockType == "Distributed" {
		lock.InitDistributedLock(variable.Config.GetString("Redis.Addr"), variable.Config.GetString("Redis.Password"))
	}

	hourSummaryCheckCacheType := variable.Config.GetString("Cache.HourSummaryCheck.Type")
	var hourSummaryCheckCache *cache.HourSummaryCheckCache
	if hourSummaryCheckCacheType == "Local" {
		hourSummaryCheckCache = cache.CreateHourSummaryCheckLocalCache()
	} else {
		hourSummaryCheckCache = cache.CreateHourSummaryCheckRemoteCache(variable.Config.GetString("Redis.Addr"), variable.Config.GetString("Redis.Password"))
	}

	client, err := gorm.GetOneMysqlClient()
	if err != nil {
		error_handler.LogAndPanic(business_error.CreateMysqlClientError(err))
	}

	// olap建议读写分别创建客户端
	olapWriteClient, err := gorm.GetOneOlapClient()
	if err != nil {
		error_handler.LogAndPanic(business_error.CreateOlapClientError(err))
	}
	olapReadClient, err := gorm.GetOneOlapClient()
	if err != nil {
		error_handler.LogAndPanic(business_error.CreateOlapClientError(err))
	}

	workplaceDao := dao.CreateWorkplaceDao(client)
	scriptDao := dao.CreateScriptDao(client)
	employeeStatusDao := dao.CreateEmployeeStatusDao(client)

	actionService := action.CreateActionService(dao.CreateActionRepository(client))

	employeeSnapshotService := employee_snapshot.CreateEmployeeSnapshotService(dao.CreateEmployeeDao(client))

	processService := process.CreateProcessServiceImpl(dao.CreateProcessPositionDao(client), dao.CreateProcessImplementDao(client), workplaceDao)

	calcDynamicParamService := calc_dynamic_param.CreateCalcDynamicParamService(dao.CreateCalcDynamicParamDao(client), workplaceDao, scriptDao)

	summarySinkService := sink.CreateSummarySinkService(olap_dao.CreateHourSummaryResultDao(olapWriteClient), hourSummaryCheckCache)

	workplaceService := workplace.CreateWorkplaceService(workplaceDao)

	employeeStatusSinkService := sink.CreateEmployeeStatusSinkService(employeeStatusDao)

	domainServiceHolder := service.DomainServiceHolder{
		EmployeeSnapshotService: employeeSnapshotService,
		ActionService:           actionService,
		ProcessService:          processService,
		CalcDynamicParamService: calcDynamicParamService,
		WorkplaceService:        workplaceService,
	}

	service.DomainHolder = domainServiceHolder

	efficiencyComputeService := service.CreateEfficiencyComputeService()
	efficiencyService := service.CreateEfficiencyService(olap_dao.CreateHourSummaryResultDao(olapReadClient), employeeStatusDao)
	serviceHolder := service.ServiceHolder{
		EfficiencyComputeService:  efficiencyComputeService,
		EfficiencyService:         efficiencyService,
		SummarySinkService:        summarySinkService,
		EmployeeStatusSinkService: employeeStatusSinkService,
	}

	service.Holder = serviceHolder

	// 注册计算节点脚本
	calc_node.Register("SetCrossDayAttendance", golang_node2.SetCrossDayAttendance)
	calc_node.Register("ComputeAttendanceDefaultEndTime", golang_node2.ComputeAttendanceDefaultEndTime)
	calc_node.Register("MarchProcess", golang_node2.MarchProcess)
	calc_node.Register("CutOffAttendanceTime", golang_node2.CutOffAttendanceTime)
	calc_node.Register("AddCrossDayData", golang_node2.AddCrossDayData)
	calc_node.Register("FilterOtherDaysData", golang_node2.FilterOtherDaysData)
	calc_node.Register("FilterExpiredData", golang_node2.FilterExpiredData)
	calc_node.Register("ComputeAttendanceDefaultStartTime", golang_node2.ComputeAttendanceDefaultStartTime)
	calc_node.Register("PaddingUnfinishedWorkEndTime", golang_node2.PaddingUnfinishedWorkEndTime)
	calc_node.Register("CutOffOvertimeWork", golang_node2.CutOffOvertimeWork)
	calc_node.Register("CutOffCrossWork", golang_node2.CutOffCrossWork)
	calc_node.Register("AddReasonableBreakTime", golang_node2.AddReasonableBreakTime)
	calc_node.Register("CutOffWorkByRest", golang_node2.CutOffWorkByRest)
	calc_node.Register("CalcWorkTransitionTime", golang_node2.CalcWorkTransitionTime)
	calc_node.Register("GenerateIdleData", golang_node2.GenerateIdleData)
	calc_node.Register("MatchRestProcess", golang_node2.MatchRestProcess)
}
