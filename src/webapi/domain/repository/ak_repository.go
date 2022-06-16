package repository

import (
	"context"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"

	"github.com/hololee2cn/pkg/ginx"

	"github.com/hololee2cn/wxpub/v1/src/pkg/redis"

	log "github.com/sirupsen/logrus"
)

type AccessTokenRepository struct {
	ak *persistence.AkRepo
}

var defaultAccessTokenRepository = &AccessTokenRepository{}

func NewAccessTokenRepository(ak *persistence.AkRepo) {
	if defaultAccessTokenRepository.ak == nil {
		defaultAccessTokenRepository = &AccessTokenRepository{
			ak: ak,
		}
	}
}

func DefaultAccessTokenRepository() *AccessTokenRepository {
	return defaultAccessTokenRepository
}

func (a *AccessTokenRepository) GetAccessToken(ctx context.Context) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("GetAccessToken traceID:%s", traceID)
	// 先从redis中取access token，没有则调用接口获取并保存
	var ak string
	var err error
	ak, err = a.ak.GetAccessTokenFromRedis(ctx)
	if err != nil {
		log.Errorf("GetAccessToken AccessTokenRepository GetAccessTokenFromRedis failed,traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	if len(ak) > 0 {
		return ak, nil
	}
	return a.FreshAccessToken(ctx)
}

func (a *AccessTokenRepository) FreshAccessToken(ctx context.Context) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("FreshAccessToken traceID:%s", traceID)
	var err error
	var oldAk string
	// 先获取旧ak
	oldAk, err = a.ak.GetAccessTokenFromRedis(ctx)
	if err != nil {
		log.Errorf("FreshAccessToken AccessTokenRepository GetAccessToken failed,traceID:%s,err:%+v", traceID, err)
	}
	// redis lock for access token to avoid racing to cover ak value from redis
	rLock := redis.NewRLock(*a.ak.Redis, consts.RedisLockAccessToken)
	// init redis lock time 2 seconds
	rLock.SetExpire(2)
	var ok bool
	var e error
	var newAk string
	// 尝试获取锁
	for i := 0; i < 5; i++ {
		ok, e = rLock.Acquire()
		if !ok || e != nil {
			log.Errorf("FreshAccessToken get redis lock failed,traceID:%s, ok:%v,err:%+v", traceID, ok, e)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}
	// 获取锁成功
	if ok {
		// 更新new ak
		newAk, err = func(ctx context.Context) (string, error) {
			defer func() {
				for i := 0; i < 5; i++ {
					ok, err = rLock.Release()
					if !ok || err != nil {
						log.Errorf("getAccessTokenAndReleaseLock delete redis lock failed,traceID:%s, ok:%v,err:%+v", traceID, ok, err)
						time.Sleep(time.Millisecond * 100)
						continue
					}
					return
				}
			}()
			newAk, err = a.getAccessTokenFromRemote(ctx)
			if err != nil {
				log.Errorf("getAccessTokenAndReleaseLock AccessTokenRepository getAccessTokenFromRemote failed,traceID:%s,err:%+v", traceID, err)
				return "", err
			}
			return newAk, nil
		}(ctx)
		return newAk, nil
	}
	// 获取不到锁，休眠100ms再从redis中取当前ak，如果ak value发生改变，证明已更新，判断5次
	for i := 0; i < 5; i++ {
		// 获取新的accessToken
		newAk, err = a.ak.GetAccessTokenFromRedis(ctx)
		if err != nil {
			log.Errorf("GetAccessTokenFromRedis AkRepo GetAccessToken failed,traceID:%s,err:%+v", traceID, err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		// 证明更新成功
		if len(newAk) > 0 && newAk != oldAk {
			return newAk, nil
		}
	}
	return "", err
}

func (a *AccessTokenRepository) getAccessTokenFromRemote(ctx context.Context) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("FreshAccessToken traceID:%s", traceID)
	akResp, err := a.ak.GetAccessTokenFromRequest(ctx)
	if err != nil {
		log.Errorf("getAccessTokenFromRemote AccessTokenRepository GetAccessTokenFromRequest failed,traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	err = a.ak.SetAccessTokenToRedis(ctx, akResp.AccessToken, int(akResp.ExpiresIn))
	if err != nil {
		log.Errorf("getAccessTokenFromRemote AccessTokenRepository  SetAccessTokenToRedis failed,traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	return akResp.AccessToken, nil
}
