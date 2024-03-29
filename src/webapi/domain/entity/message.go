package entity

import (
	"encoding/json"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/vo"

	"github.com/hololee2cn/wxpub/v1/src/utils"
)

type SendTmplMsgReq struct {
	// 模板id
	TmplMsgID string `json:"tmpl_msg_id"`
	// 接收者手机号
	Phones []string `json:"phones"`
	// 模板数据
	Data MsgData `json:"data"`
	// 跳转详情url
	URL string `json:"url"`
}

// MsgData 微信模板Data，微信规定最多只有5个keyword
type MsgData struct {
	First    DataItem `json:"first"`
	Keyword1 DataItem `json:"keyword1"`
	Keyword2 DataItem `json:"keyword2"`
	Keyword3 DataItem `json:"keyword3"`
	Keyword4 DataItem `json:"keyword4"`
	Keyword5 DataItem `json:"keyword5"`
	Remark   DataItem `json:"remark"`
}

type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type SendTmplMsgResp struct {
	// 发送请求消息id
	RequestID string `json:"request_id"`
}

// TmplMsgStatusReq 模板消息状态请求
type TmplMsgStatusReq struct {
	RequestID string `json:"request_id"`
}

// TmplMsgStatusResp 模板消息状态返回
type TmplMsgStatusResp struct {
	Lists []TmplMsgStatusItem `json:"lists"`
	Total int                 `json:"total"`
}

type TmplMsgStatusItem struct {
	// 主键id
	ID int `json:"id"`
	// 用户号码
	Phone string `json:"phone"`
	// 模板id
	TemplateID string `json:"template_id"`
	// 发送模板内容
	Content MsgData `json:"content"`
	// 跳转详情url
	URL string `json:"url"`
	// 失败原因
	Cause string `json:"cause"`
	// 发送状态，0为pending判定中，1为sending发送中，2为success成功，3为failure失败
	Status int `json:"status"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
}

type SendTmplMsgRemoteReq struct {
	// 获取到的凭证
	AccessToken string `json:"access_token"`
	// 接收者openid
	ToUser string `json:"touser"`
	// 模板ID
	TemplateID string `json:"template_id"`
	// 模板数据
	Data MsgData `json:"data"`
	// 跳转详情url
	URL string `json:"url"`
}

type SendTmplMsgRemoteResp struct {
	MsgID int64 `json:"msgid"`
	ErrorInfo
}

// MsgLog 消息日志表
type MsgLog struct {
	// 主键id
	ID int `json:"id" gorm:"id"`
	// 发送消息id
	RequestID string `json:"request_id" gorm:"request_id"`
	// 微信消息id
	MsgID int64 `json:"msg_id" gorm:"msg_id"`
	// 接收者openid
	ToUser string `json:"to_user" gorm:"to_user"`
	// 用户号码
	Phone string `json:"phone" gorm:"phone"`
	// 模板id
	TemplateID string `json:"template_id" gorm:"template_id"`
	// 发送模板内容
	Content string `json:"content" gorm:"content"`
	// 跳转详情url
	URL string `json:"url" gorm:"url"`
	// 失败原因
	Cause string `json:"cause" gorm:"cause"`
	// 发送状态，0为pending判定中，1为sending发送中，2为success成功，3为failure失败
	Status int `json:"status" gorm:"status"`
	// 发送次数
	Count int `json:"count" gorm:"count"`
	// 创建时间
	CreateTime int64 `json:"create_time" gorm:"create_time"`
	// 更新时间
	UpdateTime int64 `json:"update_time" gorm:"update_time"`
}

func (m *TmplMsgStatusReq) Validate() string {
	if len(m.RequestID) <= 0 {
		return "request_id is empty"
	}
	return ""
}

func (m MsgLog) TableName() string {
	return "msg_log"
}

func (m *MsgLog) TransferSendTmplMsgRemoteReq() (SendTmplMsgRemoteReq, error) {
	var wxContent MsgData
	err := json.Unmarshal([]byte(m.Content), &wxContent)
	if err != nil {
		return SendTmplMsgRemoteReq{}, err
	}
	return SendTmplMsgRemoteReq{
		ToUser:     m.ToUser,
		TemplateID: m.TemplateID,
		Data:       wxContent,
		URL:        m.URL,
	}, nil
}

func (m *MsgLog) TransferTmplMsgStatusItem() (TmplMsgStatusItem, error) {
	var wxContent MsgData
	err := json.Unmarshal([]byte(m.Content), &wxContent)
	if err != nil {
		return TmplMsgStatusItem{}, err
	}
	return TmplMsgStatusItem{
		ID:         m.ID,
		Phone:      m.Phone,
		TemplateID: m.TemplateID,
		Content:    wxContent,
		URL:        m.URL,
		Cause:      m.Cause,
		Status:     m.Status,
		CreateTime: m.CreateTime,
	}, nil
}

func (r *SendTmplMsgReq) Validate() (errorMsg string) {
	if len(r.TmplMsgID) <= 0 {
		errorMsg = "tmpl msg id is empty"
		return
	}
	if len(r.Phones) <= 0 {
		errorMsg = "phones is empty"
		return
	}
	// 去重
	r.Phones = utils.RemoveStringRepeated(r.Phones)
	if len(r.Phones) > 100 {
		errorMsg = "the phones slice is too long"
		return
	}
	return
}

func (r *SendTmplMsgReq) TransferPendingMsgLog(requestID string, toUser string, phone string, content string) MsgLog {
	return MsgLog{
		RequestID:  requestID,
		ToUser:     toUser,
		Phone:      phone,
		TemplateID: r.TmplMsgID,
		Content:    content,
		URL:        r.URL,
		Cause:      "",
		Status:     new(vo.MsgStatus).GetPending(),
		Count:      0,
		CreateTime: time.Now().Unix(),
		UpdateTime: 0,
	}
}

func (r *SendTmplMsgReq) TransferFailureMsgLog(requestID string, toUser string, phone string, content string) MsgLog {
	return MsgLog{
		RequestID:  requestID,
		ToUser:     toUser,
		Phone:      phone,
		TemplateID: r.TmplMsgID,
		Content:    content,
		URL:        r.URL,
		Cause:      "",
		Status:     new(vo.MsgStatus).GetFailure(),
		Count:      0,
		CreateTime: time.Now().Unix(),
		UpdateTime: 0,
	}
}
