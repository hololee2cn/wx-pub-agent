package tasks

import (
	"context"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	"github.com/hololee2cn/wxpub/v1/src/webapi/g"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	"time"

	"github.com/hololee2cn/pkg/errorx"

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

func ConsumerTask(ctx context.Context) {
	msgRepo = persistence.DefaultMessageRepo()
	akRepository = repository.DefaultAccessTokenRepository()
	// 最大发送次数消息处理
	g.Add(1)
	go handleFailSendMsg(ctx)
	// 发送消息处理
	g.Add(1)
	go handleMsg(ctx)
}

// 处理消息，状态为发送中，针对消息失败记录失败次数，失败3次，状态改为发送失败，消息写入db无序
func handleMsg(ctx context.Context) {
	defer func() {
		close(maxMsgChan)
		g.Done()
	}()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Infof("handleMsg exit..., err:%+v", ctx.Err())
			return
		case <-ticker.C:
			// 扫全表状态为pending的消息
			msgLogs, err := msgRepo.ListPendingMsgLogs(ctx)
			if err != nil {
				log.Errorf("handleMsg GetListPendingMsgLogs failed,err:%+v", err)
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
					// 获取access token
					var ak string
					ak, err = akRepository.GetAccessToken(ctx)
					if err != nil {
						log.Errorf("handleMsg GetAccessToken failed,err:%+v", err)
						return
					}
					msgSendReq := msgLog.TransferSendTmplMsgRemoteReq()
					msgSendReq.AccessToken = ak
					log.Infof("msg send req ak is %s", ak)
					var resp entity.SendTmplMsgRemoteResp
					resp, err = msgRepo.SendTmplMsgFromRequest(ctx, msgSendReq)
					if err != nil {
						log.Errorf("handleMsg SendTmplMsgFromRequest failed,msgLog:%+v,err:%+v", msgLog, err)
						// token过期重试或者失效重新刷新ak去请求重试，不记录重试次数
						if resp.ErrCode == errorx.CodeRIDExpired || resp.ErrCode == errorx.CodeRIDUnauthorized {
							_, err = akRepository.FreshAccessToken(ctx)
							if err != nil {
								log.Errorf("handleMsg FreshAccessToken failed,err:%+v", err)
								return
							}
						}
						// 记录当前重试次数+1
						msgLog.Count++
						err = msgRepo.UpdateMsgLog(ctx, msgLog)
						if err != nil {
							log.Errorf("handleMsg UpdateMsgLog failed,msgLog:%+v,err:%+v", msgLog, err)
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
						log.Errorf("handleMsg UpdateMsgLog failed,msgLog:%+v,err:%+v", msgLog, err)
						return
					}
				}(msgLog)
			}
		}
	}
}

// 扫描超过三次重试的消息改为失败状态
func handleFailSendMsg(ctx context.Context) {
	defer g.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Infof("handleFailSendMsg exit..., err:%+v", ctx.Err())
			return
		case <-ticker.C:
			// 待发送错误重试消息处理
			err := msgRepo.UpdateMaxRetryCntMsgLogs(ctx)
			if err != nil {
				log.Errorf("handleFailSendMsg UpdateMaxRetryCntMsgLogs failed,err:%+v", err)
			}
			// 发送中回调最大时间消息处理
			err = msgRepo.UpdateTimeoutMsgLogs(ctx)
			if err != nil {
				log.Errorf("handleFailSendMsg UpdateTimeoutMsgLogs failed,err:%+v", err)
			}
		}
	}
}
