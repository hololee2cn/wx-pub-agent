package server

import (
	"context"
	"image/color"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/model"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/service"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/utils"

	pb "github.com/hololee2cn/captcha/pkg/grpcIFace"
	me "github.com/hololee2cn/wxpub/v1/src/captcha/internal/errors"
	ce "github.com/hololee2cn/wxpub/v1/src/captcha/internal/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type captchaSvcServer struct {
	pb.UnimplementedCaptchaServiceServer
	captchaSvc service.CaptchaSvc
}

func NewCaptchaSvcServer(captchaSvc service.CaptchaSvc) *captchaSvcServer {
	return &captchaSvcServer{captchaSvc: captchaSvc}
}

func (s *captchaSvcServer) Get(ctx context.Context, req *pb.GetCaptchaRequest) (
	resp *pb.GetCaptchaResponse, err error) {

	var bgColor color.RGBA
	bgColor, err = utils.StrToRGBA(req.GetBgColor())
	if err != nil {
		err = ce.Wrap(err, "bg format", me.CodeInvalidBgColorFormat)
		log.Error(err)
		return
	}

	opts := model.CaptchaCommonOpts{
		Type:   req.GetType(),
		Width:  int(req.GetWidth()),
		Height: int(req.GetHeight()),
		Length: int(req.GetLength()),

		MaxAge: req.GetMaxAge(),

		AudioLanguage:   req.GetAudioLanguage(),
		NoiseCount:      int(req.GetNoiseCount()),
		ShowLineOptions: int(req.ShowLineOptions),
		BgColor:         &bgColor,
		DigitMaxSkew:    req.GetDigitMaxSkew(),
		DigitDotCount:   int(req.GetDigitDotCount()),
		Debug:           req.GetDebug(),
	}

	var captcha model.CaptchaResponse
	captcha, err = s.captchaSvc.GenCaptcha(&opts)
	if err != nil {
		err = ce.Wrap(err, "gen captcha", me.CodeFailedGenCaptcha)
		log.Error(err)
		return
	}
	resp = &pb.GetCaptchaResponse{
		ID:          captcha.ID,
		Base64Value: captcha.Base64Value,
		Answer:      captcha.Answer,
	}
	return
}

func (s *captchaSvcServer) Verify(ctx context.Context, req *pb.VerifyCaptchaRequest) (
	resp *pb.VerifyCaptchaResponse, err error) {

	match, err := s.captchaSvc.VerifyCaptcha(req.GetID(), req.GetAnswer(), true)
	if err != nil {
		err = ce.Wrap(err, "verify", me.CodeFailedVerify)
		log.Error(err)
		return
	}

	resp = &pb.VerifyCaptchaResponse{
		Data: match,
	}
	return
}
