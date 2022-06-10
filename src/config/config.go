package config

import (
	"strings"
	"sync"

	config2 "github.com/hololee2cn/wxpub/v1/src/pkg/config"

	"github.com/hololee2cn/wxpub/v1/src/consts"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel   logrus.Level
	ListenAddr = config2.DefaultString("listen_addr", ":80")
	// SMode 服务端运行状态，已知影响 gin 框架日志输出等级
	SMode = consts.ServerMode(config2.DefaultString("server_mode", string(consts.ServerModeRelease)))

	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        = config2.DefaultString("db_name", "pub_platform_mgr")
	DBMaxIdleConn = config2.DefaultInt("max_db_idle_conn", 1000)
	DBMaxOpenConn = config2.DefaultInt("max_db_open_conn", 1000)

	CaptchaRPCAddr       = config2.DefaultString("captcha_rpc_addr", "captcha.common:80")
	SmsRPCAddr           = config2.DefaultString("sms_rpc_addr", "sms-xuanwu.misc-pub:80")
	SmsContentTemplateCN = config2.DefaultString("sms_content_template_cn", "南凌科技验证码：%s。尊敬的用户，您正在绑定手机号，切勿轻易将验证码告知他人！")

	WXBaseURL        = config2.DefaultString("wx_base_url", "https://api.weixin.qq.com")
	WXAccessTokenURL = config2.DefaultString("wx_access_token_url", WXBaseURL+"/cgi-bin/token")
	WXMsgTmplSendURL = config2.DefaultString("wx_msg_tmpl_send_url", WXBaseURL+"/cgi-bin/message/template/send")

	RedisAddresses   []string
	AppID            string
	AppSecret        string
	VerifyProfileURL string
	once             sync.Once
)

func Init() {
	once.Do(func() {
		DBHost = config2.MustString("db_host")
		DBUser = config2.MustString("db_user")
		DBPassword = config2.MustString("db_password")

		if SMode == consts.ServerModeDebug {
			LogLevel = logrus.DebugLevel
		} else {
			LogLevel = logrus.InfoLevel
		}

		RedisAddresses = strings.Split(config2.MustString("redis_addresses"), ",")

		AppID = config2.MustString("app_id")
		AppSecret = config2.MustString("app_secret")
		VerifyProfileURL = config2.MustString("verify_profile_url")
	})
}
