package enums

import (
	"encoding/json"
	"errors"
)

type TradeState uint

// SUCCESS：支付成功
// REFUND：转入退款
// NOTPAY：未支付
// CLOSED：已关闭
// REVOKED：已撤销（仅付款码支付会返回）
// USERPAYING：用户支付中（仅付款码支付会返回）
// PAYERROR：支付失败（仅付款码支付会返回）
const (
	TRADE_STATE_NIL TradeState = iota
	TRADE_STATE_SUCCESSS
	TRADE_STATE_REFUND
	TRADE_STATE_NOTPAY
	TRADE_STATE_CLOSED
	TRADE_STATE_REVOKED
	TRADE_STATE_USERPAYING
	TRADE_STATE_PAYERROR
)

var (
	tradeStateValue2CodeMap map[TradeState]string
	tradeStateCode2ValueMap map[string]TradeState
)

var (
	ErrUnkownTradeState = errors.New("unknown trade state")
)

func init() {
	tradeStateValue2CodeMap = map[TradeState]string{
		TRADE_STATE_SUCCESSS:   "SUCCESS",
		TRADE_STATE_REFUND:     "REFUND",
		TRADE_STATE_NOTPAY:     "NOTPAY",
		TRADE_STATE_CLOSED:     "CLOSED",
		TRADE_STATE_REVOKED:    "REVOKED",
		TRADE_STATE_USERPAYING: "USERPAYING",
		TRADE_STATE_PAYERROR:   "PAYERROR",
	}

	tradeStateCode2ValueMap = make(map[string]TradeState)
	for k, v := range tradeStateValue2CodeMap {
		tradeStateCode2ValueMap[v] = k
	}
}

func (t TradeState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TradeState) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, ok := ParseTradeState(str)
	if !ok {
		return errors.New("unkown trade_state")
	}
	*t = *val

	return nil
}

func (t *TradeState) String() string {
	v, ok := tradeStateValue2CodeMap[*t]
	if ok {
		return v
	}
	return "UNKNOWN"
}

func ParseTradeState(code string) (*TradeState, bool) {
	v, ok := tradeStateCode2ValueMap[code]
	return &v, ok
}

func (t TradeState) IsStable() bool {
	return t != TRADE_STATE_NOTPAY && t != TRADE_STATE_USERPAYING && t != TRADE_STATE_SUCCESSS
}

func (t TradeState) IsSuccess() bool {
	return t == TRADE_STATE_SUCCESSS
}
