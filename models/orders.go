package models

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/web3-identity/conflux-pay/models/enums"
)

type OrderCore struct {
	AppName       string              `gorm:"type:varchar(32)" json:"app_name"`
	TradeProvider enums.TradeProvider `gorm:"uint" json:"trade_provider" swaggertype:"string"`
	TradeNo       string              `gorm:"type:varchar(32);uniqueIndex" json:"trade_no"`
	TradeType     enums.TradeType     `gorm:"uint" json:"trade_type" swaggertype:"string"`
	TradeState    enums.TradeState    `gorm:"uint" json:"trade_state" swaggertype:"string"`
	RefundState   enums.RefundState   `gorm:"uint" json:"refund_state" swaggertype:"string"`
	Amount        uint                `gorm:"uint" json:"amount"` // 单位为分
	Description   *string             `gorm:"type:varchar(255)" json:"description"`
	TimeExpire    *time.Time          `json:"time_expire,omitempty"`
	QrPayMode     string              `json:"qr_pay_mode,omitempty"`   // 支付二维码模式。 只有alipay，且 trade type 为 h5 模式有效; 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	QrCodeWidth   string              `json:"qr_code_width,omitempty"` // 二维码宽度。 只有alipay，且 trade type 为 h5 模式有效，qr pay mode 为4 时有效； 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	ReturnUrl     string              `json:"return_url,omitempty"`    // 付款成功后的跳转链接。只有alipay，且 trade type 为 h5 模式有效; 用法参考 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay?scene=22
	CodeUrl       *string             `gorm:"type:varchar(255)" json:"code_url,omitempty"`
	H5Url         *string             `gorm:"type:varchar(1023)" json:"h5_url,omitempty"`
	WapUrl        *string             `gorm:"type:varchar(1023)" json:"wap_url,omitempty"`
}

// 付款交易和退款交易都完成
func (o *OrderCore) IsStable() bool {
	return o.TradeState.IsStable() && o.RefundState.IsStable(o.TradeState)
}

type OrderNofity struct {
	AppPayNotifyUrl *string `gorm:"type:varchar(255)" json:"app_pay_notify_url"` // 上层应用通知url
	// PayNotifyNextTime    *time.Time `json:"pay_notify_next_time"`
	PayNotifyCount       int  `json:"pay_notify_count"`
	IsPayNotifyCompleted bool `json:"is_pay_notify_completed"`

	AppRefundNotifyUrl *string `gorm:"type:varchar(255)" json:"app_refund_notify_url"` // 上层应用通知url
	// RefundNotifyNextTime    *time.Time `json:"refund_notify_next_time"`
	RefundNotifyCount       int  `json:"refund_notify_count"`
	IsRefundNotifyCompleted bool `json:"is_refund_notify_completed"`
}

type Order struct {
	BaseModel
	OrderCore
	OrderNofity
}

func (o *Order) Save() error {
	err := GetDB().Save(o).Error
	if err != nil {
		logrus.WithError(err).Error("failed save order")
	}
	return err
}

func (o *Order) UpdateStates(tradeState enums.TradeState, refundState enums.RefundState) {
	if tradeState == o.TradeState && refundState == o.RefundState {
		return
	}

	o.TradeState = tradeState
	o.RefundState = refundState

	GetDB().Save(o)
}

func (o *Order) UpdateTradeState(tradeState enums.TradeState) error {
	if tradeState == o.TradeState {
		return nil
	}

	o.TradeState = tradeState
	return GetDB().Save(o).Error
}

func (o *Order) UpdateRefundState(refundState enums.RefundState) error {
	if refundState == o.RefundState {
		return nil
	}

	o.RefundState = refundState
	return GetDB().Save(o).Error
}

func FindOrderByTradeNo(tradeNo string) (*Order, error) {
	o := Order{}
	o.TradeNo = tradeNo
	return &o, GetDB().Where(&o).First(&o).Error
}

func FindNeedNotifyOrders(startId uint) ([]*Order, error) {
	var orders []*Order
	if err := GetDB().Where("id > ?", startId).
		Where("is_pay_notify_completed = ? or is_refund_notify_completed = ?", false, false).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func FindNeedCloseOrders() ([]*Order, error) {
	var orders []*Order
	if err := GetDB().Debug().
		Where("time_expire < ? and trade_state in ? and trim(app_name)!=''", time.Now(), enums.GetUnpayTradeStates()).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
