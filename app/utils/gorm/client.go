package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	"wagner/app/global/variable"
)

func GetOneMysqlClient() (*gorm.DB, error) {
	return getSqlDriver()
}

func getSqlDriver() (*gorm.DB, error) {
	var dbDialector gorm.Dialector
	if val, err := getDbDialector(); err != nil {
		//variable.ZapLog.Error(my_errors.ErrorsDialectorDbInitFail+sqlType, zap.Error(err))
	} else {
		dbDialector = val
	}
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		//gorm 数据库驱动初始化失败
		return nil, err
	}

	//// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	//// https://github.com/go-gorm/gorm/issues/3789  此 issue 所反映的问题就是我们本次解决掉的
	//_ = gormDb.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)
	//
	//// https://github.com/go-gorm/gorm/issues/4838
	//_ = gormDb.Callback().Create().Before("gorm:before_create").Register("CreateBeforeHook", CreateBeforeHook)
	//// 为了完美支持gorm的一系列回调函数
	//_ = gormDb.Callback().Update().Before("gorm:before_update").Register("UpdateBeforeHook", UpdateBeforeHook)

	// 为主连接设置连接池(43行返回的数据库驱动指针)
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(variable.OrmConfig.GetDuration("Gorm.mysql.SetConnMaxLifetime") * time.Second)
		rawDb.SetMaxIdleConns(variable.OrmConfig.GetInt("Gorm.mysql.Write.SetMaxIdleConns"))
		rawDb.SetMaxOpenConns(variable.OrmConfig.GetInt("Gorm.mysql.SetMaxOpenConns"))
		// 全局sql的debug配置
		if variable.OrmConfig.GetBool("Gorm.SqlDebug") {
			return gormDb.Debug(), nil
		} else {
			return gormDb, nil
		}
	}
}

// 获取一个数据库方言(Dialector),通俗的说就是根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector() (gorm.Dialector, error) {
	dsn := getDsn()
	dbDialector := mysql.Open(dsn)
	return dbDialector, nil
}

// 根据配置参数生成数据库驱动 dsn
func getDsn() string {
	Host := variable.OrmConfig.GetString("Gorm.mysql.Host")
	DataBase := variable.OrmConfig.GetString("Gorm.mysql.DataBase")
	Port := variable.OrmConfig.GetInt("Gorm.mysql.Port")
	User := variable.OrmConfig.GetString("Gorm.mysql.User")
	Pass := variable.OrmConfig.GetString("Gorm.mysql.Pass")
	Charset := variable.OrmConfig.GetString("Gorm.mysql.Charset")

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=Local", User, Pass, Host, Port, DataBase, Charset)
}
