package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/wxutil"
	log "github.com/sirupsen/logrus"
)

type WX struct {
	wx application.WXInterface
}

func NewWXController(awApp application.WXInterface) *WX {
	return &WX{
		wx: awApp,
	}
}

func (a *WX) GetWXCheckSign(c *gin.Context) {
	var param entity.WXCheckReq
	ginx.BindQuery(c, &param)
	// wx开放平台验证
	ok := wxutil.CheckSign(param.Signature, wxutil.CalcSign(param.TimeStamp, param.Nonce, consts.Token))
	if !ok {
		log.Infof("wx public platform access failed!")
		return
	}
	// 原样返回
	ginx.NewRender(c).RawString(param.EchoStr)
}

func (a *WX) HandleXML(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var param entity.WXCheckReq
	ginx.BindQuery(c, &param)
	// wx开放平台验证
	ok := wxutil.CheckSign(param.Signature, wxutil.CalcSign(param.TimeStamp, param.Nonce, consts.Token))
	if !ok {
		log.Infof("wx public platform access failed!")
		return
	}
	var reqBody *entity.TextRequestBody
	ginx.BindXML(c, &reqBody)
	// 事件xml返回
	respBody, err := a.wx.HandleXML(ctx, reqBody)
	ginx.NewRender(c).DataString(string(respBody), err)
}
