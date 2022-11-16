package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

var (
	RClient redis.UniversalClient
)

func NewRedisClient(addresses []string) (func(), redis.UniversalClient) {
	RClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addresses,
	})
	err := RClient.Ping().Err()
	if err != nil {
		panic(err)
	}
	return func() {
		_ = RClient.Close()
	}, RClient
}

func RSet(key string, value interface{}, ttl int) (err error) {
	r := RClient.Set(key, value, time.Second*time.Duration(ttl))
	if err = r.Err(); err != nil {
		log.Errorf("failed to RSet,err:%+v", err)
		return
	}
	return
}

func RGet(key string) (value []byte, err error) {
	value, err = RClient.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		log.Errorf("get key from redis failed, key:%s, err:%+v", key, err)
		return
	}
	return
}
