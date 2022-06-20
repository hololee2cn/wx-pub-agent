// +build go1.13

package errors

import (
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 如果未设置错误码则返回-1
func Code(err error) int {
	for ; err != nil; err = Unwrap(err) {
		if e, ok := err.(WithCode); ok {
			return e.Code()
		}
	}
	return -1
}

// 如果未设置错误码则返回-1
func CodeSource(err error) int {
	var last WithCode
	for ; err != nil; err = Unwrap(err) {
		if l, ok := err.(WithCode); ok {
			last = l
		}
	}
	if last != nil {
		return last.Code()
	}

	return -1
}

func UnwrapWithCode(err error) WithCode {
	for ; err != nil; err = Unwrap(err) {
		if e, ok := err.(WithCode); ok {
			return e
		}
	}
	return nil
}

func UnwrapToSource(err error) error {
	var last error
	for ; err != nil; err = Unwrap(err) {
		last = err
	}
	return last
}

// 拿到带着返回码的最深一层错误，并不代表其更深处没有错误
func UnwrapToSourceWithCode(err error) WithCode {
	var last WithCode
	for ; err != nil; err = Unwrap(err) {
		if l, ok := err.(WithCode); ok {
			last = l
		}
	}
	return last
}

// 对于gRPC的错误提取出：gRPC错误码、我们自定义的错误码和错误信息
func ParseGRPCErrMsg(err error) (grpcCode codes.Code, code int, msg string) {
	st := status.Convert(err)
	grpcCode = st.Code()
	code, msg = ParseErrMsg(st.Message())
	return
}

// 用于跨程序传递err时，从errMsg中取出code和msg
// 当然可以用unsafe指针转换来提高性能，但是感觉应该没必要那么写。。
func ParseErrMsg(rawMsg string) (code int, msg string) {
	if !strings.HasPrefix(rawMsg, "code: ") {
		return -1, rawMsg
	}
	errFields := strings.SplitN(rawMsg, ", ", 2)
	if len(errFields) == 1 {
		return -1, rawMsg
	}
	codeStr := strings.TrimPrefix(errFields[0], "code: ")
	code, e := strconv.Atoi(codeStr)
	if e != nil {
		return -1, errFields[1]
	}

	return code, errFields[1]
}
