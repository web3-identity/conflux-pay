package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/web3-identity/conflux-pay/controllers"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", indexEndpoint)

	api := router.Group("v0")
	{
		order := api.Group("orders")
		{
			order.GET("summary/:trade_no", controllers.GetOrderSummary) //provider maybe wechat/alipay
			wechat := order.Group("wechat")
			{
				ctrl := controllers.WechatOrderCtrl{}
				wechat.POST("", ctrl.MakeOrder)
				wechat.POST("/", ctrl.MakeOrder)
				wechat.PUT("/refresh-url/:trade_no", ctrl.RefreshPayUrl)
				wechat.PUT("/refund/:trade_no", ctrl.Refund)
				wechat.PUT("/close/:trade_no", ctrl.Close)
				wechat.GET("/:trade_no", ctrl.GetOrder)
				wechat.POST("/notify/:trade_no", ctrl.ReceiveNotify)
			}
			alipay := order.Group("alipay")
			{
				ctrl := controllers.AlipayOrderCtrl{}
				alipay.POST("", ctrl.MakeOrder)
				alipay.POST("/", ctrl.MakeOrder)
				alipay.GET("/:trade_no", ctrl.GetOrder)
			}
		}

	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CONFLUX_PAY"))
}
