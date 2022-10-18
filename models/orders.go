package models

import (
	"time"

	"github.com/wangdayong228/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type Order struct {
	BaseModel
	Provider    enums.TradeProvider `gorm:"uint" json:"trade_provider"`
	TradeNo     string              `gorm:"type:varchar(32);uniqueIndex" json:"trade_no"`
	TradeType   enums.TradeType     `gorm:"uint" json:"trade_type"`
	TradeState  enums.TradeState    `gorm:"uint" json:"trade_state"`
	Amount      uint                `gorm:"uint" json:"amount"` // 单位为分
	Description *string             `gorm:"type:varchar(256)" json:"description"`
	TimeExpire  *time.Time          `json:"time_expire,omitempty"`
	CodeUrl     *string             `gorm:"type:varchar(256)" json:"code_url,omitempty"`
	H5Url       *string             `gorm:"type:varchar(256)" json:"h5_url,omitempty"`
}

func FindOrderByTradeNo(tradeNo string) (*Order, error) {
	o := Order{
		TradeNo: tradeNo,
	}
	return &o, GetDB().First(&o).Error
}

type WechatOrderDetail struct {
	Amount         uint    `gorm:"type:varchar(32);" json:"amount,omitempty"`
	Appid          *string `gorm:"type:varchar(32);" json:"appid,omitempty"`
	Attach         *string `gorm:"type:varchar(32);" json:"attach,omitempty"`
	BankType       *string `gorm:"type:varchar(32);" json:"bank_type,omitempty"`
	Mchid          *string `gorm:"type:varchar(32);" json:"mchid,omitempty"`
	TradeNo        *string `gorm:"type:varchar(32);" json:"trade_no,omitempty"`
	Payer          *string `gorm:"type:varchar(32);" json:"payer,omitempty"`
	SuccessTime    *string `gorm:"type:varchar(32);" json:"success_time,omitempty"`
	TradeState     *string `gorm:"type:varchar(32);" json:"trade_state,omitempty"`
	TradeStateDesc *string `gorm:"type:varchar(256);" json:"trade_state_desc,omitempty"`
	TradeType      *string `gorm:"type:varchar(32);" json:"trade_type,omitempty"`
	TransactionId  *string `gorm:"type:varchar(32);" json:"transaction_id,omitempty"`
	// PromotionDetail []PromotionDetail `json:"promotion_detail,omitempty"`
}

func FindWechatOrderDetailByTradeNo(tradeNo string) (*WechatOrderDetail, error) {
	o := WechatOrderDetail{
		TradeNo: &tradeNo,
	}
	return &o, GetDB().First(&o).Error
}

func NewWechatOrderDetailByRaw(raw *payments.Transaction) *WechatOrderDetail {
	return &WechatOrderDetail{
		Amount:         uint(*raw.Amount.Total),
		Appid:          raw.Appid,
		Attach:         raw.Attach,
		BankType:       raw.BankType,
		Mchid:          raw.Mchid,
		TradeNo:        raw.OutTradeNo,
		Payer:          raw.Payer.Openid,
		SuccessTime:    raw.SuccessTime,
		TradeState:     raw.TradeState,
		TradeStateDesc: raw.TradeStateDesc,
		TradeType:      raw.TradeType,
		TransactionId:  raw.TransactionId,
	}
}

type AlipayOrderDetail struct {
}
