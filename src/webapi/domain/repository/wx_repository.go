package repository

import (
	"context"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"

	"github.com/hololee2cn/pkg/ginx"

	"github.com/hololee2cn/wxpub/v1/src/utils"

	log "github.com/sirupsen/logrus"
)

type WXRepository struct {
	wx   *persistence.WxRepo
	user *persistence.UserRepo
	msg  *persistence.MessageRepo
}

var defaultWXRepository = &WXRepository{}

func NewWXRepository(wx *persistence.WxRepo, user *persistence.UserRepo, msg *persistence.MessageRepo) {
	if defaultWXRepository.wx == nil {
		defaultWXRepository.wx = wx
	}
	if defaultWXRepository.user == nil {
		defaultWXRepository.user = user
	}
	if defaultWXRepository.msg == nil {
		defaultWXRepository.msg = msg
	}
}

func DefaultWXRepository() *WXRepository {
	return defaultWXRepository
}

func (a *WXRepository) GetWXSign(timestamp, nonce, token string) (sign string) {
	// 本地计算signature
	si := []string{token, timestamp, nonce}
	// 字典序排序
	sort.Strings(si)
	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for _, v := range si {
		b.WriteString(v)
	}
	return utils.Sha1(b.String())
}

func (a *WXRepository) HandleXML(ctx context.Context, reqBody *entity.TextRequestBody) ([]byte, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("HandleXML traceID:%s", traceID)
	if reqBody == nil {
		return nil, fmt.Errorf("xml request body is empty")
	}
	var responseTextBody []byte
	var err error
	switch reqBody.MsgType {
	case consts.Text:
		responseTextBody, err = a.handlerTextXML(ctx, reqBody)
		if err != nil {
			log.Errorf("HandleXML handlerTextXML failed traceID:%s,err:%+v", traceID, err)
			return nil, err
		}
	case consts.Event:
		responseTextBody, err = a.handlerEventXML(ctx, reqBody)
		if err != nil {
			log.Errorf("HandleXML handlerEventXML failed traceID:%s,err:%+v", traceID, err)
			return nil, err
		}
	}
	return responseTextBody, nil
}

func (a *WXRepository) handlerTextXML(ctx context.Context, reqBody *entity.TextRequestBody) ([]byte, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("handlerTextXML traceID:%s", traceID)
	var err error
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	msgID := a.genMsgID(reqBody.FromUserName, reqBody.CreateTime)
	exist, err := a.wx.IsExistMsgIDFromRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerTextXML IsExistMsgIDFromRedis failed,traceID:%s,err:%+v", traceID, err)
		return nil, err
	}
	// 若存在返回空串,不存在则持久化存储,并保存msg id 到 redis
	if exist {
		return nil, nil
	}
	// 关键词回复，keys:告警、alert、绑定、用户
	if !utils.CheckSubstrings(reqBody.Content, []string{"alert", "告警", "绑定", "用户"}...) {
		return nil, nil
	}
	err = a.wx.SetMsgIDToRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerTextXML WXRepository wx repo set msg id to redis failed,traceID:%s,err:%+v", traceID, err)
	}
	// 用户绑定链接推送
	return a.makeTextResponseBody(reqBody.ToUserName, reqBody.FromUserName, fmt.Sprintf("%s%s", consts.SubscribeRespContent, config.VerifyProfileURL))
}

