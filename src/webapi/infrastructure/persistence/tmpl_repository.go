package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/pkg/httputil"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	log "github.com/sirupsen/logrus"
)

type TmplRepo struct {
}

var defaultTmplRepo *TmplRepo

func NewTmplRepo() {
	if defaultTmplRepo == nil {
		defaultTmplRepo = &TmplRepo{}
	}
}

func DefaultTmplRepo() *TmplRepo {
	return defaultTmplRepo
}

func (t *TmplRepo) ListTmplFromRequest(ctx context.Context, ak string) (entity.TemplateList, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("ListTmplFromRequest traceID:%s", traceID)
	requestProperty := httputil.GetRequestProperty(http.MethodPost, config.WXListTmplURL+fmt.Sprintf("?access_token=%s", ak),
		nil, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("ListTmplFromRequest request wx list tmpl send failed, traceID:%s, error:%+v", traceID, err)
		return entity.TemplateList{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("ListTmplFromRequest request wx list tmpl send failed, statusCode:%d,traceID:%s, error:%+v", statusCode, traceID, err)
		return entity.TemplateList{}, err
	}
	var msgResp entity.TemplateList
	err = json.Unmarshal(body, &msgResp)
	if err != nil {
		log.Errorf("ListTmplFromRequest get wx list tmpl send failed by unmarshal, resp:%s, traceID:%s, err:%+v", string(body), traceID, err)
		return entity.TemplateList{}, err
	}
	return msgResp, nil
}
