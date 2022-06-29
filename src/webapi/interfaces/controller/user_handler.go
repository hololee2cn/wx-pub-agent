package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/errorx"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/utils"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	log "github.com/sirupsen/logrus"
)

type User struct {
	user application.UserInterface
}

func NewUserController(user application.UserInterface) *User {
	return &User{
		user: user,
	}
}

func (u *User) GenCaptcha(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	width := ginx.QueryStr(c, "width", strconv.Itoa(consts.CaptchaDefaultWidth))
	w, err := strconv.ParseInt(width, 10, 64)
	if err != nil {
		log.Errorf("%s get width failed,err: %+v", traceID, err)
		ginx.BombErr(errorx.CodeInternalServerError, errorx.GetErrorMessage(errorx.CodeInternalServerError))
	}
	height := ginx.QueryStr(c, "height", strconv.Itoa(consts.CaptchaDefaultHeight))
	h, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		log.Errorf("%s get height failed,err: %+v", traceID, err)
		ginx.CustomErr(err)
	}

	captchaID, captchaBase64Value, err := u.user.GenCaptcha(ctx, int32(w), int32(h))
	if err != nil {
		log.Errorf("%s get captchaID and captchaBase64Value failed,err: %+v", traceID, err)
		ginx.CustomErr(err)
	}

	CaptchaResp := entity.CaptchaResp{
		CaptchaID:          captchaID,
		CaptchaBase64Value: captchaBase64Value,
	}
	ginx.NewRender(c).Data(CaptchaResp, nil)
}

func (u *User) SendSms(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var req entity.SendSmsReq
	ginx.BindJSON(c, &req)

	if utils.VerifyMobilePhoneFormat(req.Phone) {
		log.Errorf("invalid phone number: %s, traceID:%s", req.Phone, traceID)
		ginx.BombErr(errorx.CodeInvalidParams, "invalid phone number")
	}

	ok, err := u.user.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaAnswer)
	if err != nil {
		log.Errorf("VerifyCaptcha failed, traceID:%s, err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	if !ok {
		log.Errorf("wrong captcha answer: %s, traceID:%s", req.CaptchaAnswer, traceID)
		ginx.BombErr(errorx.CodeForbidden, "wrong captcha answer")
	}
	exist, err := u.user.VerifyPhone(ctx, req.Phone)
	if err != nil {
		log.Errorf("VerifyPhone failed, traceID:%s, err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	if exist {
		log.Errorf("exist phone:%s", req.Phone)
		ginx.BombErr(errorx.CodeResourcesHasExist, "phone is exist")
	}
	err = u.user.SendSms(ctx, req)
	if err != nil {
		log.Errorf("SendSms failed, traceID:%s, err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	ginx.NewRender(c).Message(nil)
}

func (u *User) VerifyAndUpdatePhone(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var req entity.VerifyCodeReq
	ginx.BindJSON(c, &req)

	if utils.VerifyMobilePhoneFormat(req.Phone) {
		log.Errorf("invalid phone number: %s, traceID:%s", req.Phone, traceID)
		ginx.BombErr(errorx.CodeInvalidParams, "invalid phone number")
	}

	ok, isExpire, err := u.user.VerifySmsCode(ctx, req)
	if err != nil {
		log.Errorf("user VerifySmsCode failed, traceID:%s, err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	if !ok {
		log.Errorf("verify code is not correct, code: %s, traceID: %s", req.VerifyCode, traceID)
		ginx.BombErr(errorx.CodeForbidden, "verify code is not correct")
	}
	if isExpire {
		log.Errorf("sms code is expired, code: %s, traceID: %s", req.VerifyCode, traceID)
		ginx.BombErr(errorx.CodeTokenExpire, "sms code is expired")
	}

	user, err := u.user.GetUserByOpenID(ctx, req.OpenID)
	if err != nil {
		log.Errorf("VerifyAndUpdatePhone get user by open_id error: %+v, traceID: %s", err, traceID)
		ginx.CustomErr(err)
	}

	user.Phone = req.Phone
	user.Name = req.Name
	err = u.user.SaveUser(ctx, user, true)
	if err != nil {
		log.Errorf("VerifyAndUpdatePhone update user error: %+v, traceID: %s", err, traceID)
		ginx.CustomErr(err)
	}
	ginx.NewRender(c).Message(nil)
}
