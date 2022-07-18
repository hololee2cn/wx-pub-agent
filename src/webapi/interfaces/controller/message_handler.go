package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/errorx"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
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
	// 模板查找
	_, err := repository.DefaultTmplRepository().GetTemplate(ctx, param.TmplMsgID)
	if err != nil {
		log.Errorf("SendTmplMessage validate GetTemplate req param failed, traceID:%s, err:%+v", traceID, err)
		ginx.NewRender(c, http.StatusBadRequest).Message(err.Error())
		return
	}
	var msgResp entity.SendTmplMsgResp
	msgResp, err = a.message.SendTmplMsg(ctx, param)
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
	var msgStatusResp *entity.TmplMsgStatusResp
	msgStatusResp, err := a.message.TmplMsgStatus(ctx, param.RequestID)
	if msgStatusResp == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such msg status by request id")
		return
	}
	ginx.NewRender(c).Data(msgStatusResp, err)
}
