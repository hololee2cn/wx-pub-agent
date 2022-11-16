package utils

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/model"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/config"

	"github.com/mojocn/base64Captcha"
)

var (
	_reRGBAStr = regexp.MustCompile(`^#[0-9A-Fa-f]{8}$`) // #rrggbbaa
)

func StrToRGBA(s string) (rgba color.RGBA, err error) {
	if s == "" {
		return
	}

	if !_reRGBAStr.MatchString(s) {
		err = fmt.Errorf("rgba str must match regexp: %s", _reRGBAStr.String())
		return
	}

	s = strings.TrimPrefix(s, "#")
	colors := [4]int64{}

	for i := range colors {
		colors[i], err = strconv.ParseInt(s[i*2:i*2+2], 16, 32)
		if err != nil {
			return
		}
	}

	rgba = color.RGBA{
		R: uint8(colors[0]),
		G: uint8(colors[1]),
		B: uint8(colors[2]),
		A: uint8(colors[3]),
	}
	return
}

// captcha redis set key
func FullID(id string) string {
	if strings.HasPrefix(id, config.Module) {
		return id
	}
	return config.Module + "-" + id
}

// 通过传入参数构建统一的 captcha driver
func UnifyCaptchaDriver(opts *model.CaptchaCommonOpts) (driver base64Captcha.Driver, err error) {
	if opts == nil {
		err = fmt.Errorf("nil opts passed")
		return
	}
	if opts.Type == "" {
		opts.Type = config.CaptchaTypeString
	}

	// 统一设置默认值
	// 虽然有的选项不是通用的, 如: AudioDefaultLanguage, 统一设置也无妨
	if opts.Width <= 0 {
		opts.Width = config.CaptchaDefaultWidth
	}
	if opts.Height <= 0 {
		opts.Height = config.CaptchaDefaultHeight
	}
	if opts.Length <= 0 {
		opts.Length = config.CaptchaDefaultLength
	}
	if opts.AudioLanguage == "" {
		opts.AudioLanguage = config.AudioDefaultLanguage
	}
	if opts.MaxAge <= 0 {
		opts.MaxAge = config.CaptchaDefaultMaxAge
	}
	if opts.NoiseCount <= 0 {
		opts.NoiseCount = config.CaptchaDefaultNoiseCount
	}
	if opts.ShowLineOptions <= 0 {
		opts.ShowLineOptions = config.OptionsShowAllLines
	}

	switch opts.Type {
	case config.CaptchaTypeAudio:
		driver = &base64Captcha.DriverAudio{
			Length:   opts.Length,
			Language: opts.AudioLanguage,
		}
	case config.CaptchaTypeDigit:
		driver = &base64Captcha.DriverDigit{
			Width:    opts.Width,
			Height:   opts.Height,
			Length:   opts.Length,
			MaxSkew:  opts.DigitMaxSkew,
			DotCount: opts.DigitDotCount,
		}
	case config.CaptchaTypeMath:
		driver = (&base64Captcha.DriverMath{
			Width:           opts.Width,
			Height:          opts.Height,
			NoiseCount:      opts.NoiseCount,
			ShowLineOptions: opts.ShowLineOptions,
			BgColor:         opts.BgColor,
			Fonts:           config.DefaultFonts,
		}).ConvertFonts()
	case config.CaptchaTypeString:
		driver = (&base64Captcha.DriverString{
			Width:           opts.Width,
			Height:          opts.Height,
			Length:          opts.Length,
			NoiseCount:      opts.NoiseCount,
			ShowLineOptions: opts.ShowLineOptions,
			BgColor:         opts.BgColor,
			Source:          config.TxtNumbersAndAlphabet,
			Fonts:           config.DefaultFonts,
		}).ConvertFonts()
	case config.CaptchaTypeChinese:
		driver = (&base64Captcha.DriverString{
			Width:           opts.Width,
			Height:          opts.Height,
			Length:          opts.Length,
			NoiseCount:      opts.NoiseCount,
			ShowLineOptions: opts.ShowLineOptions,
			BgColor:         opts.BgColor,
			Source:          config.CommonlyUsedChinese,
			Fonts:           config.DefaultFonts,
		}).ConvertFonts()
	default:
		err = fmt.Errorf("unsupported type: %s", opts.Type)
		return
	}

	return
}
