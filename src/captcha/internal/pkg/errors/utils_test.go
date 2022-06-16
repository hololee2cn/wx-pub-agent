package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	err := ErrorfWithCode(123, "sdlkfj%v", "abc")
	assert.Equal(t, 123, Code(err))

	err1 := Errorf("sdlkfj: %w", err)
	assert.Equal(t, 123, Code(err1))

	err1 = Wrap(err, "abc")
	assert.Equal(t, 123, Code(err1))

	err1 = Wrap(err, "abc", 234)
	assert.Equal(t, 234, Code(err1))

	err1 = Errorf("sdlkfj: %w", err)
	assert.Equal(t, 123, Code(err1))

	err2 := ErrorfWithCode(345, "sabc: %w", err1)
	assert.Equal(t, 345, Code(err2))

	err3 := Wrap(err2, "abc", 456)
	assert.Equal(t, 456, Code(err3))
}

func TestCodeSource(t *testing.T) {
	err := ErrorfWithCode(123, "sdlkfj%v", "abc")
	assert.Equal(t, 123, CodeSource(err))

	err1 := Errorf("sdlkfj: %w", err)
	assert.Equal(t, 123, CodeSource(err1))

	err1 = Wrap(err, "abc")
	assert.Equal(t, 123, CodeSource(err1))

	err1 = Wrap(err, "abc", 234)
	assert.Equal(t, 123, CodeSource(err1))

	err1 = Errorf("sdlkfj: %w", err)
	assert.Equal(t, 123, CodeSource(err1))

	err2 := ErrorfWithCode(345, "sabc: %w", err1)
	assert.Equal(t, 123, CodeSource(err2))

	err3 := Wrap(err2, "abc", 456)
	assert.Equal(t, 123, CodeSource(err3))
}

func TestParseErrMsg(t *testing.T) {
	code, msg := ParseErrMsg("code: 123, abc")
	assert.Equal(t, 123, code)
	assert.Equal(t, "abc", msg)

	code, msg = ParseErrMsg("code: 12a3, abc")
	assert.Equal(t, -1, code)
	assert.Equal(t, "abc", msg)

	code, msg = ParseErrMsg("code: 123,abc")
	assert.Equal(t, -1, code)
	assert.Equal(t, "code: 123,abc", msg)

	code, msg = ParseErrMsg("code:123, abc")
	assert.Equal(t, -1, code)
	assert.Equal(t, "code:123, abc", msg)

	code, msg = ParseErrMsg("c123, abc")
	assert.Equal(t, -1, code)
	assert.Equal(t, "c123, abc", msg)
}
