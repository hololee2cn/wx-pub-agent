package entity

type SendSmsReq struct {
	// 微信用户的openID
	OpenID string `json:"open_id"`
	// 目标手机号
	Phone string `json:"phone"`
	// 图形验证码id
	CaptchaID string `json:"captcha_id"`
	// 图形验证码
	CaptchaAnswer string `json:"captcha_answer"`
}

type VerifyCodeRedisValue struct {
	VerifyCodeAnswer     string `json:"verify_code_answer"`
	VerifyCodeCreateTime int64  `json:"verify_code_create_time"`
}

type VerifyCodeReq struct {
	// 微信用户的openID
	OpenID string `json:"open_id"`
	// 用户名字
	Name string `json:"name"`
	// 目标手机号
	Phone string `json:"phone"`
	// 验证码
	VerifyCode string `json:"verify_code"`
}

type CaptchaResp struct {
	// 验证码，返回验证码id
	CaptchaID string `json:"captcha_id"`
	// 验证码，返回验证码Base64值
	CaptchaBase64Value string `json:"captcha_base_64_value"`
}

type VerifyPhoneReq struct {
	// 目标手机号
	Phone string `json:"phone"`
}
