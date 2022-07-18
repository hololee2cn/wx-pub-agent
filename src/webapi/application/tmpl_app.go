package application

import (
	"context"

	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
)

type tmplApp struct {
	tmpl repository.TmplRepository
}

// tmplApp implements the TmplInterface
var _ TmplInterface = &tmplApp{}

type TmplInterface interface {
	ListTemplate(ctx context.Context) (entity.ListTmplResp, error)
	GetTemplate(ctx context.Context, templateID string) (*entity.ListTmplItem, error)
	FreshTemplate(ctx context.Context) (entity.ListTmplResp, error)
}

func (u *tmplApp) ListTemplate(ctx context.Context) (entity.ListTmplResp, error) {
	return u.tmpl.ListTemplate(ctx)
}

func (u *tmplApp) GetTemplate(ctx context.Context, templateID string) (*entity.ListTmplItem, error) {
	return u.tmpl.GetTemplate(ctx, templateID)
}

func (u *tmplApp) FreshTemplate(ctx context.Context) (entity.ListTmplResp, error) {
	return u.tmpl.FreshTemplate(ctx)
}
