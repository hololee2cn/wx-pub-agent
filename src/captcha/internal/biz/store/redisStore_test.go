package store

import (
	"github.com/hololee2cn/wxpub/v1/src/pkg/config"
	"github.com/hololee2cn/wxpub/v1/src/pkg/redis"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

// 依赖于redis
func TestRedisStore(t *testing.T) {
	redisAddrs := config.DefaultString("redis_addrs", "")
	if len(redisAddrs) == 0 {
		t.SkipNow()
	}
	redis.NewRedisClient(strings.Split(redisAddrs, ","))

	t.Log("use redis")
	const MaxAge = 10
	var id string
	var err error
	var match bool
	var plain string

	rs1 := NewRedisStore(MaxAge)
	_uuid := uuid.New().String()

	values := []string{"world", "世界", "千山鸟飞绝"}
	for i, value := range values {
		id = _uuid + strconv.Itoa(i)
		err = rs1.Set(id, value)
		if err != nil {
			t.Error("set failed")
		}

		match, err = rs1.Verify(id, value, false)
		assert.Nil(t, err)
		assert.True(t, match)

		match, err = rs1.Verify(id, value, true)
		assert.Nil(t, err)
		assert.True(t, match)

		plain, err = rs1.Get(id, true)
		assert.NotNil(t, err)
		assert.Equal(t, "", plain)
	}
}
