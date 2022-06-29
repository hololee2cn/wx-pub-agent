package repository

import (
	"context"
	"fmt"

	"github.com/hololee2cn/wxpub/v1/src/utils"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
)

type UserRepository struct {
	user        *persistence.UserRepo
	phoneVerify *persistence.PhoneVerifyRepo
}

var defaultUserRepository = &UserRepository{}

func NewUserRepository(user *persistence.UserRepo, phoneVerify *persistence.PhoneVerifyRepo) {
	if defaultUserRepository.user == nil {
		defaultUserRepository.user = user
	}
	if defaultUserRepository.phoneVerify == nil {
		defaultUserRepository.phoneVerify = phoneVerify
	}
}

func DefaultUserRepository() *UserRepository {
	return defaultUserRepository
}

func (a *UserRepository) GetUserByOpenID(ctx context.Context, openID string) (entity.User, error) {
	return a.user.GetUserByOpenID(ctx, openID)
}

func (a *UserRepository) SaveUser(ctx context.Context, user entity.User, isUpdateAll bool) error {
	return a.user.SaveUser(ctx, user, isUpdateAll)
}

func (a *UserRepository) GenCaptcha(ctx context.Context, width int32, height int32) (string, string, error) {
	return a.phoneVerify.GenCaptcha(ctx, width, height)
}

func (a *UserRepository) VerifyCaptcha(ctx context.Context, captchaID string, captchaAnswer string) (bool, error) {
	return a.phoneVerify.VerifyCaptcha(ctx, captchaID, captchaAnswer)
}

func (a *UserRepository) VerifyPhone(ctx context.Context, phone string) (bool, error) {
	return a.user.IsExistPhone(ctx, phone)
}

func (a *UserRepository) SendSms(ctx context.Context, req entity.SendSmsReq) error {
	verifyCodeAnswer := utils.GenVerifySmsCode()
	err := a.phoneVerify.SetVerifyCodeSmsStorage(ctx, req.OpenID+req.Phone, verifyCodeAnswer)
	if err != nil {
		return err
	}

	content := fmt.Sprintf(config.SmsContentTemplateCN, verifyCodeAnswer)
	sender := consts.SmsSender
	return a.phoneVerify.SendSms(ctx, content, sender, req.Phone)
}

func (a *UserRepository) VerifySmsCode(ctx context.Context, req entity.VerifyCodeReq) (bool, bool, error) {
	return a.phoneVerify.VerifySmsCode(ctx, req.OpenID+req.Phone, req.VerifyCode, consts.RedisAuthTTL)
}
