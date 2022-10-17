package controllers

import (
	"github.com/gin-gonic/gin"
	cns_errors "github.com/wangdayong228/conflux-pay/pay_errors"
	"github.com/wangdayong228/conflux-pay/services"
	"github.com/wangdayong228/conflux-pay/utils/ginutils"
)

type WechatOrderCtrl struct {
	service services.WechatOrderService
}

func (w *WechatOrderCtrl) MakeOrder(c *gin.Context) {
	req := services.MakeWechatOrderReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
	}
	// TODO: 身份认证 APIKEY -> App
	resp, err := w.service.MakeOrder("cns", req)
	ginutils.RenderResp(c, resp, err)
}

func (w *WechatOrderCtrl) GetOrder(c *gin.Context) {
	trandeNo := c.Param("trade_no")
	o, err := w.service.GetOrderDetail(trandeNo)
	ginutils.RenderResp(c, o, err)
}
