package repository

import (
	"context"
	"fmt"

	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/infrastructure/persistence"
	log "github.com/sirupsen/logrus"
)

type TmplRepository struct {
	tmpl *persistence.TmplRepo
}

var defaultTmplRepository = &TmplRepository{}

func NewTmplRepository(tmpl *persistence.TmplRepo) {
	if defaultTmplRepository.tmpl == nil {
		defaultTmplRepository.tmpl = tmpl
	}
}

func DefaultTmplRepository() *TmplRepository {
	return defaultTmplRepository
}

func (t *TmplRepository) ListTemplate(ctx context.Context) (entity.ListTmplResp, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	tmpls, err := t.tmpl.GetTmplsForRedis(ctx)
	if err != nil {
		log.Errorf("TmplRepository GetTmplsForRedis failed,traceID:%s,err:%+v", traceID, err)
	}
	if len(tmpls.Lists) > 0 {
		return tmpls, nil
	}
	tmplResp, err := t.FreshTemplate(ctx)
	if err != nil {
		log.Errorf("TmplRepository StoreAndGetTmpl failed,traceID:%s,err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	return tmplResp, nil
}

func (t *TmplRepository) GetTemplate(ctx context.Context, templateID string) (*entity.ListTmplItem, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	tmpls, err := t.tmpl.GetTmplsForRedis(ctx)
	if err != nil {
		log.Errorf("TmplRepository GetTmplsForRedis failed,traceID:%s,err:%+v", traceID, err)
	}
	if len(tmpls.Lists) > 0 {
		return t.findTmplByID(tmpls.Lists, templateID)
	}
	tmplResp, err := t.FreshTemplate(ctx)
	if err != nil {
		log.Errorf("TmplRepository StoreAndGetTmpl failed,traceID:%s,err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	return t.findTmplByID(tmplResp.Lists, templateID)
}

func (t *TmplRepository) FreshTemplate(ctx context.Context) (tmplResp entity.ListTmplResp, err error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	ak, err := DefaultAccessTokenRepository().GetAccessToken(ctx)
	if err != nil {
		log.Errorf("storeTmpl GetAccessToken failed,traceID:%s,err:%+v", traceID, err)
		return
	}
	tmplList, err := t.tmpl.ListTmplFromRequest(ctx, ak)
	if err != nil {
		log.Errorf("storeTmpl TmplRepository ListTemplate request list template failed,traceID:%s,err:%+v", traceID, err)
		return
	}
	// 缓存模板内容，过期时间为一天
	tmplResp = tmplList.TransferListTmplResp()
	err = t.tmpl.SetTmplsToRedis(ctx, tmplResp)
	if err != nil {
		log.Errorf("storeTmpl SetTmplsToRedis failed,traceID:%s,err:%+v", traceID, err)
		return
	}
	return
}

func (t *TmplRepository) findTmplByID(tmpls []entity.ListTmplItem, tmplID string) (*entity.ListTmplItem, error) {
	for _, tmpl := range tmpls {
		if tmpl.TemplateID == tmplID {
			return &tmpl, nil
		}
	}
	return &entity.ListTmplItem{}, fmt.Errorf("template is not found")
}
