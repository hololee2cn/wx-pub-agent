package persistence

import (
	"fmt"
	oslog "log"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	captchaPb "github.com/hololee2cn/captcha/pkg/grpcIFace"
	smsPb "github.com/hololee2cn/sms-xuanwu/pkg/grpcIFace"
	redis3 "github.com/hololee2cn/wxpub/v1/src/pkg/redis"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repositories struct {
	DB                *gorm.DB
	Redis             *redis.UniversalClient
	SmsGRPCClient     smsPb.SenderClient
	CaptchaGRPCClient captchaPb.CaptchaServiceClient
}

type DBConfig struct {
	DBDriver, DBUser, DBPassword, DBHost, DBName string
	MaxIdleConn, MaxOpenConn                     int
}

var CommonRepositories Repositories

func NewRepositories(DBConfig DBConfig, redisAddresses []string, smsRPCAddr, captchaRPCAddr string, debugMode bool) (func(), error) {
	dbClean, err := NewDBRepositories(DBConfig, debugMode)
	if err != nil {
		return nil, err
	}
	redisClean, err := NewRedisRepositories(redisAddresses)
	if err != nil {
		return nil, err
	}
	smsClean, err := NewSmsGRPCClientRepositories(smsRPCAddr)
	if err != nil {
		return nil, err
	}
	captchaClean, err := NewCaptchaGRPCClientRepositories(captchaRPCAddr)
	if err != nil {
		return nil, err
	}

	// persistence repo init
	NewAkRepo()
	NewMessageRepo()
	NewUserRepo()
	NewWxRepo()
	NewPhoneVerifyRepo()
	NewTmplRepo()

	// release all the resources
	return func() {
		dbClean()
		redisClean()
		smsClean()
		captchaClean()
	}, nil
}

func NewDBRepositories(config DBConfig, debugMode bool) (func(), error) {
	dbUser, dbPassword, dbHost, dbName := config.DBUser, config.DBPassword, config.DBHost, config.DBName
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&interpolateParams=true", dbUser, dbPassword, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	if config.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour) // 设置最大连接超时
	}
	if debugMode {
		newLogger := logger.New(
			oslog.New(os.Stdout, "\r\n", oslog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,         // Disable color
			},
		)
		db.Logger = newLogger
	}

	CommonRepositories.DB = db
	return func() {
		_ = sqlDB.Close()
	}, nil
}

func NewRedisRepositories(addresses []string) (func(), error) {
	redisClean, redisClient := redis3.NewRedisClient(addresses)
	err := redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}
	CommonRepositories.Redis = &redisClient
	log.Info("redis client init success")
	return redisClean, nil
}

func NewSmsGRPCClientRepositories(smsRPCAddr string) (func(), error) {
	smsConn, err := grpc.Dial(smsRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to dial sms grpc server: %+v", err)
		return nil, err
	}
	smsClient := smsPb.NewSenderClient(smsConn)
	CommonRepositories.SmsGRPCClient = smsClient
	return func() {
		_ = smsConn.Close()
	}, nil
}

func NewCaptchaGRPCClientRepositories(captchaRPCAddr string) (func(), error) {
	captchaConn, err := grpc.Dial(captchaRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to dial captcha grpc server: %+v", err)
		return nil, err
	}

	captchaClient := captchaPb.NewCaptchaServiceClient(captchaConn)
	CommonRepositories.CaptchaGRPCClient = captchaClient
	return func() {
		_ = captchaConn.Close()
	}, nil
}
