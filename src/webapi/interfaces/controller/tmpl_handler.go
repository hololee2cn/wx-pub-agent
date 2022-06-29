package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	log "github.com/sirupsen/logrus"
)

type Tmpl struct {
	tmpl application.TmplInterface
}

func NewTmplController(tmpl application.TmplInterface) *Tmpl {
	return &Tmpl{
		tmpl: tmpl,
	}
}

func (a *Tmpl) ListTemplate(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	templates, err := a.tmpl.ListTemplate(ctx)
	ginx.NewRender(c).Data(templates, err)
}
