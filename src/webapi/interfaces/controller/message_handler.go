package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/errorx"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	message application.MessageInterface
}

func NewMessageController(msg application.MessageInterface) *Message {
	return &Message{
		message: msg,
	}
}

func (a *Message) SendTmplMessage(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var param entity.SendTmplMsgReq
	ginx.BindJSON(c, &param)
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("SendTmplMessage validate sendmsg req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		ginx.BombErr(errorx.CodeInvalidParams, errMsg)
	}
	var msgResp entity.SendTmplMsgResp
	msgResp, err := a.message.SendTmplMsg(ctx, param)
	ginx.NewRender(c).Data(msgResp, err)
}

func (a *Message) TmplMsgStatus(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var param entity.TmplMsgStatusReq
	param.RequestID = ginx.URLParamStr(c, "id")
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("TmplMsgStatus validate req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		ginx.BombErr(errorx.CodeInvalidParams, errMsg)
	}
	var msgStatusResp entity.TmplMsgStatusResp
	msgStatusResp, err := a.message.TmplMsgStatus(ctx, param.RequestID)
	ginx.NewRender(c).Data(msgStatusResp, err)
}
