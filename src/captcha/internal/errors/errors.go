package errros

import (
	ce "github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors"
)

// 本模块所有的错误码在此定义
// 6301-6500 shared-captcha(通用的验证码模块)
const (
	CodeNoTraceID = 1001 // 未带trace-id

	CodeInvalidParams = 6301 // 参数错误

	CodeInvalidBgColorFormat = 6302 // 背景颜色格式不对
	CodeFailedVerify         = 6303 // 检查验证码失败
	CodeFailedGenCaptcha     = 6004 // 生成验证码失败
	CodeRedisOp              = 6005 // 　redis操作失败
	CodeNotExist             = 6006 // 不存在

	CodeUnknown = 6500 // 未知错误
)

// 将一些常用的错误封装为函数
// 这里写成函数的形式是为了调试时可以看出调用栈
var (
	NoTraceID     = func(msg ...string) error { return ce.New(format("no x-trace-id", msg...), CodeNoTraceID) }
	InvalidParams = func(msg ...string) error { return ce.New(format("invalid params", msg...), CodeInvalidParams) }
)

func format(msg1 string, msg2 ...string) string {
	if len(msg2) == 0 {
		return msg1
	}
	return msg1 + ": " + msg2[0]
}
