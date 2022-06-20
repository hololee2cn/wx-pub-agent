// +build go1.13

package errors

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors/internal/errors"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors/internal/xerrors"
)

type WithCode interface {
	error
	Code() int
}

// 请依据xerrors实现来写format
func Errorf(format string, args ...interface{}) (err error) {
	err = errors.Errorf(format, args...)
	if isXErrorsUnwrapError(err) { // 说明wrap失败，请依据xerrors实现来检查语法
		_, format2, _ := xerrors.ParsePercentW(format)
		errCore := getXErrorsUnwrapErrErr(err)
		if errCore != nil {
			err = errors.Wrapf(errCore, format2, args...)
			return
		}
		err = errors.Errorf(format2, args...)
		return
	}
	return
}

// 大部分逻辑与上个函数一致，不复用因为取栈层数会发生变化
func ErrorfWithCode(code int, format string, args ...interface{}) (err error) {
	err = errors.Errorf(format, args...)
	if isXErrorsUnwrapError(err) { // 说明wrap失败，请依据xerrors实现来检查语法
		_, format2, _ := xerrors.ParsePercentW(format)
		errCore := getXErrorsUnwrapErrErr(err)
		if errCore != nil {
			err = errors.Wrapf(errCore, format2, args...)
			return withCode(err, code)
		}
		err = errors.Errorf(format2, args...)
		return withCode(err, code)
	}
	return withCode(err, code)
}
func Wrap(err error, msg string, code ...int) error {
	if err == nil {
		return nil
	}
	if isPkgError(err) {
		err = xerrors.Errorf("%v ->: %w", msg, err)
	} else {
		err = errors.Wrap(err, fmt.Sprintf("%v ->", msg))
	}
	if len(code) > 0 {
		err = withCode(err, code[0])
	}
	return err
}
func New(text string, code ...int) error {
	err := errors.New(text)
	if len(code) > 0 {
		err = withCode(err, code[0])
	}
	return err
}

func withCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &errWithCode{
		error: err,
		code:  code,
	}
}

type errWithCode struct {
	error
	code int
}

func (e errWithCode) Code() int {
	return e.code
}
func (e errWithCode) Unwrap() error {
	return e.error
}
func (e errWithCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		_, _ = io.WriteString(s, fmt.Sprintf("code: %v, ", e.code))
	}

	if f, ok := e.error.(fmt.Formatter); ok {
		f.Format(s, verb)
	} else {
		_, _ = io.WriteString(s, e.error.Error())
	}
}
func (e errWithCode) Error() string {
	return fmt.Sprintf("%v", e)
}

const recursiveDepth = 10
const pkgErrImportPath = "github.com/hololee2cn/captcha/internal/pkg/errors/internal/errors"
const typeXErrorsUnwrapError = "xerrors.noWrapError"

func isPkgError(err error) bool {
	e := err
	for {
		v := reflect.ValueOf(&e)

		for i := 0; (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && i < recursiveDepth; i++ { // 循环直至解出接口/指针后的真实类型(结构体)
			v = v.Elem()
		}
		if v.Type().PkgPath() == pkgErrImportPath {
			return true
		}

		e = Unwrap(e)
		if e == nil {
			return false
		}
	}
}

func isXErrorsUnwrapError(err error) bool {
	e := err
	for {
		v := reflect.ValueOf(&e)
		for i := 0; (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && i < recursiveDepth; i++ {
			v = v.Elem()
		}
		if v.Type().String() == typeXErrorsUnwrapError {
			return true
		}

		e = Unwrap(e)
		if e == nil {
			return false
		}
	}
}

// 有时unwrapError里边其实会包含着一个error，这里是尝试把它取出来
func getXErrorsUnwrapErrErr(err error) error {
	e := err
	var v reflect.Value
	for {
		v = reflect.ValueOf(&e)
		for i := 0; (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && i < recursiveDepth; i++ {
			v = v.Elem()
		}
		if v.Type().String() == typeXErrorsUnwrapError {
			break
		}

		e = Unwrap(e)
		if e == nil {
			return nil
		}
	}
	ev := v.FieldByName("err")
	if ev.IsNil() {
		return nil
	}

	p := ev.InterfaceData()
	c := unsafe.Pointer(&p)
	z := (*error)(c)
	return *z
}
