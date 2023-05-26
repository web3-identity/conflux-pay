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
			// trader := order.Group(":provider")
			ctrl := controllers.NewOrderCtrl()
			order.POST("", ctrl.MakeOrder)
			order.GET("/:trade_no", ctrl.GetOrder) //provider maybe trader/alipay
			order.PUT("/refresh-url/:trade_no", ctrl.RefreshPayUrl)
			order.PUT("/refund/:trade_no", ctrl.Refund)
			order.PUT("/close/:trade_no", ctrl.Close)

			order.POST("/notify-pay/:trade_no", ctrl.ReceivePayNotify)
			order.POST("/notify-refund/:trade_no", ctrl.ReceiveRefundNotify)

			// alipay := order.Group("alipay")
			// {
			// 	ctrl := controllers.AlipayOrderCtrl{}
			// 	alipay.POST("", ctrl.MakeOrder)
			// 	alipay.POST("/", ctrl.MakeOrder)
			// 	alipay.GET("/:trade_no", ctrl.GetOrder)
			// }
		}
		cmb := api.Group("cmb")
		{
			cmb.GET("/history", controllers.QueryCmbRecords)
			cmb.GET("/history/recent", controllers.QueryRecentCmbRecords)
			cmb.POST("/unit-account", controllers.AddUnitAccount)
			cmb.POST("/unit-account/relation", controllers.SetUnitAccountRelation)
		}
	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CONFLUX_PAY"))
}
