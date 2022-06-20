package store

import (
	"strings"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/utils"
	me "github.com/hololee2cn/wxpub/v1/src/captcha/internal/errors"
	ce "github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors"
	"github.com/hololee2cn/wxpub/v1/src/pkg/redis"

	log "github.com/sirupsen/logrus"
)

// 使用redis而不是基于内存来存储验证码相关信息
type redisStore struct {
	maxAge int64 // 验证码存储时长, 单位为秒
}

func NewRedisStore(maxAge int64) *redisStore {
	return &redisStore{
		maxAge: maxAge,
	}
}

//
// https://github.com/mojocn/base64Captcha/blob/master/interface_store.go

// func (rs *redisStore) Set(id string, value string) {
// 	var err error

// 	if id = strings.TrimSpace(id); id == "" {
// 		err = myErrors.EmptyID
// 		return
// 	}

// 	if value = strings.TrimSpace(value); value == "" {
// 		err = myErrors.EmptyValue
// 		return
// 	}

// 	rc := redis.RedisClient.Get()
// 	if err = rc.Err(); err != nil {
// 		return
// 	}

// 	defer func() {
// 		if err1 := rc.Close(); err1 != nil {
// 			panic(err)
// 		}
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 	}()

// 	_, err = redis.String(rc.Do(redis.RC_SET, id, value, "EX", rs.maxAge))
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	return
// }

// 返回error, 没有实现https://github.com/mojocn/base64Captcha/blob/master/interface_store.go
func (rs *redisStore) Set(id string, value string) (err error) {
	if id = strings.TrimSpace(id); id == "" {
		err = me.InvalidParams("empty id")
		log.Error(err)
		return
	}

	if value = strings.TrimSpace(value); value == "" {
		err = me.InvalidParams("empty value")
		log.Error(err)
		return
	}

	id = utils.FullID(id)
	expire := time.Duration(rs.maxAge) * time.Second

	log.Debugf("redis set key: %v, val: %v, expire: %v", id, value, expire)
	rStatus := redis.RClient.Set(id, value, expire)
	if err = rStatus.Err(); err != nil {
		err = ce.Wrap(err, "set", me.CodeRedisOp)
		log.Error(err)
		return
	}

	return
}

func (rs *redisStore) Get(id string, clear bool) (plain string, err error) {
	if id = strings.TrimSpace(id); id == "" {
		err = me.InvalidParams("empty id")
		log.Error(err)
		return
	}

	id = utils.FullID(id)
	// 如果获取之后就要清除
	if clear {
		defer func() {
			var deleted int64
			rInt := redis.RClient.Del(id)
			if err = rInt.Err(); err != nil {
				log.Error(err)
				return
			}

			deleted, _ = rInt.Result()
			if deleted != 1 {
				log.Warnf("key: %s already deleted or expired", id)
			} else {
				log.Infof("delete key %s success", id)
			}
		}()
	}

	rStr := redis.RClient.Get(id)
	log.Debugf("redis get key: %v, res: %v", id, rStr)
	if err = rStr.Err(); err != nil {
		log.Error(err)
		return
	}
	plain, _ = rStr.Result()
	log.Infof("id: %s, plain: %s", id, plain)
	return
}

func (rs *redisStore) Verify(id, answer string, clear bool) (match bool, err error) {
	var tmpStr string
	tmpStr, err = rs.Get(id, clear)
	if err != nil {
		log.Error(err)
		return
	}
	match = strings.ToLower(tmpStr) == strings.ToLower(answer)
	return
}
