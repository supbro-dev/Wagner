package container

import (
	"github.com/dgraph-io/ristretto"
	"sync"
)

// 定义一个全局键值对存储容器
var sMap sync.Map

// getOrCreateContainer 创建一个容器工厂
func GetOrCreateContainer(cacheKey string) *container {
	value, _ := sMap.Load(cacheKey)
	if value == nil {
		theCache, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,     // 键追踪数量（通常设为最大缓存的10倍）
			MaxCost:     1 << 30, // 最大缓存成本（例如 1GB）
			BufferItems: 64,      // 性能优化参数
		})
		if err != nil {
			panic(err)
		}
		value = &container{cache: theCache}
		sMap.Store(cacheKey, value)
		return value.(*container)
	}

	if v, ok := value.(*container); ok {
		return v
	} else {
		panic("contains类型有误")
	}
}

type container struct {
	cache *ristretto.Cache
}

// Set  1.以键值对的形式将代码注册到容器
func (c *container) Set(key string, value interface{}) (res bool) {

	if _, exists := c.KeyIsExists(key); exists == false {
		c.cache.Set(key, value, 1)
		res = true
	} else {

	}
	return
}

// Delete  2.删除
func (c *container) Delete(key string) {
	c.cache.Del(key)
}

// Get 3.传递键，从容器获取值
func (c *container) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

// KeyIsExists 4. 判断键是否被注册
func (c *container) KeyIsExists(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *container) ClearCache() {
	c.cache.Clear()
}
