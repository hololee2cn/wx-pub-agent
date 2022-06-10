package main

import (
	"context"

	"github.com/hololee2cn/wxpub/v1/src/interfaces/webapi"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func main() {
	globalCtx, globalCancel = context.WithCancel(context.Background())
	webapi.Run(globalCtx, globalCancel)
}
