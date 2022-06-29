package entity

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

func (m *TemplateList) TransferListTmplResp() ListTmplResp {
	var ret ListTmplResp
	for _, v := range m.TemplateList {
		ret.Lists = append(ret.Lists, ListTmplItem{
			TemplateID: v.TemplateID,
			Title:      v.Title,
			Content:    v.Content,
			Example:    v.Example,
		})
	}
	ret.Total = len(m.TemplateList)
	return ret
}
