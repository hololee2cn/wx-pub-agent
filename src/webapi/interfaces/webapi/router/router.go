package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	"github.com/hololee2cn/wxpub/v1/src/webapi/interfaces/controller"
	"github.com/hololee2cn/wxpub/v1/src/webapi/interfaces/middleware"
)

var (
	wx   *controller.WX
	user *controller.User
	msg  *controller.Message
	tmpl *controller.Tmpl
)

func registerController() {
	wx = controller.NewWXController(
		repository.DefaultWXRepository())
	user = controller.NewUserController(
		repository.DefaultUserRepository())
	msg = controller.NewMessageController(
		repository.DefaultMessageRepository())
	tmpl = controller.NewTmplController(
		repository.DefaultTmplRepository())
}

func New() *gin.Engine {
	gin.SetMode(string(config.SMode))

	if strings.ToLower(string(config.SMode)) == gin.ReleaseMode {
		ginx.DisableConsoleColor()
	}
	registerController()
	engine := gin.New()
	engine.Use(ginx.Recovery())
	initRouter(engine)

	return engine
}

func initRouter(router *gin.Engine) {
	root := router.Group("")
	// web
	routerWeb(root)

	router.Use(middleware.GinContext)
	api := root.Group("/api/v1/")
	// wx api
	routerWX(api)
	// user info verify and binding
	routerVerify(api)
	// msg handler
	routerMsg(api)
	// template management api
	routerTmpl(api)
}

func routerWeb(open *gin.RouterGroup) {
	open.Static("/static", "doc/img")
	open.StaticFile("/favicon.icon", "static/favicon.icon")
	open.StaticFile("", "doc/static/index.html")
}

func routerWX(router *gin.RouterGroup) {
	wxGroup := router.Group("/wxapi")
	{
		// wx开放平台接入测试接口
		wxGroup.GET("", wx.GetWXCheckSign)
		// todo: 暂时先用明文传输，后续补充aes加密传输
		// wx开放平台事件接收
		wxGroup.POST("", wx.HandleXML)
	}
}

func routerVerify(router *gin.RouterGroup) {
	smsProfileGroup := router.Group("/user")
	{
		smsProfileGroup.POST("/send-sms", user.SendSms)
		smsProfileGroup.POST("/verify-sms", user.VerifyAndUpdatePhone)
		smsProfileGroup.GET("/captcha", user.GenCaptcha)
	}
}

func routerMsg(router *gin.RouterGroup) {
	msgGroup := router.Group("/message")
	{
		// tmpl msg pusher
		pushSubGroup := msgGroup.Group("/tmpl-push")
		{
			pushSubGroup.POST("", msg.SendTmplMessage)
		}
		// tmpl msg status
		statusSubGroup := msgGroup.Group("/status")
		{
			statusSubGroup.GET("/:id", msg.TmplMsgStatus)
		}
	}
}

func routerTmpl(router *gin.RouterGroup) {
	tmplGroup := router.Group("/tmpl")
	{
		tmplGroup.GET("", tmpl.ListTemplate)
	}
}
