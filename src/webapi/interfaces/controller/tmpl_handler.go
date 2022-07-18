package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/errorx"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/application"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
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

func (a *Tmpl) GetTemplate(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)
	var param entity.GetTemplateReq
	param.TemplateID = ginx.URLParamStr(c, "id")
	template, err := a.tmpl.GetTemplate(ctx, param.TemplateID)
	if template == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such template")
		return
	}
	ginx.NewRender(c).Data(template, err)
}

func (a *Tmpl) FreshTemplate(c *gin.Context) {
	ctx := ginx.DefaultTodoContext(c)
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)
	var param entity.FreshTemplateReq
	ginx.BindJSON(c, &param)
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("FreshTemplate validate FreshTemplateReq req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		ginx.BombErr(errorx.CodeInvalidParams, errMsg)
	}
	_, err := a.tmpl.FreshTemplate(ctx)
	ginx.NewRender(c).Data(nil, err)
}