func (a *WXRepository) handlerEventXML(ctx context.Context, reqBody *entity.TextRequestBody) ([]byte, error) {
	var respContent string
	var err error
	// 事件类型
	switch reqBody.Event {
	case consts.SubscribeEvent: // 关注订阅
		if respContent, err = a.handlerSubscribeEvent(ctx, reqBody); err != nil {
			return nil, err
		}
	case consts.UnsubscribeEvent: // 取消关注订阅
		if respContent, err = a.handlerUnSubscribeEvent(ctx, reqBody); err != nil {
			return nil, err
		}
	case consts.TEMPLATESENDJOBFINISHEvent: // 事件回调内部系统错误重发
		if respContent, err = a.handlerTEMPLATESENDJOBFINISHEvent(ctx, reqBody); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	return a.makeTextResponseBody(reqBody.ToUserName, reqBody.FromUserName, respContent)
}

func (a *WXRepository) handlerSubscribeEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("handlerSubscribeEvent traceID:%s", traceID)
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	msgID := a.genMsgID(reqBody.FromUserName, reqBody.CreateTime)
	exist, err := a.isExistUserMsgID(ctx, msgID, reqBody.FromUserName, reqBody.CreateTime)
	if err != nil {
		log.Errorf("handlerSubscribeEvent WXRepository wx repo isExistUserMsgID traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	if exist {
		return "", nil
	}
	// 持久化保存
	u := entity.User{
		OpenID:     reqBody.FromUserName,
		CreateTime: reqBody.CreateTime,
		DeleteTime: 0,
	}
	err = a.user.SaveUser(ctx, u, false)
	if err != nil {
		log.Errorf("handlerSubscribeEvent WXRepository wx repo SaveUser traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	err = a.wx.SetMsgIDToRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerSubscribeEvent WXRepository wx repo set msg id to redis failed,traceID:%s,err:%+v", traceID, err)
	}
	return fmt.Sprintf("%s%s", consts.SubscribeRespContent, config.VerifyProfileURL), nil
}

func (a *WXRepository) handlerUnSubscribeEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("handlerUnSubscribeEvent traceID:%s", traceID)
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	msgID := a.genMsgID(reqBody.FromUserName, reqBody.CreateTime)
	exist, err := a.isExistUserMsgID(ctx, msgID, reqBody.FromUserName, reqBody.CreateTime)
	if err != nil {
		log.Errorf("handlerUnSubscribeEvent WXRepository wx repo isExistUserMsgID traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	if exist {
		return "", nil
	}
	// 删除用户信息
	u := entity.User{
		OpenID: reqBody.FromUserName,
	}
	err = a.user.DelUser(ctx, u)
	if err != nil {
		log.Errorf("handlerUnSubscribeEvent WXRepository wx repo SaveUser traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	err = a.wx.SetMsgIDToRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerUnSubscribeEvent WXRepository wx repo set msg id to redis failed,traceID:%s,err:%+v", traceID, err)
	}
	return consts.UnSubscribeRespContent, nil
}

func (a *WXRepository) handlerTEMPLATESENDJOBFINISHEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("handlerTEMPLATESENDJOBFINISHEvent traceID:%s", traceID)
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	exist, err := a.isExistTemplateSendJobMsgID(ctx, a.genMsgID(reqBody.FromUserName, reqBody.CreateTime), reqBody.FromUserName, reqBody.CreateTime)
	if err != nil {
		log.Errorf("handlerTEMPLATESENDJOBFINISHEvent WXRepository wx repo isExistTemplateSendJobMsgID traceID:%s,err:%+v", traceID, err)
		return "", err
	}
	if exist {
		return "", nil
	}
	// 对事件推送由于其他原因发送失败的消息进行重发
	// 判断当前发送次数是否小于等于最大重发次数,是则重发
	msg, err := a.msg.GetMsgLogByMsgID(ctx, reqBody.MsgID)
	if err != nil {
		log.Errorf("handlerTEMPLATESENDJOBFINISHEvent GetMsgLogByMsgID failed,traceID:%s,err:%+v", traceID, err)
		return consts.TEMPLATESENDJOBFINISHRespContent, err
	}
	if reqBody.Status == consts.TemplateSendSuccessStatus { // 发送成功，改状态
		updateItem := entity.MsgLog{
			ID:         msg.ID,
			Status:     consts.SendSuccess,
			Cause:      consts.TemplateSendSuccessStatus,
			UpdateTime: reqBody.CreateTime,
		}
		err = a.msg.UpdateMsgLog(ctx, updateItem)
		if err != nil {
			log.Errorf("handlerTEMPLATESENDJOBFINISHEvent UpdateMsgLog TemplateSendSuccessStatus failed,traceID:%s,err:%+v", traceID, err)
			return consts.TEMPLATESENDJOBFINISHRespContent, err
		}
	} else if reqBody.Status == consts.TemplateSendUserBlockStatus { // 发送失败，用户拒接
		updateItem := entity.MsgLog{
			ID:         msg.ID,
			Status:     consts.SendFailure,
			Cause:      consts.TemplateSendUserBlockStatus,
			UpdateTime: reqBody.CreateTime,
		}
		err = a.msg.UpdateMsgLog(ctx, updateItem)
		if err != nil {
			log.Errorf("handlerTEMPLATESENDJOBFINISHEvent UpdateMsgLog TemplateSendUserBlockStatus failed,traceID:%s,err:%+v", traceID, err)
			return consts.TEMPLATESENDJOBFINISHRespContent, err
		}
	} else if reqBody.Status == consts.TemplateSendFailedStatus { // 发送失败，内部错误，重发，改变发送状态为0，重试次数+1
		if msg.Count < consts.MaxRetryCount {
			// 更新发送状态
			updateItem := entity.MsgLog{
				ID:         msg.ID,
				Cause:      consts.TemplateSendFailedStatus,
				Status:     consts.SendPending,
				Count:      msg.Count + 1,
				UpdateTime: reqBody.CreateTime,
			}
			err = a.msg.UpdateMsgLogSendStatus(ctx, updateItem)
			if err != nil {
				log.Errorf("handlerTEMPLATESENDJOBFINISHEvent UpdateMsgLogSendStatus TemplateSendFailedStatus failed,item:%+v,traceID:%s,err:%+v", updateItem, traceID, err)
				return consts.TEMPLATESENDJOBFINISHRespContent, err
			}
		} else {
			// 改变发送状态为失败
			updateItem := entity.MsgLog{
				ID:         msg.ID,
				Cause:      consts.TemplateSendFailedStatus,
				Status:     consts.SendFailure,
				UpdateTime: reqBody.CreateTime,
			}
			err = a.msg.UpdateMsgLog(ctx, updateItem)
			if err != nil {
				log.Errorf("handlerTEMPLATESENDJOBFINISHEvent UpdateMsgLog send status TemplateSendFailedStatus failed,traceID:%s,err:%+v", traceID, err)
				return consts.TEMPLATESENDJOBFINISHRespContent, err
			}
		}
	}
	return consts.TEMPLATESENDJOBFINISHRespContent, nil
}

