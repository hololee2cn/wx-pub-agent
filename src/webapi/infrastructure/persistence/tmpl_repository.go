package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	redis2 "github.com/hololee2cn/wxpub/v1/src/pkg/redis"

	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/pkg/httputil"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	log "github.com/sirupsen/logrus"
)

type TmplRepo struct {
	Redis *redis.UniversalClient
}

var defaultTmplRepo *TmplRepo

func NewTmplRepo() {
	if defaultTmplRepo == nil {
		defaultTmplRepo = &TmplRepo{
			Redis: CommonRepositories.Redis,
		}
	}
}

func DefaultTmplRepo() *TmplRepo {
	return defaultTmplRepo
}

func (t *TmplRepo) ListTmplFromRequest(ctx context.Context, ak string) (entity.TemplateList, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("ListTmplFromRequest traceID:%s", traceID)
	requestProperty := httputil.GetRequestProperty(http.MethodPost, fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=%s", ak),
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

func (t *TmplRepo) SetTmplsToRedis(ctx context.Context, val entity.ListTmplResp) error {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("SetMsgIDToRedis traceID:%s", traceID)
	var err error
	bs, err := json.Marshal(val)
	if err != nil {
		log.Errorf("SetTmplsToRedis json marshal failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	for i := 0; i < 3; i++ {
		err = redis2.RSet(config.RedisKeyTmpl, bs, config.RedisTmplTTL)
		if err != nil {
			log.Errorf("SetTmplsToRedis TmplRepo redis set tmpls failed,traceID:%s,err:%+v", traceID, err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}
	return err
}

func (t *TmplRepo) GetTmplsForRedis(ctx context.Context) (entity.ListTmplResp, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("GetTmplsForRedis traceID:%s", traceID)
	bs, err := redis2.RGet(config.RedisKeyTmpl)
	if err != nil {
		log.Errorf("TmplRepository TmplRepo get tmpls for redis failed,err:+%v", err)
		return entity.ListTmplResp{}, err
	}
	var tmpls entity.ListTmplResp
	err = json.Unmarshal(bs, &tmpls)
	if err != nil {
		log.Errorf("TmplRepository TmplRepo json unmarshal failed,err:+%v", err)
		return entity.ListTmplResp{}, err
	}
	return tmpls, nil
}
