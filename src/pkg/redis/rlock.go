package redis

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 可重入加锁
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	// watch dog续约
	watchDogCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
	return 0
end
`
	randomLen = 16
	// 默认超时时间，防止死锁
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

// A RLock is a redis lock.
type RLock struct {
	// redis客户端
	store redis.UniversalClient
	// 超时时间
	seconds uint32
	// 锁key
	key string
	// 锁value，防止锁被别人获取到
	id string
	// watch dog stop channel
	stop chan bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRLock returns a RedisLock.
func NewRLock(store redis.UniversalClient, key string) *RLock {
	return &RLock{
		store: store,
		key:   key,
		id:    randomStr(randomLen),
		stop:  make(chan bool, 1),
	}
}

// Acquire acquires the lock.
// 加锁
func (rl *RLock) Acquire() (bool, error) {
	// 获取过期时间
	seconds := atomic.LoadUint32(&rl.seconds)
	// 默认锁过期时间为500ms，防止死锁
	resp, err := rl.store.Eval(lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		// watch dog 续约
		go rl.WatchDog()
		return true, nil
	}
	log.Errorf("Unknown reply when acquiring lock for %s: %+v", rl.key, resp)
	return false, nil
}

// Release releases the lock.
// 释放锁
func (rl *RLock) Release() (bool, error) {
	defer func() {
		rl.stop <- true
	}()
	resp, err := rl.store.Eval(delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	log.Infof("release lock success!!!")
	return reply == 1, nil
}

// SetExpire sets the expire.
// 需要注意的是需要在Acquire()之前调用
// 不然默认为500ms自动释放
func (rl *RLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func (rl *RLock) WatchDog() {
	ticker := time.NewTicker(time.Millisecond * tolerance)
	defer ticker.Stop()
	for {
		select {
		case <-rl.stop:
			log.Infof("watchdog release success!!!")
			return
		case <-ticker.C:
			resp, err := rl.store.Eval(watchDogCommand, []string{rl.key}, []string{
				rl.id, strconv.Itoa(int(rl.seconds) * millisPerSecond),
			}).Result()
			if err != nil {
				log.Errorf("WatchDog watchdog failed,err:%+v", err)
			}
			log.Debugf("WatchDog redis lock!!! resp is %+v", resp)
		}
	}
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
