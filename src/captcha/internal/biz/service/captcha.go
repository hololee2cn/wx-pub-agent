package service

import (
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/model"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/store"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/utils"
	me "github.com/hololee2cn/wxpub/v1/src/captcha/internal/errors"
	ce "github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors"
	"github.com/mojocn/base64Captcha"
	log "github.com/sirupsen/logrus"
)

type CaptchaSvc interface {
	GenCaptcha(opts *model.CaptchaCommonOpts) (c model.CaptchaResponse, err error)
	VerifyCaptcha(id string, answer string, clear bool) (match bool, err error)
}

type defaultCaptchaSvc struct {
	store store.Store
}

func NewDefaultCaptchaSvc(store store.Store) *defaultCaptchaSvc {
	return &defaultCaptchaSvc{store: store}
}

// 生成验证码
func (s *defaultCaptchaSvc) GenCaptcha(opts *model.CaptchaCommonOpts) (c model.CaptchaResponse, err error) {
	var driver base64Captcha.Driver
	if driver, err = utils.UnifyCaptchaDriver(opts); err != nil {
		log.Error(err)
		return
	}

	// c.ID = uuid.Get()
	id, content, answer := driver.GenerateIdQuestionAnswer()
	c.ID = id
	item, err := driver.DrawCaptcha(content)
	if err != nil {
		err = ce.Wrap(err, "generate captcha", me.CodeFailedGenCaptcha)
		log.Error(err)
		return
	}
	if err = s.store.Set(c.ID, answer); err != nil {
		log.Error(err)
		return
	}
	c.Base64Value = item.EncodeB64string()

	if opts.Debug {
		c.Answer = answer
	}
	return
}

// 校验验证码填写对不对
func (s *defaultCaptchaSvc) VerifyCaptcha(id string, answer string, clear bool) (match bool, err error) {
	// 校验时maxAge传递为0就行
	return s.store.Verify(id, answer, clear)
}
