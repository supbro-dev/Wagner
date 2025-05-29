package main

import (
	"wagner/app/global/variable"
	"wagner/app/utils/startup_banner"
	_ "wagner/bootstrap"
	"wagner/routers"
)

func main() {
	router := routers.InitApiRouter()
	startup_banner.Run()
	_ = router.Run(variable.Config.GetString("HttpServer.Api.Port"))
}