func (a *WXRepository) isExistUserMsgID(ctx context.Context, msgID string, fromUserName string, createTime int64) (bool, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserMsgID traceID:%s", traceID)
	exist, err := a.wx.IsExistMsgIDFromRedis(ctx, msgID)
	if err != nil {
		log.Errorf("isExistUserMsgID IsExistMsgIDFromRedis failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	// 若存在返回空串,不存在则持久化存储,并保存msg id 到 redis
	if exist {
		return true, nil
	}
	// 从db上找，存在则返回空串
	exist, err = a.user.IsExistUserMsgFromDB(ctx, fromUserName, createTime)
	if err != nil {
		log.Errorf("isExistUserMsgID IsExistUserFromDB failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	if exist {
		// 回写到redis中
		err = a.wx.SetMsgIDToRedis(ctx, msgID)
		if err != nil {
			log.Errorf("isExistUserMsgID WXRepository wx repo set msg id to redis failed,traceID:%s,err:%+v", traceID, err)
		}
		return true, nil
	}
	return false, nil
}

func (a *WXRepository) isExistTemplateSendJobMsgID(ctx context.Context, msgID string, fromUserName string, createTime int64) (bool, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("isExistTemplateSendJobMsgID traceID:%s", traceID)
	exist, err := a.wx.IsExistMsgIDFromRedis(ctx, msgID)
	if err != nil {
		log.Errorf("isExistTemplateSendJobMsgID IsExistMsgIDFromRedis exist msg id from redis,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	if exist {
		return true, nil
	}
	// db查当前消息是否存在
	exist, err = a.msg.IsExistMsgLogFromDB(ctx, fromUserName, createTime)
	if err != nil {
		log.Errorf("isExistTemplateSendJobMsgID IsExistMsgLogFromDB failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	if exist {
		// 回写到redis
		err = a.wx.SetMsgIDToRedis(ctx, msgID)
		if err != nil {
			log.Errorf("isExistUserMsgID WXRepository wx repo set msg id to redis failed,traceID:%s,err:%+v", traceID, err)
		}
		return true, nil
	}
	return false, nil
}

func (a *WXRepository) makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &entity.TextResponseBody{}
	textResponseBody.FromUserName = a.value2CDATA(fromUserName)
	textResponseBody.ToUserName = a.value2CDATA(toUserName)
	textResponseBody.MsgType = a.value2CDATA("text")
	textResponseBody.Content = a.value2CDATA(content)
	textResponseBody.CreateTime = time.Now().Unix()
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func (a *WXRepository) value2CDATA(v string) entity.CDATAText {
	return entity.CDATAText{Text: "<![CDATA[" + v + "]]>"}
}

func (a *WXRepository) genMsgID(fromUserName string, createTime int64) string {
	return fmt.Sprintf("%s%d", fromUserName, createTime)
}
