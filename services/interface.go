package services

import (
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
)

type Trader interface {
	// precreate
	PreCreate(tradeNo string, req MakeOrderReq) (*models.OrderCore, error)
	// get trade state
	GetTradeState(tradeNo string) (enums.TradeState, error)
	// refund
	Refund(tradeNo string, req RefundReq) error
	// get refund state
	GetRefundState(tradeNo string) (enums.RefundState, error)
	// close
	Close(tradeNo string) error
}
