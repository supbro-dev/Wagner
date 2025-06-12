/*
* @Author: supbro
* @Date:   2025/6/2 10:48
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/2 10:48
 */
package container

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"sync"
)

// 定义一个全局键值对存储容器
var sMap sync.Map

// GenericCache 泛型缓存包装器
type GenericCache[K comparable, V any] struct {
	name  ContainerKey
	cache *ristretto.Cache
}

type ContainerKey string

var (
	CONFIG           ContainerKey = "config"
	DYNAMIC_PARAM    ContainerKey = "dynamic_param"
	HOUR_SUMMARY_MD5 ContainerKey = "hour_sum_md5"
)

// 根据cacheName创建或获取一个泛型缓存
// Parameters: 缓存名称
// Returns: 缓存对象
func GetOrCreateCache[K comparable, V any](cacheName ContainerKey) (*GenericCache[K, V], error) {
	// 检查是否已存在该名称的缓存
	if val, ok := sMap.Load(cacheName); ok {
		if cache, ok := val.(*GenericCache[K, V]); ok {
			return cache, nil
		}
		// 类型不匹配时返回错误
		return nil, fmt.Errorf("cache '%s' already exists with different type", cacheName)
	}

	// 创建新的缓存实例
	cache, err := newGenericCache[K, V](cacheName)
	if err != nil {
		return nil, err
	}

	// 存储并返回新创建的缓存
	sMap.Store(cacheName, cache)
	return cache, nil
}

func newGenericCache[K comparable, V any](key ContainerKey) (*GenericCache[K, V], error) {
	var maxEntries int64 = 100
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxEntries * 10, // 建议为最大键数的10倍
		MaxCost:     maxEntries,      // 最大容量（按条目数）
		BufferItems: 64,              // 建议值
	})
	if err != nil {
		return nil, err
	}

	return &GenericCache[K, V]{
		name:  key,
		cache: cache,
	}, nil
}

// Set  1.以键值对的形式将代码注册到容器
func (c *GenericCache[K, V]) Set(key K, value V) (res bool) {
	if _, exists := c.KeyIsExists(key); exists == false {
		return c.cache.Set(key, value, 1)
	} else {
		return c.cache.Set(key, value, 1)
	}
}

// Delete  2.删除
func (c *GenericCache[K, V]) Delete(key string) {
	c.cache.Del(key)
}

// Get 3.传递键，从容器获取值
func (c *GenericCache[K, V]) Get(key K) (V, bool) {
	return c.KeyIsExists(key)
}

// KeyIsExists 4. 判断键是否被注册
func (c *GenericCache[K, V]) KeyIsExists(key K) (V, bool) {
	if get, b := c.cache.Get(key); b {
		if v, ok := get.(V); ok {
			return v, ok
		} else {
			var zero V
			return zero, false
		}
	} else {
		var zero V
		return zero, false
	}
}

func (c *GenericCache[K, V]) ClearCache() {
	c.cache.Clear()
}
