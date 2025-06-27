package variable

import (
	"os"
	"strings"
	"wagner/app/global/business_error"
	"wagner/app/utils/log"
	"wagner/app/utils/yml_config"
)

var (
	BasePath  string
	Config    yml_config.YmlConfigInterf // 全局配置文件指针
	OrmConfig yml_config.YmlConfigInterf // 数据库配置文件
)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.LogBusinessError(business_error.ServerOccurredError(business_error.OsError, err))
	}

}
