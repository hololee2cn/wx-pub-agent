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
	// init log
	extra.Default(config.Get().LogSvc.LogLevel)
	// init service
	cleanFunc, err := InitService()
	if err != nil {
		return nil, err
	}
	// task start
	tasks.CronTasks(ctx)

	engine := router.New()
	httpClean := httpx.Init(config.Get().HttpServer.ListenAddr, engine)
	go g.Wait()
	return func() {
		cleanFunc()
		httpClean()
	}, nil
}

func InitService() (func(), error) {
	dbConf := persistence.DBCfg{
		DBUser:      config.Get().MySQL.DBUser,
		DBPassword:  config.Get().MySQL.DBPassword,
		DBHost:      config.Get().MySQL.DBHost,
		DBName:      config.Get().MySQL.DBName,
		MaxIdleConn: config.Get().MySQL.MaxIdleConn,
		MaxOpenConn: config.Get().MySQL.MaxOpenConn,
		DebugMode:   config.Get().HttpServer.SMode == config.ServerModeDebug,
	}
	var err error
	cleanFunc, err := persistence.NewRepositories(
		persistence.NewDBRepositories(dbConf),
		persistence.NewRedisRepositories(persistence.RedisCfg{RedisAddr: config.Get().Redis.ClusterAddr}),
		persistence.NewCaptchaGRPCClientRepositories(persistence.CaptchaCfg{CaptchaRPCAddr: config.Get().CaptchaSvc.RPCAddr}),
		persistence.NewSmsGRPCClientRepositories(persistence.SmsCfg{SmsRPCAddr: config.Get().SmsSvc.RPCAddr}),
	)
	if err != nil {
		return nil, err
	}
	// persistence repo init
	persistence.NewSingletonRepo()

	// repository init
	repository.NewSingletonRepository()

	return cleanFunc, nil
}
