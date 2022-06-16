package application

import (
	"context"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
)

type wxApp struct {
	wx repository.WXRepository
}

// wxApp implements the WXInterface
var _ WXInterface = &wxApp{}

type WXInterface interface {
	HandleXML(ctx context.Context, reqBody *entity.TextRequestBody) (respBody []byte, err error)
}

func (w *wxApp) HandleXML(ctx context.Context, reqBody *entity.TextRequestBody) (respBody []byte, err error) {
	return w.wx.HandleXML(ctx, reqBody)
}
