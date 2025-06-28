/*
* @Author: supbro
* @Date:   2025/6/12 15:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 15:29
 */
package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"wagner/app/utils/log"

	"time"
	"wagner/app/global/business_error"
	"wagner/app/global/container"
	"wagner/app/global/error_handler"
)

type HourSummaryCheckCache struct {
	localCache  *container.GenericCache[string, string]
	remoteCache *redis.Client
}

func CreateHourSummaryCheckLocalCache() *HourSummaryCheckCache {
	// 设置1000的成本上线，超过1000后按LRU淘汰
	var maxCost int64 = 1000
	if cache, err := container.GetOrCreateCacheWithMaxCost[string, string](container.HOUR_SUMMARY_MD5, maxCost); err != nil {
		error_handler.LogAndPanic(business_error.ServerErrorCausedBy(err))
		return nil
	} else {
		return &HourSummaryCheckCache{localCache: cache}
	}
}

func CreateHourSummaryCheckRemoteCache(redisAddr, password string) *HourSummaryCheckCache {
	// 创建 Redis 客户端
	remoteCache := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})

	return &HourSummaryCheckCache{remoteCache: remoteCache}
}

const keyFmt = "%v-%v-%v"

func (c *HourSummaryCheckCache) PutResultMd5(employeeNumber, workplaceCode string, operateDay time.Time, md5 string) bool {
	key := fmt.Sprintf(keyFmt, employeeNumber, workplaceCode, operateDay)
	if c.localCache != nil {
		return c.localCache.Set(key, md5)
	} else {
		// 设置键值对，72小时后过期
		err := c.remoteCache.Set(context.Background(), key, md5, 72*time.Hour).Err()
		if err != nil {
			log.LogBusinessError(business_error.SetToRedisError(err))
			return false
		} else {
			return true
		}
	}
}

func (c *HourSummaryCheckCache) GutResultMd5(employeeNumber, workplaceCode string, operateDay time.Time) (string, bool) {
	key := fmt.Sprintf(keyFmt, employeeNumber, workplaceCode, operateDay)
	if c.localCache != nil {
		return c.localCache.Get(key)
	} else {
		if result, err := c.remoteCache.Get(context.Background(), key).Result(); err == nil {
			return result, true
		} else {
			return "", false
		}
	}
}
