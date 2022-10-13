package models

type TradeType uint

const (
	TRADE_TYPE_NATIVE = iota + 1
	TRADE_TYPE_H5
)

type TradeState uint

const (
	TRADE_STATE_SUCCESSS = iota + 1
	TRADE_STATE_REFUND
	TRADE_STATE_NOTPAY
	TRADE_STATE_CLOSED
	TRADE_STATE_REVOKED
	TRADE_STATE_USERPAYING
	TRADE_STATE_PAYERROR
)

type Order struct {
	BaseModel
	Commithash string     `gorm:"type:varchar(256)" json:"commit_hash"`
	PayChannel string     `gorm:"type:varchar(256)" json:"pay_channel"`
	TradeType  TradeType  `gorm:"uint" json:"trade_type"`
	TradeState TradeState `gorm:"uint" json:"trade_state"`
	Amount     uint       `gorm:"uint" json:"amount"` // 单位为分
}
