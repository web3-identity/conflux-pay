package enums

func GetUnpayTradeStates() []TradeState {
	return []TradeState{
		TRADE_STATE_NIL,
		TRADE_STATE_NOTPAY,
		TRADE_STATE_USERPAYING,
	}
}

func GetUncompleteRefundStates() []RefundState {
	return []RefundState{
		REFUND_STATE_ABNORMAL,
		REFUND_STATE_CLOSED,
		REFUND_STATE_NIL,
	}
}
