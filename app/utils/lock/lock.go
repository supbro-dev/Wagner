/*
* @Author: supbro
* @Date:   2025/5/31 23:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/5/31 23:06
 */
package lock

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
	"wagner/app/global/business_error"
	"wagner/app/global/container"
	"wagner/app/global/error_handler"
	"wagner/app/utils/datetime_util"
)

var lockFmt = "lock:%s:%s"

func Lock(employeeNumber string, operateDay time.Time, expireSec int) (bool, error) {
	lockKey := fmt.Sprintf(lockFmt, employeeNumber, datetime_util.FormatDate(operateDay))
	if t == local {
		return lockLocal(lockKey)
	} else {
		return lockDistributed(lockKey, expireSec)
	}
}

func lockLocal(lockKey string) (bool, error) {
	lock, exists := localLockMap.Get(lockKey)
	if !exists {
		// 创建互斥锁
		var theLock sync.Mutex
		theLock.Lock()
		localLockMap.Set(lockKey, &theLock)
		return true, nil
	} else {
		tryLock := lock.TryLock()
		return tryLock, nil
	}
}

func lockDistributed(lockKey string, expireSec int) (bool, error) {
	theLock, exists := distributedLockMap.Get(lockKey)
	if !exists {
		// 创建互斥锁
		theLock = rs.NewMutex(lockKey,
			redsync.WithExpiry(time.Second*time.Duration(expireSec)), // 锁过期时间
			redsync.WithTries(5),                         // 最大尝试次数
			redsync.WithRetryDelay(500*time.Millisecond), // 重试间隔
		)
	}

	if err := theLock.Lock(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Unlock(employeeNumber string, operateDay time.Time) (bool, error) {
	lockKey := fmt.Sprintf(lockFmt, employeeNumber, datetime_util.FormatDate(operateDay))
	if t == local {
		return unlockLocal(lockKey)
	} else {
		return unlockDistributed(lockKey)
	}
}

func unlockLocal(lockKey string) (bool, error) {
	if lock, exists := localLockMap.Get(lockKey); exists {
		lock.Unlock()
	}
	return true, nil
}

func unlockDistributed(lockKey string) (bool, error) {
	if theLock, exists := distributedLockMap.Get(lockKey); exists {
		unlockRes, err := theLock.Unlock()
		//if err == nil {
		//	distributedLockMap.Delete(lockKey)
		//}

		return unlockRes, err
	} else {
		return true, nil
	}
}

func InitLocalLock() {
	c, err := container.GetOrCreateCacheWithMaxCost[string, *sync.Mutex](container.LOCK, 1000)
	if err != nil {
		error_handler.LogAndPanic(business_error.ServerErrorCausedBy(err))
	}

	localLockMap = c

	t = local
}

var rs *redsync.Redsync
var distributedLockMap *container.GenericCache[string, *redsync.Mutex]
var localLockMap *container.GenericCache[string, *sync.Mutex]
var t lockType

type lockType string

var (
	local       lockType = "Local"
	distributed lockType = "Distributed"
)

func InitDistributedLock(redisAddr string, password string) {
	c, err := container.GetOrCreateCacheWithMaxCost[string, *redsync.Mutex](container.LOCK, 1000)
	if err != nil {
		error_handler.LogAndPanic(business_error.ServerErrorCausedBy(err))
	}

	distributedLockMap = c
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})

	// 创建锁管理器
	pool := goredis.NewPool(client)
	rs = redsync.New(pool)

	t = distributed
}
