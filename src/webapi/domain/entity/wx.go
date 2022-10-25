package entity

import (
	"encoding/xml"
)

type MsgType string

type WXCheckReq struct {
	Signature string `json:"signature" form:"signature"`
	TimeStamp string `json:"time_stamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	EchoStr   string `json:"echo_str" form:"echostr"`
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      MsgType
	Content      string
	MsgID        int64
	Event        string
	Ticket       string
	EventKey     string
	Status       string
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   int64
	MsgType      CDATAText
	Content      CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

func (u *WXCheckReq) Validate() (errorMessage string) {
	return
}
