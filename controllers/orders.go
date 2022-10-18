package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/conflux-pay/services"
	"github.com/wangdayong228/conflux-pay/utils/ginutils"
)

var (
	orderService = services.NewOrderService()
)

func GetOrderSummary(c *gin.Context) {
	tradeNo := c.Param("trade_no")
	o, err := orderService.GetOrderSummary(tradeNo)
	ginutils.RenderResp(c, o, err)
}
