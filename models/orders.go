package models

import (
	"time"

	"github.com/wangdayong228/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type Order struct {
	BaseModel
	OrderCore
}

type OrderCore struct {
	AppName     string              `gorm:"type:varchar(32)"`
	Provider    enums.TradeProvider `gorm:"uint" json:"trade_provider" swaggertype:"string"`
	TradeNo     string              `gorm:"type:varchar(32);uniqueIndex" json:"trade_no"`
	TradeType   enums.TradeType     `gorm:"uint" json:"trade_type" swaggertype:"string"`
	TradeState  enums.TradeState    `gorm:"uint" json:"trade_state" swaggertype:"string"`
	RefundState enums.RefundState   `gorm:"uint" json:"refund_state" swaggertype:"string"`
	Amount      uint                `gorm:"uint" json:"amount"` // 单位为分
	Description *string             `gorm:"type:varchar(255)" json:"description"`
	TimeExpire  *time.Time          `json:"time_expire,omitempty"`
	CodeUrl     *string             `gorm:"type:varchar(255)" json:"code_url,omitempty"`
	H5Url       *string             `gorm:"type:varchar(255)" json:"h5_url,omitempty"`
}

func (o *OrderCore) IsStable() bool {
	return o.TradeState.IsStable() && o.RefundState.IsStable(o.TradeState)
}

func FindOrderByTradeNo(tradeNo string) (*Order, error) {
	o := Order{}
	o.TradeNo = tradeNo
	return &o, GetDB().Where(&o).First(&o).Error
}

type WechatOrderDetail struct {
	BaseModel
	Amount         uint    `gorm:"type:varchar(32);" json:"amount,omitempty"`
	Appid          *string `gorm:"type:varchar(32);" json:"appid,omitempty"`
	Attach         *string `gorm:"type:varchar(32);" json:"attach,omitempty"`
	BankType       *string `gorm:"type:varchar(32);" json:"bank_type,omitempty"`
	Mchid          *string `gorm:"type:varchar(32);" json:"mchid,omitempty"`
	TradeNo        *string `gorm:"type:varchar(32);uniqueIndex" json:"trade_no,omitempty"`
	Payer          *string `gorm:"type:varchar(32);" json:"payer,omitempty"`
	SuccessTime    *string `gorm:"type:varchar(32);" json:"success_time,omitempty"`
	TradeState     *string `gorm:"type:varchar(32);" json:"trade_state,omitempty"`
	TradeStateDesc *string `gorm:"type:varchar(255);" json:"trade_state_desc,omitempty"`
	TradeType      *string `gorm:"type:varchar(32);" json:"trade_type,omitempty"`
	TransactionId  *string `gorm:"type:varchar(32);" json:"transaction_id,omitempty"`
	RefundNo       *string `gorm:"type:varchar(32);" json:"refund_no,omitempty"`
	RefundStatus   *string `gorm:"type:varchar(32);" json:"refresh_status,omitempty"`
	// PromotionDetail []PromotionDetail `json:"promotion_detail,omitempty"`
}

func FindWechatOrderDetailByTradeNo(tradeNo string) (*WechatOrderDetail, error) {
	o := WechatOrderDetail{
		TradeNo: &tradeNo,
	}
	return &o, GetDB().Where(&o).First(&o).Error
}

func NewWechatOrderDetailByRaw(raw *payments.Transaction) *WechatOrderDetail {
	payer := (*string)(nil)
	if raw.Payer != nil {
		payer = raw.Payer.Openid
	}

	amount := uint(0)
	if raw != nil && raw.Amount != nil && raw.Amount.Total != nil {
		amount = uint(*raw.Amount.Total)
	}

	return &WechatOrderDetail{
		Amount:         amount,
		Appid:          raw.Appid,
		Attach:         raw.Attach,
		BankType:       raw.BankType,
		Mchid:          raw.Mchid,
		TradeNo:        raw.OutTradeNo,
		Payer:          payer,
		SuccessTime:    raw.SuccessTime,
		TradeState:     raw.TradeState,
		TradeStateDesc: raw.TradeStateDesc,
		TradeType:      raw.TradeType,
		TransactionId:  raw.TransactionId,
	}
}

func UpdateWechatOrderDetail(val *WechatOrderDetail) error {
	valInDb, err := FindWechatOrderDetailByTradeNo(*val.TradeNo)
	if err != nil {
		return err
	}

	val.BaseModel = valInDb.BaseModel
	return GetDB().Save(val).Error
}

type AlipayOrderDetail struct {
}
