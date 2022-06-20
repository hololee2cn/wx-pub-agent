package entity

type ErrorInfo struct {
	// 错误码
	ErrCode int64 `json:"errcode"`
	// 错误信息
	ErrMsg string `json:"errmsg"`
}
