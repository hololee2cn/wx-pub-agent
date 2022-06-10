package application

import (
	"context"

	"github.com/hololee2cn/wxpub/v1/src/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/domain/repository"
)

type messageApp struct {
	message repository.MessageRepository
}

// messageApp implements the MessageInterface
var _ MessageInterface = &messageApp{}

type MessageInterface interface {
	SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error)
	TmplMsgStatus(ctx context.Context, requestID string) (entity.TmplMsgStatusResp, error)
}

func (u *messageApp) SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	return u.message.SendTmplMsg(ctx, param)
}

func (u *messageApp) TmplMsgStatus(ctx context.Context, requestID string) (entity.TmplMsgStatusResp, error) {
	return u.message.TmplMsgStatus(ctx, requestID)
}
