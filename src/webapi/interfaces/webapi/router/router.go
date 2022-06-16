package router

import (
	"fmt"
	"github.com/hololee2cn/wxpub/v1/src/webapi/config"
	"github.com/hololee2cn/wxpub/v1/src/webapi/consts"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/repository"
	"github.com/hololee2cn/wxpub/v1/src/webapi/interfaces/controller"
	"github.com/hololee2cn/wxpub/v1/src/webapi/interfaces/middleware"
	"github.com/hololee2cn/wxpub/v1/src/webapi/wxutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/utils"
)

var (
	wx   *controller.WX
	user *controller.User
	msg  *controller.Message
)

func registerController() {
	wx = controller.NewWXController(
		repository.DefaultWXRepository())
	user = controller.NewUserController(
		repository.DefaultUserRepository())
	msg = controller.NewMessageController(
		repository.DefaultMessageRepository())
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
	open := router.Group("/api/v1/")
	// web
	routerWeb(open)
	// wx api
	routerWX(open)
	// user info verify and binding
	routerVerify(open)

	router.Use(middleware.GinContext)

	// msg handler
	routerMsg(open)
}

func routerWeb(open *gin.RouterGroup) {
	open.Static("/static", "doc")
	open.StaticFile("/favicon.icon", "doc/img/favicon.icon")
	open.GET("/", func(context *gin.Context) {
		ts := time.Now().Unix()
		nonce := utils.GenRandNonce()
		sig := wxutil.CalcSign(fmt.Sprintf("%d", ts), nonce, consts.Token)
		context.HTML(http.StatusOK, "index.html", gin.H{
			"app_id":      config.AppID,
			"timestamp":   ts,
			"nonce_str":   nonce,
			"signature":   sig,
			"js_api_list": []string{""},
		})
	})
}

func routerWX(router *gin.RouterGroup) {
	wxGroup := router.Group("/wx")
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
		smsProfileGroup.GET("/send-sms", user.SendSms)
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
