package repository

import (
	"context"

	"github.com/hololee2cn/pkg/ginx"

	"github.com/hololee2cn/wxpub/v1/src/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/infrastructure/persistence"
	"github.com/hololee2cn/wxpub/v1/src/utils"

	log "github.com/sirupsen/logrus"
)

type MessageRepository struct {
	msg  *persistence.MessageRepo
	user *persistence.UserRepo
}

var defaultMessageRepository = &MessageRepository{}

func NewMessageRepository(msg *persistence.MessageRepo, user *persistence.UserRepo) {
	if defaultMessageRepository.msg == nil {
		defaultMessageRepository.msg = msg
	}
	if defaultMessageRepository.user == nil {
		defaultMessageRepository.user = user
	}
}

func DefaultMessageRepository() *MessageRepository {
	return defaultMessageRepository
}

func (t *MessageRepository) SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsg traceID:%s", traceID)
	var resp entity.SendTmplMsgResp
	var err error
	// 生成request id
	requestID, err := utils.GetUUID()
	if err != nil {
		log.Errorf("SendTmplMsg MessageRepository GetUUID failed,traceID:%s,err:%+v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	resp.RequestID = requestID
	users, err := t.user.ListUserByPhones(ctx, param.Phones)
	if err != nil {
		log.Errorf("SendTmplMsg ListUserByPhones failed,traceID:%s,err:%+v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	userPhoneMap := make(map[string]entity.User)
	for _, user := range users {
		userPhoneMap[user.Phone] = user
	}
	msgLogs := make([]entity.MsgLog, 0)
	for _, phone := range param.Phones {
		if user, ok := userPhoneMap[phone]; ok {
			msgLogs = append(msgLogs, param.TransferPendingMsgLog(requestID, user.OpenID, user.Phone))
		} else {
			msgLogs = append(msgLogs, param.TransferFailureMsgLog(requestID, "", phone))
		}
	}
	// 消息批量存入db
	err = t.msg.BatchSaveMsgLog(ctx, msgLogs)
	if err != nil {
		log.Errorf("SendTmplMsg BatchSaveMsgLog failed,traceID:%s,err:%+v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	return resp, nil
}

func (t *MessageRepository) TmplMsgStatus(ctx context.Context, requestID string) (entity.TmplMsgStatusResp, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("TmplMsgStatus traceID:%s", traceID)
	var resp entity.TmplMsgStatusResp
	resp.Lists = make([]entity.TmplMsgStatusItem, 0)
	var count int64
	var err error
	count, err = t.msg.ListMsgLogsByReqIDCnt(ctx, requestID)
	if err != nil {
		log.Errorf("TmplMsgStatus ListMsgLogsByReqIDCnt failed,traceID:%s,err:%+v", traceID, err)
		return entity.TmplMsgStatusResp{}, err
	}
	items, err := t.msg.ListMsgLogsByReqID(ctx, requestID)
	if err != nil {
		log.Errorf("TmplMsgStatus ListMsgLogsByRequestID failed,traceID:%s,err:%+v", traceID, err)
		return entity.TmplMsgStatusResp{}, err
	}
	for _, item := range items {
		resp.Lists = append(resp.Lists, item.TransferTmplMsgStatusItem())
	}
	resp.Total = int(count)
	return resp, nil
}
