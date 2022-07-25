package tasks

import (
	"context"

	"github.com/hololee2cn/pkg/errorx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	"github.com/hololee2cn/wxpub/v1/src/webapi/g"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	log "github.com/sirupsen/logrus"
)

var (
	msgRepo      *persistence.MessageRepo
	akRepository *repository.AccessTokenRepository
	maxMsgChan   chan struct{}
)

func init() {
	maxMsgChan = make(chan struct{}, 100)
}

// HandleMsg 处理消息，状态为发送中，针对消息失败记录失败次数，失败3次，状态改为发送失败，消息写入db无序
func HandleMsg(ctx context.Context) {
	msgRepo = persistence.DefaultMessageRepo()
	akRepository = repository.DefaultAccessTokenRepository()
	select {
	case <-ctx.Done():
		log.Infof("HandleMsg exit..., err:%+v", ctx.Err())
		return
	default:
		// 扫全表状态为pending的消息
		msgLogs, e := msgRepo.ListPendingMsgLogs(ctx)
		if e != nil {
			log.Errorf("HandleMsg GetListPendingMsgLogs failed,err:%+v", e)
			break
		}
		// 发送消息并修改发送状态，允许消息重复消费
		for _, msgLog := range msgLogs {
			maxMsgChan <- struct{}{}
			g.Add(1)
			go func(msgLog entity.MsgLog) {
				defer func() {
					g.Done()
					<-maxMsgChan
				}()
				handlePerMsg(ctx, msgLog)
			}(msgLog)
		}
	}
}

// handlePerMsg 处理每个待发送消息
func handlePerMsg(ctx context.Context, msgLog entity.MsgLog) {
	// 获取access token
	var ak string
	var err error
	ak, err = akRepository.GetAccessToken(ctx)
	if err != nil {
		log.Errorf("handlePerMsg GetAccessToken failed,err:%+v", err)
		return
	}
	var msgSendReq entity.SendTmplMsgRemoteReq
	msgSendReq, err = msgLog.TransferSendTmplMsgRemoteReq()
	if err != nil {
		log.Errorf("handlePerMsg TransferSendTmplMsgRemoteReq failed,err:%v", err)
		return
	}
	msgSendReq.AccessToken = ak
	log.Infof("msg send req ak is %s", ak)
	var resp entity.SendTmplMsgRemoteResp
	resp, err = msgRepo.SendTmplMsgFromRequest(ctx, msgSendReq)
	if err != nil {
		log.Errorf("handlePerMsg SendTmplMsgFromRequest failed,msgLog:%+v,err:%+v", msgLog, err)
		// token过期重试或者失效重新刷新ak去请求重试，不记录重试次数
		if resp.ErrCode == errorx.CodeRIDExpired || resp.ErrCode == errorx.CodeRIDUnauthorized || resp.ErrCode == errorx.CodeRIDOpenIDInvalid {
			_, err = akRepository.FreshAccessToken(ctx)
			if err != nil {
				log.Errorf("handlePerMsg FreshAccessToken failed,err:%+v", err)
				return
			}
		}
		// 记录当前重试次数+1
		msgLog.Count++
		err = msgRepo.UpdateMsgLog(ctx, msgLog)
		if err != nil {
			log.Errorf("handlePerMsg UpdateMsgLog failed,msgLog:%+v,err:%+v", msgLog, err)
			return
		}
		return
	}
	// 记录当前消息发送状态为发送中等待回调状态确认
	msgLog.MsgID = resp.MsgID
	msgLog.Status = consts.Sending
	msgLog.Count++
	err = msgRepo.UpdateMsgLog(ctx, msgLog)
	if err != nil {
		log.Errorf("handlePerMsg UpdateMsgLog failed,msgLog:%+v,err:%+v", msgLog, err)
		return
	}
}

// HandleFailSendMsg 扫描超过三次重试的消息改为失败状态
func HandleFailSendMsg(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Infof("HandleFailSendMsg exit..., err:%+v", ctx.Err())
		return
	default:
		// 待发送错误重试消息处理
		err := persistence.DefaultMessageRepo().UpdateMaxRetryCntMsgLogsStatus(ctx)
		if err != nil {
			log.Errorf("HandleFailSendMsg UpdateMaxRetryCntMsgLogsStatus failed,err:%+v", err)
		}
		// 发送中回调最大时间消息处理
		err = persistence.DefaultMessageRepo().UpdateTimeoutMsgLogsStatus(ctx)
		if err != nil {
			log.Errorf("HandleFailSendMsg UpdateTimeoutMsgLogsStatus failed,err:%+v", err)
		}
	}
}
