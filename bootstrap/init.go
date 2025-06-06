package bootstrap

import (
	"fmt"
	"wagner/app/global/variable"
	"wagner/app/service"
	"wagner/app/service/action"
	"wagner/app/service/calc_dynamic_param"
	"wagner/app/service/employee_snapshot"
	"wagner/app/service/sink"
	"wagner/app/service/standard_position"
	"wagner/app/utils/gorm"
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

	fmt.Println("finished init")

	client, err := gorm.GetOneMysqlClient()
	if err != nil {
		panic(err)
	}

	olapClient, err := gorm.GetOneOlapClient()

	if err != nil {
		panic(err)
	}

	workplaceDao := dao.CreateWorkplaceDao(client)
	scriptDao := dao.CreateScriptDao(client)

	actionService := action.CreateActionService(dao.CreateActionRepository(client))

	employeeSnapshotService := employee_snapshot.CreateEmployeeSnapshotService(dao.CreateEmployeeDao(client))

	standardPositionService := standard_position.CreateStandardPositionService(dao.CreateStandardPositionDao(client), workplaceDao)

	calcDynamicParamService := calc_dynamic_param.CreateCalcDynamicParamService(dao.CreateCalcDynamicParamDao(client), workplaceDao, scriptDao)

	summarySinkService := sink.CreateSummarySinkService(olap_dao.CreateHourSummaryResultDao(olapClient))

	domainServiceHolder := service.DomainServiceHolder{
		EmployeeSnapshotService: employeeSnapshotService,
		ActionService:           actionService,
		StandardPositionService: standardPositionService,
		CalcDynamicParamService: calcDynamicParamService,
	}

	service.DomainHolder = domainServiceHolder

	efficiencyComputeService := service.EfficiencyComputeService{}
	serviceHolder := service.ServiceHolder{
		EfficiencyComputeService: &efficiencyComputeService,
		SummarySinkService:       summarySinkService,
	}

	service.Holder = serviceHolder
}
