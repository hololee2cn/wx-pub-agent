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
}

func (u *tmplApp) ListTemplate(ctx context.Context) (entity.ListTmplResp, error) {
	return u.tmpl.ListTemplate(ctx)
}
