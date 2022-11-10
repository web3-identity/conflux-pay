package models

import "github.com/wechatpay-apiv3/wechatpay-go/services/payments"

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
