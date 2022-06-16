package service

import (
	"testing"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/model"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/store"

	"github.com/stretchr/testify/assert"
)

func Test_Captcha(t *testing.T) {
	captchaSvc := NewDefaultCaptchaSvc(store.NewMemoryStore(3))

	opts := model.CaptchaCommonOpts{
		Debug: true,
	}
	c, err := captchaSvc.GenCaptcha(&opts)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.Answer)

	match, err := captchaSvc.VerifyCaptcha(c.ID, c.Answer, false)
	assert.Nil(t, err)
	assert.True(t, match)

	// 测试一个验证码填写错误的
	match, err = captchaSvc.VerifyCaptcha(c.ID, c.Answer+"a", true)
	assert.Nil(t, err)
	assert.False(t, match)
}
