package yml_config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
	"time"
	"wagner/app/global/business_error"
	"wagner/app/global/container"
	"wagner/app/global/error_handler"
	"wagner/app/global/variable"
	"wagner/app/utils/yml_config"
)

type ymlConfig struct {
	viper *viper.Viper
	mu    *sync.Mutex
}

var lastChangeTime time.Time
var containerFactory *container.GenericCache[string, interface{}]

func init() {
	lastChangeTime = time.Now()
	cache, err := container.GetOrCreateCache[string, interface{}](container.CONFIG)
	if err != nil {
		error_handler.LogAndPanic(business_error.CreateCacheError(err))
	}
	containerFactory = cache
}

// CreateYamlFactory 创建一个yaml配置文件工厂
// 参数设置为可变参数的文件名，这样参数就可以不需要传递，如果传递了多个，我们只取第一个参数作为配置文件名
func CreateYamlFactory(fileName ...string) yml_config.YmlConfigInterf {

	yamlConfig := viper.New()
	// 配置文件所在目录
	yamlConfig.AddConfigPath(variable.BasePath + "/config")
	// 需要读取的文件名,默认为：config
	if len(fileName) == 0 {
		yamlConfig.SetConfigName("config")
	} else {
		yamlConfig.SetConfigName(fileName[0])
	}
	//设置配置文件类型(后缀)为 yml
	yamlConfig.SetConfigType("yml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		error_handler.LogAndPanic(business_error.ServerOccurredError(business_error.OsError, err))
	}

	return &ymlConfig{
		viper: yamlConfig,
		mu:    new(sync.Mutex),
	}
}

// ConfigFileChangeListen 监听文件变化
func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// keyIsCache 判断相关键是否已经缓存
func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := containerFactory.KeyIsExists(keyName); exists {
		return true
	} else {
		return false
	}
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{}) bool {
	// 避免瞬间缓存键、值时，程序提示键名已经被注册的日志输出
	y.mu.Lock()
	defer y.mu.Unlock()
	if _, exists := containerFactory.KeyIsExists(keyName); exists {
		return true
	}
	return containerFactory.Set(keyName, value)
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	if get, b := containerFactory.Get(keyName); b {
		return get
	} else {
		return nil
	}
}

// 清空已经缓存的配置项信息
func (y *ymlConfig) clearCache() {
	containerFactory.ClearCache()
}

// Get 一个原始值
func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetString 字符串格式返回值
func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

// GetBool 布尔格式返回值
func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt 整数格式返回值
func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32 整数格式返回值
func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64 整数格式返回值
func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetFloat64 小数格式返回值
func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetDuration 时间单位格式返回值
func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetStringSlice 字符串切片数格式返回值
func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}
