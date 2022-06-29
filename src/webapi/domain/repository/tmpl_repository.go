package repository

import (
	"context"

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
	ak, err := DefaultAccessTokenRepository().GetAccessToken(ctx)
	if err != nil {
		log.Errorf("TmplRepository ListTemplate get ak failed,traceID: %s,err: %+v", traceID, err)
		ginx.CustomErr(err)
	}
	templateList, err := t.tmpl.ListTmplFromRequest(ctx, ak)
	if err != nil {
		log.Errorf("TmplRepository ListTemplate request list template failed,traceID:%s,err:%+v", traceID, err)
		ginx.CustomErr(err)
	}
	return templateList.TransferListTmplResp(), nil
}
