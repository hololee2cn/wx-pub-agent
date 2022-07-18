package tasks

import (
	"context"

	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	log "github.com/sirupsen/logrus"
)

func SyncTemplates(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Infof("sync templates exit,err:%v", ctx.Err())
		return
	default:
		_, err := repository.DefaultTmplRepository().FreshTemplate(ctx)
		if err != nil {
			log.Errorf("SyncTemplates FreshTemplate failed,err:%+v", err)
		}
	}
}
