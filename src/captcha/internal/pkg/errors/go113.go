// +build go1.13

package errors

import stdErrors "errors"

func Unwrap(err error) error {
	return stdErrors.Unwrap(err)
}
func Is(err, target error) bool {
	return stdErrors.Is(err, target)
}
func As(err error, target interface{}) bool {
	return stdErrors.As(err, target)
}
