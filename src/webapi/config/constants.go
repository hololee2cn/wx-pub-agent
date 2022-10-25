package config

const (
	// ProjectName 项目名称
	ProjectName = "wx-public-agent"

	// HTTPTimeoutHeader http timeout header
	HTTPTimeoutHeader = "x-nova-timeout"

	// DefaultHTTPTimeOut default http timeout
	DefaultHTTPTimeOut = 60

	// Authorization 认证header
	Authorization = "Authorization"

	// InternalAPITokenHeader internal api token header
	InternalAPITokenHeader = "x-auth-token"

	// DefaultPage default http request page
	DefaultPage = 1

	// DefaultPageSize default http request page size
	DefaultPageSize = 20

	// MaxLimitSize max http request limit size
	MaxLimitSize = 100 // 最大只能查询100，且默认为100

	// Token wx 公众号token
	Token = "nova"

	// RedisKeyAccessToken redis key for access token
	RedisKeyAccessToken = ProjectName + "-access_token"

	// RedisKeyTmpl redis key for template
	RedisKeyTmpl = ProjectName + "-tmpl"

	// RedisKeyPrefixMsgID redis key prefix - sent message
	RedisKeyPrefixMsgID = ProjectName + "-msg_id_"

	// RedisKeyPrefixChallenge redis key prefix - challenge
	RedisKeyPrefixChallenge = ProjectName + "-challenge_"

	// DLockPrefix redis lock prefix
	DLockPrefix = "__dlock-"

	// RedisLockAccessToken redis lock key for access token
	RedisLockAccessToken = DLockPrefix + RedisKeyAccessToken

	// RedisLockTask redis lock key for task
	RedisLockTask = DLockPrefix + "task"

	// MaxRetryCount 消息最大失败重试次数，实际调用接口次数为3*3(http client repeated count)=9
	MaxRetryCount = 3

	// MaxWXCallBackTime 微信回调最大时间
	MaxWXCallBackTime = 15 // 微信回调最大时间

	// MaxHandleMsgCount 最大处理消息条数
	MaxHandleMsgCount = 100

	// VerifyCodeSmsChallengeTTL sms challenge expire ttl 验证短信时候期限设置为30分钟，期间可以重发短信
	VerifyCodeSmsChallengeTTL = 1800

	// RedisMsgIDTTL redis msg key expire ttl
	RedisMsgIDTTL = 30

	// RedisAuthTTL redis auth key expire ttl
	RedisAuthTTL = 300

	// RedisTmplTTL redis template expire ttl 模板缓存时间为1天
	RedisTmplTTL = 3600 * 24

	// SmsContentTemplateCN 短信验证提示语
	SmsContentTemplateCN = "验证码：%s。尊敬的用户，您正在绑定手机号，切勿轻易将验证码告知他人！"
)
