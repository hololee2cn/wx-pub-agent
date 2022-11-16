package config

import (
	"strings"
	"sync"
)

var (
	// 验证码有效时间, 单位秒
	// CaptchaMaxAge = AppConfig.DefaultInt64("captcha_max_age", 180)

	RPCAddr = DefaultString("rpc_addr", ":50051")

	LogLevel = DefaultInt("log_level", 4)

	RedisAddrs []string
	once       sync.Once
)

func Init() {
	InitC()
	once.Do(func() {
		redisAddrs := MustString("redis_addresses")
		if len(redisAddrs) > 0 {
			RedisAddrs = strings.Split(redisAddrs, ",")
		}
	})
}
