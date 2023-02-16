package services

import (
	"encoding/json"
	"time"

	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type TradeStateChangeHandler func(o *models.Order)
type RefundStateChangeHandler func(o *models.Order)

type MakeOrderReq struct {
	AppName       string          `json:"app_name" binding:"required"`
	TradeProvider string          `json:"trade_provider" swaggertype:"string" binding:"required,oneof=wechat alipay"`
	TradeType     enums.TradeType `json:"trade_type" binding:"required" swaggertype:"string"`
	Description   *string         `json:"description" binding:"required"`
	TimeExpire    int64           `json:"time_expire,omitempty" binding:"required"` // alipay 当面付无效，当面付固定过期时间为2小时
	Amount        int64           `json:"amount" binding:"required"`
	NotifyUrl     *string         `json:"notify_url,omitempty"`
	QrPayMode     string          `json:"qr_pay_mode,omitempty"`   // 支付二维码模式。 只有alipay，且 trade type 为 h5 模式有效; 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	QrCodeWidth   string          `json:"qr_code_width,omitempty"` // 二维码宽度。 只有alipay，且 trade type 为 h5 模式有效，qr pay mode 为4 时有效； 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	ReturnUrl     string          `json:"return_url,omitempty"`    // 付款成功后的跳转链接。只有alipay，且 trade type 为 h5 模式有效; 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
}

func NewMakeOrderReqFromOrder(o *models.Order) *MakeOrderReq {
	return &MakeOrderReq{
		TradeProvider: o.TradeProvider.String(),
		TradeType:     o.TradeType,
		Description:   o.Description,
		TimeExpire:    o.TimeExpire.Unix(),
		Amount:        int64(o.Amount),
		NotifyUrl:     o.AppPayNotifyUrl,
		QrPayMode:     o.QrPayMode,
		QrCodeWidth:   o.QrCodeWidth,
		ReturnUrl:     o.ReturnUrl,
	}
}

func (req *MakeOrderReq) FillToOrder(o *models.Order) {
	expire := time.Unix(req.TimeExpire, 0)
	o.TradeProvider = req.MustGetTradeProvider()
	o.TradeType = req.TradeType
	o.Description = req.Description
	o.TimeExpire = &expire
	o.Amount = uint(req.Amount)
	o.AppPayNotifyUrl = req.NotifyUrl
	o.QrPayMode = req.QrPayMode
	o.QrCodeWidth = req.QrCodeWidth
	o.ReturnUrl = req.ReturnUrl
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
		TradeProvider: raw.TradeProvider,
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
