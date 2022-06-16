package model

import (
	"image/color"
)

type CaptchaResponse struct {
	ID          string `json:"id"`
	Base64Value string `json:"value"` // 验证码以base64形式来存储

	Answer string `json:"answer"` // 调试时可打印出来
}

// 获取验证码时的通用选项(不管什么类型的验证码, 某些选项对某种类型的验证码不适用, 其实没关系, 设置不起作用而已)
type CaptchaCommonOpts struct {
	// 验证码类型
	Type string

	// from audio
	// Language possible values for lang are "en", "ja", "ru", "zh".
	AudioLanguage string

	// 图片大小(像素)
	Width  int
	Height int

	// 输出多少个字符
	Length int

	// from string
	// NoiseCount 噪音值
	NoiseCount int

	// ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	// Source is a unicode which is the rand string from.
	// 这个也不要外部指定, 免得别人弄一些政治敏感词 或者 指定的个数太少容易被识别
	// Source string

	// BgColor captcha image background color (optional)
	BgColor *color.RGBA

	// Fonts loads by name see fonts.go's comment
	// 不传默认用所有内置的字体
	// Fonts      []string
	// fontsArray []*truetype.Font

	// from chinese
	// all redeclared

	// from math
	// all redeclared

	// from digit
	DigitMaxSkew float64
	// DotCount Number of background circles.
	DigitDotCount int

	// 上述的 BgColor 对 digit类型不适用, 这其实无所谓, 传了也不影响, 不用额外说明

	// 下面的选项属于业务上的

	// 验证码有效期, 单位为秒
	MaxAge int64

	// 调试模式
	Debug bool
}
