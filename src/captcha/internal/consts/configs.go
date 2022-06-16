package consts

import (
	"github.com/hololee2cn/wxpub/v1/src/pkg/config"
	"strings"
	"sync"
)

var (
	// 验证码有效时间, 单位秒
	// CaptchaMaxAge = AppConfig.DefaultInt64("captcha_max_age", 180)

	RPCAddr = config.DefaultString("rpc_addr", ":50051")

	LogLevel = config.DefaultInt("log_level", 4)

	RedisAddrs []string
	once       sync.Once
)

func Init() {
	once.Do(func() {
		redisAddrs := config.MustString("redis_addresses")
		if len(redisAddrs) > 0 {
			RedisAddrs = strings.Split(redisAddrs, ",")
		}
	})
}
