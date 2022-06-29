package webapi

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hololee2cn/pkg/extra"
	"github.com/hololee2cn/wxpub/v1/src/pkg/httpx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	"github.com/hololee2cn/wxpub/v1/src/webapi/g"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	"github.com/hololee2cn/wxpub/v1/src/webapi/interfaces/webapi/router"
	"github.com/hololee2cn/wxpub/v1/src/webapi/tasks"
	log "github.com/sirupsen/logrus"
)

// Run run webapi
func Run(ctx context.Context, cancelFunc context.CancelFunc) {
	code := 1
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cleanFunc, err := initialize(ctx)
	if err != nil {
		fmt.Println("webapi init fail:", err)
		os.Exit(code)
	}

EXIT:
	for {
		sig := <-quit
		log.Infoln("received signal:", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			code = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	log.Infoln("Shutting down server...")

	cleanFunc()
	cancelFunc()
	fmt.Println("webapi exited")

	os.Exit(code)
}

func initialize(ctx context.Context) (func(), error) {
	// init config
	config.Init()
	// init log
	extra.Default(config.LogLevel)
	// init service
	cleanFunc, err := InitService()
	if err != nil {
		return nil, err
	}
	// task start
	tasks.ConsumerTask(ctx)

	engine := router.New()
	httpClean := httpx.Init(config.ListenAddr, engine)
	go g.Wait()
	return func() {
		cleanFunc()
		httpClean()
	}, nil
}

func InitService() (func(), error) {
	debugMode := config.SMode == consts.ServerModeDebug
	dbConf := persistence.DBConfig{
		DBUser:      config.DBUser,
		DBPassword:  config.DBPassword,
		DBHost:      config.DBHost,
		DBName:      config.DBName,
		MaxIdleConn: config.DBMaxIdleConn,
		MaxOpenConn: config.DBMaxOpenConn,
	}
	cleanFunc, err := persistence.NewRepositories(dbConf, config.RedisAddresses, config.SmsRPCAddr, config.CaptchaRPCAddr, debugMode)
	if err != nil {
		return nil, err
	}
	// repository init
	repository.NewWXRepository(
		persistence.DefaultWxRepo(), persistence.DefaultUserRepo(), persistence.DefaultMessageRepo())
	repository.NewAccessTokenRepository(
		persistence.DefaultAkRepo())
	repository.NewUserRepository(
		persistence.DefaultUserRepo(), persistence.DefaultPhoneVerifyRepo())
	repository.NewMessageRepository(
		persistence.DefaultMessageRepo(), persistence.DefaultUserRepo())
	repository.NewTmplRepository(
		persistence.DefaultTmplRepo())
	return cleanFunc, nil
}
