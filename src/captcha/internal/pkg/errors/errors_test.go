package errors

import (
	"fmt"
	"testing"

	perrors "github.com/hololee2cn/captcha/internal/pkg/errors/internal/errors"
	"github.com/hololee2cn/captcha/internal/pkg/errors/internal/xerrors"
	"github.com/stretchr/testify/assert"
)

func Test_isPkgError(t *testing.T) {
	err := fmt.Errorf("abc")
	assert.Equal(t, false, isPkgError(err))

	err = xerrors.Errorf("abc%v", "sss")
	assert.Equal(t, false, isPkgError(err))

	err = perrors.Errorf("abc%v", "bcd")
	assert.Equal(t, true, isPkgError(err))

	err = Errorf("abc")
	assert.Equal(t, true, isPkgError(err))

	err1 := Errorf("a : %w", err)
	assert.Equal(t, true, isPkgError(err1))

	err1 = Errorf("a : %v", err)
	assert.Equal(t, true, isPkgError(err1))

	err1 = Errorf("a : %s", err)
	assert.Equal(t, true, isPkgError(err1))

	err1 = Errorf("a : %+v", err)
	assert.Equal(t, true, isPkgError(err1))

	err1 = Wrap(err, "abc")
	assert.Equal(t, true, isPkgError(err1))

	err = fmt.Errorf("abc: %w", err)
	assert.Equal(t, true, isPkgError(err))
}

func Test_isXErrorsUnwrapError(t *testing.T) {
	err := xerrors.Errorf("abc")
	assert.Equal(t, true, isXErrorsUnwrapError(err))

	err1 := xerrors.Errorf("sdklfj%s", err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfljk%v", err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%w%w", err, err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%w: %w", err, err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%s%w", err, err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%s: %w", err, err)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err2 := perrors.New("abc")

	err1 = xerrors.Errorf("sdklfj%s", err2)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfljk%v", err2)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%w%w", err2, err2)
	assert.Equal(t, true, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%w: %w", err2, err2)
	assert.Equal(t, false, isXErrorsUnwrapError(err1))

	err1 = xerrors.Errorf("aslkdfl%s%w", err2, err2)
	c := isXErrorsUnwrapError(err1)
	assert.Equal(t, false, c)

	err1 = xerrors.Errorf("aslkdfl%s: %w", err2, err2)
	assert.Equal(t, false, isXErrorsUnwrapError(err1))
}

func Test_getXErrorsUnwrapErrErr(t *testing.T) {
	e := xerrors.Errorf("%s: %s", "abc", "xxx")
	x := getXErrorsUnwrapErrErr(e)
	assert.Equal(t, nil, x)

	e = perrors.New("abc")
	x = getXErrorsUnwrapErrErr(e)
	assert.Equal(t, nil, x)

	e = xerrors.Errorf("xxx: %s", e)
	x = getXErrorsUnwrapErrErr(e)
	assert.NotEqual(t, nil, x)

	e = xerrors.Errorf("xxx: %w", e)
	x = getXErrorsUnwrapErrErr(e)
	assert.NotEqual(t, nil, x)

	e = xerrors.Errorf("xxx: %w", perrors.New("abc"))
	x = getXErrorsUnwrapErrErr(e)
	assert.Equal(t, nil, x)

	e = xerrors.Errorf("xxx: %w", "abc")
	x = getXErrorsUnwrapErrErr(e)
	assert.Equal(t, nil, x)
}
