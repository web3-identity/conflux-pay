package controllers

import (
	"github.com/gin-gonic/gin"
	cns_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/web3-identity/conflux-pay/services"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

type OrderCtrl struct {
	service services.OrderService
}

func NewOrderCtrl() *OrderCtrl {
	return &OrderCtrl{
		service: *services.NewOrderService(),
	}
}

// @Tags        Orders
// @ID          MakeOrder
// @Summary     Make Order
// @Description make order
// @Produce     json
// @Param       make_ord_req body     services.MakeOrderReq true "make_wechat_order_req"
// @Success     200          {object} models.Order
// @Failure     400          {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500          {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders [post]
func (w *OrderCtrl) MakeOrder(c *gin.Context) {
	req := services.MakeOrderReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	// TODO: 身份认证 APIKEY -> App
	resp, err := w.service.MakeOrder(req)
	ginutils.RenderResp(c, resp, err)
}

// @Tags        Orders
// @ID          RefreshPayUrl
// @Summary     refresh pay url
// @Description refresh pay url
// @Produce     json
// @Param       trade_no path     string true "trade no"
// @Success     200      {object} services.MakeOrderResp
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/refresh-url/{trade_no} [put]
func (w *OrderCtrl) RefreshPayUrl(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.RefreshUrl(trandeNo)
	ginutils.RenderResp(c, o, err)
}

// @Tags        Orders
// @ID          QueryOrder
// @Summary     query order by trade no
// @Description query order by trade no
// @Produce     json
// @Param       trade_no path     string                            true "trade no"
// @Success     200      {object} models.Order                      "order"
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/{trade_no} [get]
func (w *OrderCtrl) GetOrder(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.GetOrder(trandeNo)
	ginutils.RenderResp(c, o, err)
}

// @Tags        Orders
// @ID          Refund
// @Summary     refund pay
// @Description refund pay
// @Produce     json
// @Param       trade_no   path     string                            true "trade no"
// @Param       refund_req body     services.RefundReq                true "refund_req"
// @Success     200        {object} models.OrderCore                  "order"
// @Failure     400        {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500        {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/refund/{trade_no} [put]
func (w *OrderCtrl) Refund(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	req := services.RefundReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	o, err := w.service.Refund(trandeNo, req)
	ginutils.RenderResp(c, o, err)
}

// @Tags        Orders
// @ID          Close
// @Summary     close order
// @Description close order
// @Produce     json
// @Param       trade_no path     string                            true "trade no"
// @Success     200      {object} models.OrderCore                  "order"
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/close/{trade_no} [put]
func (w *OrderCtrl) Close(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.Close(trandeNo)
	ginutils.RenderResp(c, o, err)
}

func (w *OrderCtrl) ReceivePayNotify(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	err := w.service.PayNotifyHandler(trandeNo, c.Request)
	ginutils.RenderResp(c, nil, err)
}

func (w *OrderCtrl) ReceiveRefundNotify(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	err := w.service.RefundNotifyHandler(trandeNo, c.Request)
	ginutils.RenderResp(c, nil, err)
}
