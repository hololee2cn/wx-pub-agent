package persistence

import (
	"context"
	"time"

	"github.com/hololee2cn/pkg/ginx"

	redis3 "github.com/hololee2cn/wxpub/v1/src/pkg/redis"

	"github.com/go-redis/redis/v7"
	"github.com/hololee2cn/wxpub/v1/src/consts"
	log "github.com/sirupsen/logrus"
)

type WxRepo struct {
	Redis *redis.UniversalClient
}

var defaultWxRepo *WxRepo

func NewWxRepo() {
	if defaultWxRepo == nil {
		defaultWxRepo = &WxRepo{
			Redis: CommonRepositories.Redis,
		}
	}
}

func DefaultWxRepo() *WxRepo {
	return defaultWxRepo
}

func (a *WxRepo) SetMsgIDToRedis(ctx context.Context, msgID string) error {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("SetMsgIDToRedis traceID:%s", traceID)
	var err error
	for i := 0; i < 3; i++ {
		err = redis3.RSet(consts.RedisKeyMsgID+msgID, "", consts.RedisMsgIDTTL)
		if err != nil {
			log.Errorf("SetMsgIDToRedis WxRepo redis set msg id failed,traceID:%s,err:%+v", traceID, err)
			continue
		}
		break
	}
	return err
}

func (a *WxRepo) IsExistMsgIDFromRedis(ctx context.Context, msgID string) (bool, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("IsExistMsgIDFromRedis traceID:%s", traceID)
	var err error
	for i := 0; i < 3; i++ {
		_, err = redis3.RGet(consts.RedisKeyMsgID + msgID)
		if err != nil {
			if err == redis.Nil {
				return false, nil
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		log.Errorf("IsExistMsgIDFromRedis get wx msg id from redis failed,traceID:%s, err:%+v", traceID, err)
		return false, err
	}
	return true, nil
}
