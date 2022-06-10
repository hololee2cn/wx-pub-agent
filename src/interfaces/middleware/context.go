package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/hololee2cn/pkg/ginx"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/wxpub/v1/src/consts"
	"github.com/hololee2cn/wxpub/v1/src/utils"
	log "github.com/sirupsen/logrus"
)

func GinContext(ctx *gin.Context) {
	traceID := ctx.GetHeader(ginx.HTTPTraceIDHeader)
	timeoutStr := ctx.GetHeader(consts.HTTPTimeoutHeader)
	if !validTraceID(traceID) {
		var err error
		log.Warnf("Request %s doesn't input a trace id", ctx.Request.URL.Path)
		traceID, err = utils.GetUUID()
		if err != nil {
			log.Errorf("Request %s new uuid failed, err:%s", ctx.Request.URL.Path, err.Error())
		}
	}
	timeoutSec, _ := strconv.Atoi(timeoutStr)
	if timeoutSec < 1 || timeoutSec > consts.DefaultHTTPTimeOut {
		log.Warnf("Request %s doesn't input a timeout argument or it's invalid: %s", ctx.Request.URL.Path, timeoutStr)
		timeoutSec = consts.DefaultHTTPTimeOut
	}

	c, cancelF := context.WithTimeout(context.WithValue(context.Background(), ginx.ContextTraceID, traceID), time.Second*time.Duration(timeoutSec))
	defer cancelF()
	ctx.Set(ginx.GinContextContext, c)
	ctx.Next()
}

func validTraceID(id string) bool {
	return len(id) > 0
}
