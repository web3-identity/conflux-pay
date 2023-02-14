package controllers

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/web3-identity/conflux-pay/services"
// 	"github.com/web3-identity/conflux-pay/utils/ginutils"
// )

// var (
// 	orderService = services.NewWechatOrderService()
// )

// // @Tags        Orders
// // @ID          QueryOrderSummary
// // @Summary     query order summary by trade no
// // @Description query order summary by trade no
// // @Produce     json
// // @Param       trade_no path     string true "trade no"
// // @Success     200      {object} models.Order
// // @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// // @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// // @Router      /orders/summary/{trade_no} [get]
// func GetOrder(c *gin.Context) {
// 	tradeNo := c.Param("trade_no")
// 	o, err := orderService.GetOrder(tradeNo)
// 	ginutils.RenderResp(c, o, err)
// }
