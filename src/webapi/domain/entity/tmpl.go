package entity

import "github.com/hololee2cn/wxpub/v1/src/webapi/config"

// ListTmplResp 模板消息状态返回
type ListTmplResp struct {
	Lists []ListTmplItem `json:"lists"`
	Total int            `json:"total"`
}

// ListTmplItem 模板列表返回
type ListTmplItem struct {
	TemplateID string `json:"template_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Example    string `json:"example"`
}

type TemplateList struct {
	TemplateList []TemplateItem `json:"template_list"`
}

type TemplateItem struct {
	TemplateID      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

// GetTemplateReq 模板内容请求
type GetTemplateReq struct {
	TemplateID string `json:"template_id"`
}

// FreshTemplateReq 刷新模板请求
type FreshTemplateReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func (m *TemplateList) TransferListTmplResp() ListTmplResp {
	var ret ListTmplResp
	for _, v := range m.TemplateList {
		ret.Lists = append(ret.Lists, *v.TransferListTmplItem())
	}
	ret.Total = len(m.TemplateList)
	return ret
}

func (m *TemplateItem) TransferListTmplItem() *ListTmplItem {
	return &ListTmplItem{
		TemplateID: m.TemplateID,
		Title:      m.Title,
		Content:    m.Content,
		Example:    m.Example,
	}
}

func (f *FreshTemplateReq) Validate() (errMsg string) {
	if f.AppID != config.Get().WxSvc.AppID {
		errMsg = "appid is not right"
	}
	if f.AppSecret != config.Get().WxSvc.AppSecret {
		errMsg = "app secret is not right"
	}
	return
}
