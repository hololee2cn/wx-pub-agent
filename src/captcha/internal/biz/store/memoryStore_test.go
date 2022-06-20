package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	const (
		MaxAge = 3
		Key    = "hello"
		Value  = "world"
	)

	s := NewMemoryStore(MaxAge)

	err := s.Set(Key, Value)
	assert.Nil(t, err)

	plain, err := s.Get(Key, false)
	assert.Nil(t, err)
	assert.Equal(t, Value, plain)

	match, err := s.Verify(Key, Value, false)
	assert.Nil(t, err)
	assert.True(t, match)

	time.Sleep(time.Duration(MaxAge) * time.Second)
	plain, err = s.Get(Key, false)
	assert.NotNil(t, err)
	assert.Equal(t, "", plain)
}
