package store

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// 可以利用此来做单元测试
type memoryStore struct {
	maxAge time.Duration
	c      *cache.Cache
}

// maxAge: second
func NewMemoryStore(maxAge int64) *memoryStore {
	return &memoryStore{
		maxAge: time.Duration(maxAge) * time.Second,
		c:      cache.New(cache.NoExpiration, time.Minute),
	}
}

func (s *memoryStore) Set(id string, value string) error {
	s.c.Set(id, value, s.maxAge)
	return nil
}

func (s *memoryStore) Get(id string, clear bool) (plain string, err error) {
	defer func() {
		if clear {
			s.c.Delete(id)
		}
	}()

	v, found := s.c.Get(id)
	if !found {
		err = fmt.Errorf("no such key: %s", id)
		return
	}

	var ok bool
	plain, ok = v.(string)
	if !ok {
		err = fmt.Errorf("answer not string")
		return
	}

	return
}

func (s *memoryStore) Verify(id, answer string, clear bool) (match bool, err error) {
	plain, err := s.Get(id, clear)
	if err != nil {
		return
	}

	if plain != answer {
		match = false
		return
	}
	match = true
	return
}
