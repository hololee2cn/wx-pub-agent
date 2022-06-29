package application

import (
	"context"

	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
)

type userApp struct {
	user repository.UserRepository
}

// userApp implements the UserInterface
var _ UserInterface = &userApp{}

type UserInterface interface {
	GetUserByOpenID(ctx context.Context, openID string) (entity.User, error)
	SaveUser(ctx context.Context, user entity.User, isUpdateAll bool) error
	GenCaptcha(ctx context.Context, width int32, height int32) (string, string, error)
	VerifyCaptcha(ctx context.Context, captchaID string, captchaAnswer string) (bool, error)
	VerifyPhone(ctx context.Context, phone string) (bool, error)
	SendSms(ctx context.Context, req entity.SendSmsReq) error
	VerifySmsCode(ctx context.Context, req entity.VerifyCodeReq) (bool, bool, error)
}

func (u *userApp) GetUserByOpenID(ctx context.Context, openID string) (entity.User, error) {
	return u.user.GetUserByOpenID(ctx, openID)
}

func (u *userApp) SaveUser(ctx context.Context, user entity.User, isUpdateAll bool) error {
	return u.user.SaveUser(ctx, user, isUpdateAll)
}

func (u *userApp) GenCaptcha(ctx context.Context, width int32, height int32) (string, string, error) {
	return u.user.GenCaptcha(ctx, width, height)
}

func (u *userApp) VerifyCaptcha(ctx context.Context, captchaID string, captchaAnswer string) (bool, error) {
	return u.user.VerifyCaptcha(ctx, captchaID, captchaAnswer)
}

func (u *userApp) VerifyPhone(ctx context.Context, phone string) (bool, error) {
	return u.user.VerifyPhone(ctx, phone)
}

func (u *userApp) SendSms(ctx context.Context, req entity.SendSmsReq) error {
	return u.user.SendSms(ctx, req)
}

func (u *userApp) VerifySmsCode(ctx context.Context, req entity.VerifyCodeReq) (bool, bool, error) {
	return u.user.VerifySmsCode(ctx, req)
}
