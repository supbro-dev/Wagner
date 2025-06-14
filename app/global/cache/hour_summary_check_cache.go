/*
* @Author: supbro
* @Date:   2025/6/12 15:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 15:29
 */
package cache

import (
	"fmt"
	"time"
	"wagner/app/global/container"
)

type HourSummaryCheckCache struct {
	localCache  *container.GenericCache[string, string]
	remoteCache interface{}
}

func CreateHourSummaryCheckLocalCache() *HourSummaryCheckCache {
	cache, err := container.GetOrCreateCache[string, string](container.HOUR_SUMMARY_MD5)
	if err != nil {
		panic(err)
	}
	return &HourSummaryCheckCache{localCache: cache}
}

func (c *HourSummaryCheckCache) PutResultMd5(employeeNumber, workplaceCode string, operateDay time.Time, md5 string) bool {
	return c.localCache.Set(fmt.Sprintf("%v-%v-%v", employeeNumber, workplaceCode, operateDay), md5)
}

func (c *HourSummaryCheckCache) GutResultMd5(employeeNumber, workplaceCode string, operateDay time.Time) (string, bool) {
	return c.localCache.Get(fmt.Sprintf("%v-%v-%v", employeeNumber, workplaceCode, operateDay))
}
