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

// Options application options at all
type Options struct {
	DBConfig      DBCfg
	RedisConfig   RedisCfg
	SmsConfig     SmsCfg
	CaptchaConfig CaptchaCfg
}

type DBCfg struct {
	DBDriver    string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBName      string
	MaxIdleConn int
	MaxOpenConn int
	DebugMode   bool
}

type RedisCfg struct {
	RedisAddr []string
}

type CaptchaCfg struct {
	CaptchaRPCAddr string
}

type SmsCfg struct {
	SmsRPCAddr string
}

type (
	OptionFunc func(*Options)
	RepoFunc   func(*Repositories) error
	CleanFunc  func()
)

var CommonRepositories *Repositories

func NewRepositories(options ...OptionFunc) (func(), error) {
	// load options
	opts := loadOptions(options...)
	// apply config
	applyDB, cleanDB := opts.DBConfig.applyDBConfig()
	applyRedis, cleanRedis := opts.RedisConfig.applyRedisConfig()
	applySms, cleanSms := opts.SmsConfig.applySmsConfig()
	applyCaptcha, cleanCaptcha := opts.CaptchaConfig.applyCaptchaConfig()
	// clean func
	cleanFunc := func() {
		cleanDB()
		cleanRedis()
		cleanSms()
		cleanCaptcha()
	}
	var err error
	CommonRepositories, err = NewRepo(applyDB, applyRedis, applySms, applyCaptcha)
	if err != nil {
		return nil, err
	}
	return cleanFunc, err
}

func NewSingletonRepo() {
	// persistence repo init
	NewAkRepo()
	NewMessageRepo()
	NewUserRepo()
	NewWxRepo()
	NewPhoneVerifyRepo()
	NewTmplRepo()
}

func loadOptions(opts ...OptionFunc) *Options {
	options := &Options{}
	for _, option := range opts {
		option(options)
	}
	return options
}

func NewDBRepositories(dbConfig DBCfg) OptionFunc {
	return func(opt *Options) {
		opt.DBConfig = dbConfig
	}
}

func NewRedisRepositories(redisConfig RedisCfg) OptionFunc {
	return func(opt *Options) {
		opt.RedisConfig = redisConfig
	}
}

func NewSmsGRPCClientRepositories(smsConfig SmsCfg) OptionFunc {
	return func(opt *Options) {
		opt.SmsConfig = smsConfig
	}
}

func NewCaptchaGRPCClientRepositories(captchaConfig CaptchaCfg) OptionFunc {
	return func(opt *Options) {
		opt.CaptchaConfig = captchaConfig
	}
}

func NewRepo(applyFuncs ...RepoFunc) (*Repositories, error) {
	r := &Repositories{}
	for _, applyFunc := range applyFuncs {
		if err := applyFunc(r); err != nil {
			return r, err
		}
	}
	return r, nil
}

func (d *DBCfg) applyDBConfig() (RepoFunc, CleanFunc) {
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&interpolateParams=true", d.DBUser, d.DBPassword, d.DBHost, d.DBName),
	))
	if err != nil {
		return func(*Repositories) error {
			return err
		}, nil
	}
	sqlDB, _ := db.DB()
	if d.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(d.MaxIdleConn)
	}
	if d.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(d.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour) // 设置最大连接超时
	}
	if d.DebugMode {
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
	return func(r *Repositories) error {
		r.DB = db
		return nil
	}, func() { _ = sqlDB.Close() }
}

func (r *RedisCfg) applyRedisConfig() (RepoFunc, CleanFunc) {
	redisClean, redisClient := redis3.NewRedisClient(r.RedisAddr)
	err := redisClient.Ping().Err()
	if err != nil {
		return func(*Repositories) error {
			return err
		}, func() { redisClean() }
	}
	return func(r *Repositories) error {
		r.Redis = &redisClient
		return nil
	}, func() { redisClean() }
}

func (c *CaptchaCfg) applyCaptchaConfig() (RepoFunc, CleanFunc) {
	captchaConn, err := grpc.Dial(c.CaptchaRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return func(*Repositories) error {
			return err
		}, nil
	}
	return func(r *Repositories) error {
		r.CaptchaGRPCClient = captchaPb.NewCaptchaServiceClient(captchaConn)
		return nil
	}, func() { _ = captchaConn.Close() }
}

func (s *SmsCfg) applySmsConfig() (RepoFunc, CleanFunc) {
	smsConn, err := grpc.Dial(s.SmsRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return func(*Repositories) error {
			return err
		}, nil
	}
	return func(r *Repositories) error {
		r.SmsGRPCClient = smsPb.NewSenderClient(smsConn)
		return nil
	}, func() { _ = smsConn.Close() }
}
