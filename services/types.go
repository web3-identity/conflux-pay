package services

import (
	"encoding/json"

	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type TradeStateChangeHandler func(o *models.Order)
type RefundStateChangeHandler func(o *models.Order)

type MakeOrderReq struct {
	TradeProvider string          `json:"trade_provider" swaggertype:"string" binding:"required,oneof=wechat alipay"`
	TradeType     enums.TradeType `json:"trade_type" binding:"required" swaggertype:"string"`
	Description   *string         `json:"description" binding:"required"`
	TimeExpire    int64           `json:"time_expire,omitempty" binding:"required"` // alipay 当面付无效，当面付固定过期时间为2小时
	Amount        int64           `json:"amount" binding:"required"`
	NotifyUrl     *string         `json:"notify_url,omitempty"`
	QrPayMode     string          `json:"qr_mode,omitempty"`    // 只有alipay，且 trade type 为 h5 模式有效; 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	ReturnUrl     string          `json:"return_url,omitempty"` // 只有alipay，且 trade type 为 h5 模式有效
}

func (m *MakeOrderReq) MustGetTradeProvider() enums.TradeProvider {
	val, ok := enums.ParseTradeProviderByName(m.TradeProvider)
	if !ok {
		panic("unkown trade provider")
	}
	return *val
}

type MakeOrderResp struct {
	TradeProvider enums.TradeProvider `json:"trade_provider" swaggertype:"string"`
	TradeType     enums.TradeType     `json:"trade_type" swaggertype:"string"`
	TradeNo       string              `json:"trade_no"`
	CodeUrl       *string             `json:"code_url,omitempty"`
	H5Url         *string             `json:"h5_url,omitempty"`
}

func NewMakeOrderRespFromRaw(raw *models.Order) *MakeOrderResp {
	return &MakeOrderResp{
		TradeProvider: raw.Provider,
		TradeType:     raw.TradeType,
		TradeNo:       raw.TradeNo,
		CodeUrl:       raw.CodeUrl,
		H5Url:         raw.H5Url,
	}
}

type RefundReq struct {
	Reason    string  `json:"reason" binding:"required"`
	NotifyUrl *string `json:"notify_url"`
}

type refundWithRefundStatus struct {
	refunddomestic.Refund
	RefundStatus *string `gorm:"-" json:"refund_status"`
}

type PrepayRequest native.PrepayRequest

func (r *PrepayRequest) toH5() (val *h5.PrepayRequest) {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(j, &val); err != nil {
		panic(err)
	}
	return
}