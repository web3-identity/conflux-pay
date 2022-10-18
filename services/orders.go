package services

import (
	"github.com/wangdayong228/conflux-pay/models"
	"github.com/wangdayong228/conflux-pay/models/enums"
	cns_errors "github.com/wangdayong228/conflux-pay/pay_errors"
)

type OrderService struct {
	wechatService *WechatOrderService
	// alipayService AlipayOrderService
}

func NewOrderService() *OrderService {
	return &OrderService{
		wechatService: &WechatOrderService{},
	}
}

// 数据库查询，如果状态未稳定，wechatpay查询
// TODO: 优化策略
// 1. 在微信主动通知平均时间内直接返回NOTPAY,之后如果还未收到notify，主动查询
// 2. 或者 开启异步同步结果服务
func (w *OrderService) GetOrderSummary(tradeNo string) (*models.Order, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}
	if o.TradeState.IsStable() {
		return o, nil
	}

	switch o.Provider {
	case enums.TRADE_PROVIDER_WECHAT:
		w.wechatService.GetOrderDetailAndSave(tradeNo)
		return models.FindOrderByTradeNo(tradeNo)
	}

	return nil, cns_errors.ERR_PROVIDER_UNSUPPORT
}
