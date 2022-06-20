package entity

type AccessTokenReq struct {
	//	获取access_token填写client_credential
	GrantType string `json:"grant_type"`
	//	第三方用户唯一凭证
	AppID string `json:"appid"`
	//	第三方用户唯一凭证密钥，即appsecret
	Secret string `json:"secret"`
}

type AccessTokenResp struct {
	// 获取到的凭证
	AccessToken string `json:"access_token"`
	// 凭证有效时间，单位：秒
	ExpiresIn int64 `json:"expires_in"`
	ErrorInfo
}

type GetAccessTokenResp struct {
	// 请求wx凭证ak
	AccessToken string `json:"access_token"`
}
