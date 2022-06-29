package config

import (
	"strings"
	"sync"

	"github.com/hololee2cn/wxpub/v1/src/pkg/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel   logrus.Level
	ListenAddr = config.DefaultString("listen_addr", ":80")
	// SMode 服务端运行状态，已知影响 gin 框架日志输出等级
	SMode = consts.ServerMode(config.DefaultString("server_mode", string(consts.ServerModeRelease)))

	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        = config.DefaultString("db_name", "pub_platform_mgr")
	DBMaxIdleConn = config.DefaultInt("max_db_idle_conn", 1000)
	DBMaxOpenConn = config.DefaultInt("max_db_open_conn", 1000)

	CaptchaRPCAddr       = config.DefaultString("captcha_rpc_addr", "captcha.common:80")
	SmsRPCAddr           = config.DefaultString("sms_rpc_addr", "sms-xuanwu.misc-pub:80")
	SmsContentTemplateCN = config.DefaultString("sms_content_template_cn", "南凌科技验证码：%s。尊敬的用户，您正在绑定手机号，切勿轻易将验证码告知他人！")

	WXBaseURL        = config.DefaultString("wx_base_url", "https://api.weixin.qq.com")
	WXAccessTokenURL = config.DefaultString("wx_access_token_url", WXBaseURL+"/cgi-bin/token")
	WXMsgTmplSendURL = config.DefaultString("wx_msg_tmpl_send_url", WXBaseURL+"/cgi-bin/message/template/send")
	WXListTmplURL    = config.DefaultString("wx_list_tmpl_url", WXBaseURL+"/cgi-bin/template/get_all_private_template")

	RedisAddresses   []string
	AppID            string
	AppSecret        string
	VerifyProfileURL string
	once             sync.Once
)

func Init() {
	once.Do(func() {
		DBHost = config.MustString("db_host")
		DBUser = config.MustString("db_user")
		DBPassword = config.MustString("db_password")

		if SMode == consts.ServerModeDebug {
			LogLevel = logrus.DebugLevel
		} else {
			LogLevel = logrus.InfoLevel
		}

		RedisAddresses = strings.Split(config.MustString("redis_addresses"), ",")

		AppID = config.MustString("app_id")
		AppSecret = config.MustString("app_secret")
		VerifyProfileURL = config.MustString("verify_profile_url")
	})
}
