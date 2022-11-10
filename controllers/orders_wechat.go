package controllers

import (
	"github.com/gin-gonic/gin"
	cns_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/web3-identity/conflux-pay/services"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

type WechatOrderCtrl struct {
	service services.WechatOrderService
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
// @Router      /orders/wechat [post]
func (w *WechatOrderCtrl) MakeOrder(c *gin.Context) {
	req := services.MakeOrderReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	// TODO: 身份认证 APIKEY -> App
	resp, err := w.service.MakeOrder("cns", req)
	ginutils.RenderResp(c, resp, err)
}

// @Tags        Orders
// @ID          RefreshPayUrl
// @Summary     refresh pay url
// @Description refresh pay url
// @Produce     json
// @Param       trade_no   path     string                            true "trade no"
// @Success     200      {object} services.MakeOrderResp
// @Failure     400        {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500        {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/wechat/refresh-url/{trade_no} [put]
func (w *WechatOrderCtrl) RefreshPayUrl(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.RefreshUrl(trandeNo)
	ginutils.RenderResp(c, o, err)
}

// @Tags        Orders
// @ID          QueryWechatOrderDetail
// @Summary     query order by trade no
// @Description query order by trade no
// @Produce     json
// @Param       trade_no path     string true "trade no"
// @Success     200      {object} models.WechatOrderDetail
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/wechat/{trade_no} [get]
func (w *WechatOrderCtrl) GetOrder(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.GetOrderDetail(trandeNo)
	ginutils.RenderResp(c, o, err)
}

// @Tags        Orders
// @ID          Refund
// @Summary     refund pay
// @Description refund pay
// @Produce     json
// @Param       trade_no path     string true "trade no"
// @Param       refund_req body     services.RefundReq                true "refund_req"
// @Success     200        {object} models.WechatRefundDetail         "refund_detail"
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/wechat/refund/{trade_no} [put]
func (w *WechatOrderCtrl) Refund(c *gin.Context) {
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
// @Param       trade_no path     string true "trade no"
// @Success     200      {object} models.WechatOrderDetail
// @Failure     400      {object} cns_errors.RainbowErrorDetailInfo "Invalid request"
// @Failure     500      {object} cns_errors.RainbowErrorDetailInfo "Internal Server error"
// @Router      /orders/wechat/close/{trade_no} [put]
func (w *WechatOrderCtrl) Close(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.Close(trandeNo)
	ginutils.RenderResp(c, o, err)
}

func (w *WechatOrderCtrl) ReceivePayNotify(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	err := w.service.PayNotifyHandler(trandeNo, c.Request)
	ginutils.RenderResp(c, nil, err)
}

func (w *WechatOrderCtrl) ReceiveRefundNotify(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	err := w.service.RefundNotifyHandler(trandeNo, c.Request)
	ginutils.RenderResp(c, nil, err)
}
