package tasks

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/webapi/config"

	"github.com/hololee2cn/wxpub/v1/src/pkg/redis"
	"github.com/hololee2cn/wxpub/v1/src/webapi/g"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	log "github.com/sirupsen/logrus"
)

const (
	MsgInterval     = 10
	TmplInterval    = config.RedisTmplTTL - config.DefaultHTTPTimeOut
	FailMsgInterval = 1
)

var RLock uint32

func CronTasks(ctx context.Context) {
	// 定时消费消息
	AddJob(ctx, time.Second*time.Duration(MsgInterval), HandleMsg, func() { close(maxMsgChan) })
	// 定时更新模板内容
	AddJob(ctx, time.Second*time.Duration(TmplInterval), SyncTemplates, nil)
	// 定时更新失败消息
	AddJob(ctx, time.Second*time.Duration(FailMsgInterval), HandleFailSendMsg, nil)
}

func AddJob(ctx context.Context, interval time.Duration, runner func(ctx context.Context), closeFunc func()) {
	g.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("cron recovered: %+v", r)
			}
			g.Done()

			if ctx.Err() == nil {
				AddJob(ctx, interval, runner, closeFunc) // try not stop
			}
		}()
		// redis lock for process to avoid racing from redis
		rLock := redis.NewRLock(*persistence.DefaultAkRepo().Redis, config.RedisLockTask)
		// init redis lock time 2 seconds
		rLock.SetExpire(2)
		var ok bool
		var e error
		defer func() {
			for i := 0; i < 5; i++ {
				ok, e = rLock.Release()
				if !ok || e != nil {
					time.Sleep(time.Millisecond * 100)
					continue
				}
				return
			}
		}()
		t := time.NewTimer(interval)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				closeFunc()
				return
			case <-t.C:
				if atomic.LoadUint32(&RLock) == 0 {
					// 尝试获取锁
					for i := 0; i < 5; i++ {
						ok, e = rLock.Acquire()
						if !ok || e != nil {
							time.Sleep(time.Millisecond * 100)
							continue
						}
						break
					}
					if !ok {
						break
					}
					atomic.StoreUint32(&RLock, 1)
				}
				runner(ctx)
				t.Reset(interval)
			}
		}
	}()
}
