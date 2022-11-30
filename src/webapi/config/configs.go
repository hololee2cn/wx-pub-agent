package config

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/wxpub/v1/src/pkg/env"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var wConfig = new(WConfig)

type ServerMode string

const (
	ServerModeDebug   ServerMode = gin.DebugMode
	ServerModeRelease ServerMode = gin.ReleaseMode
)

type WConfig struct {
	MySQL      MySQLCfg      `toml:"mysql"`
	HttpServer HTTPSvrCfg    `toml:"httpServer"`
	CaptchaSvc CaptchaSvcCfg `toml:"captchaSvc"`
	SmsSvc     SmsSvcCfg     `toml:"smsSvc"`
	LogSvc     LogCfg        `toml:"logSvc"`
	WxSvc      WxCfg         `toml:"wxSvc"`
	Redis      RedisCfg      `toml:"redis"`
	VerifySvc  VerifySvcCfg  `toml:"verifySvc"`
}

type MySQLCfg struct {
	DBHost      string `toml:"dbHost"`
	DBUser      string `toml:"dbUser"`
	DBPassword  string `toml:"dbPassword"`
	DBName      string `toml:"dbName"`
	MaxIdleConn int    `toml:"maxIdleConn"`
	MaxOpenConn int    `toml:"maxOpenConn"`
}

type HTTPSvrCfg struct {
	ListenAddr string `toml:"listenAddr"`
	// SMode 服务端运行状态，已知影响 gin 框架日志输出等级
	SMode ServerMode `toml:"sMode"`
}

type CaptchaSvcCfg struct {
	RPCAddr string `toml:"rpcAddr"`
}

type SmsSvcCfg struct {
	RPCAddr string `toml:"rpcAddr"`
}

type LogCfg struct {
	LogLevel logrus.Level `toml:"logLevel"`
}

type WxCfg struct {
	AppID     string `toml:"appId"`
	AppSecret string `toml:"appSecret"`
}

type RedisCfg struct {
	ClusterAddr []string `toml:"clusterAddr"`
}

type VerifySvcCfg struct {
	Addr string `toml:"addr"`
}

var (
	//go:embed dev_configs.toml
	devConfigs []byte

	//go:embed fat_configs.toml
	fatConfigs []byte

	//go:embed uat_configs.toml
	uatConfigs []byte

	//go:embed pro_configs.toml
	proConfigs []byte
)

func init() {
	var r io.Reader

	switch env.Active().Value() {
	case "dev":
		r = bytes.NewReader(devConfigs)
	case "fat":
		r = bytes.NewReader(fatConfigs)
	case "uat":
		r = bytes.NewReader(uatConfigs)
	case "pro":
		r = bytes.NewReader(proConfigs)
	default:
		r = bytes.NewReader(devConfigs)
	}

	// read config file
	viper.SetConfigType("toml")
	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(wConfig); err != nil {
		panic(err)
	}
	viper.SetConfigName(env.Active().Value() + "_configs")
	viper.AddConfigPath("./src/webapi/config")

	configFile := "./src/webapi/config/" + env.Active().Value() + "_configs.toml"
	if _, ok := isExists(configFile); !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}()

		if err = viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(wConfig); err != nil {
			panic(err)
		}
	})
	// 设置日志等级
	wConfig.setLogLevel()
}

// IsExists 文件是否存在
func isExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

func (c *WConfig) setLogLevel() {
	if c.HttpServer.SMode == ServerModeDebug {
		c.LogSvc.LogLevel = logrus.DebugLevel
	} else if c.HttpServer.SMode == ServerModeRelease {
		c.LogSvc.LogLevel = logrus.InfoLevel
	}
}

func Get() WConfig {
	return *wConfig
}
