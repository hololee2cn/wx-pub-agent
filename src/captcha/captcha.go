package captcha

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpcMW "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hololee2cn/pkg/extra"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/rpc/server"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/service"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/store"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/config"
	me "github.com/hololee2cn/wxpub/v1/src/captcha/internal/errors"
	"github.com/hololee2cn/wxpub/v1/src/pkg/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	log "github.com/sirupsen/logrus"
)

var (
	captchaSvc = service.NewDefaultCaptchaSvc(store.NewRedisStore(config.CaptchaDefaultMaxAge))
)

// Run run captcha
func Run() {
	code := 1
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cleanFunc, err := initialize()
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
	fmt.Println("captcha exited")

	os.Exit(code)
}

func initialize() (func(), error) {
	config.Init()
	extra.Default(log.Level(config.LogLevel))
	log.Info("captcha starting...")

	redisClean, _ := redis.NewRedisClient(config.RedisAddrs)

	captchaSvcServer := server.NewCaptchaSvcServer(captchaSvc)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcMW.ChainUnaryServer(
			getFullMethodName,
			logRequest,
		)),
	}
	rs := server.NewRpcServer(config.RPCAddr, captchaSvcServer, opts...)
	rs.Start()
	return func() {
		redisClean()
		rs.Stop()
	}, nil
}

func getFullMethodName(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fullMethodName := info.FullMethod
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		log.Errorf("no metadata in context")
	} else {
		md.Append("full_method_name", fullMethodName)
	}
	return handler(ctx, req)
}

func logRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		log.Errorf("no metadata in context. info: %+v", *info)
	} else {
		log.Infof("metadata: %+v", md)
		if slice := md.Get(config.KeyTraceID); len(slice) == 0 {
			err := me.NoTraceID(fmt.Sprintf("in metadata: %+v", md))
			log.Error(err)
			return nil, err
		}
	}
	return handler(ctx, req)
}
