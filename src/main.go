package main

import (
	"context"
	"github.com/hololee2cn/wxpub/v1/src/captcha"
	"github.com/hololee2cn/wxpub/v1/src/webapi"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
	VERSION      = "1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "wxpub"
	app.Version = VERSION
	app.Usage = "wxpub, wechat public platform backend server"
	app.Commands = []*cli.Command{
		newCaptchaCmd(),
		newWebapiCmd(),
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func newWebapiCmd() *cli.Command {
	return &cli.Command{
		Name:  "webapi",
		Usage: "run webapi",
		Action: func(ctx *cli.Context) error {
			globalCtx, globalCancel = context.WithCancel(context.Background())
			webapi.Run(globalCtx, globalCancel)
			return nil
		},
	}
}

func newCaptchaCmd() *cli.Command {
	return &cli.Command{
		Name:  "captcha",
		Usage: "run captcha",
		Action: func(ctx *cli.Context) error {
			captcha.Run()
			return nil
		},
	}
}
