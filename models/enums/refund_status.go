package enums

import (
	"encoding/json"
	"errors"
)

type RefundState uint

// SUCCESS ：退款成功
// CLOSED ：退款关闭
// PROCESSING ： 退款处理中
// ABNORMAL ：退款异常
const (
	REFUND_STATE_NIL = iota
	REFUND_STATE_SUCCESSS
	REFUND_STATE_CLOSED
	REFUND_STATE_PROCESSING
	REFUND_STATE_ABNORMAL
)

var (
	refundStateValue2CodeMap map[RefundState]string
	refundStateCode2ValueMap map[string]RefundState
)

func init() {
	refundStateValue2CodeMap = map[RefundState]string{
		REFUND_STATE_NIL:        "NIL",
		REFUND_STATE_SUCCESSS:   "SUCCESS",
		REFUND_STATE_CLOSED:     "CLOSED",
		REFUND_STATE_PROCESSING: "PROCESSING",
		REFUND_STATE_ABNORMAL:   "ABNORMAL",
	}

	refundStateCode2ValueMap = make(map[string]RefundState)
	for k, v := range refundStateValue2CodeMap {
		refundStateCode2ValueMap[v] = k
	}
}

func (t RefundState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *RefundState) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, ok := ParserefundState(str)
	if !ok {
		return errors.New("unkown refund_state")
	}
	*t = *val

	return nil
}

func (t *RefundState) String() string {
	v, ok := refundStateValue2CodeMap[*t]
	if ok {
		return v
	}
	return "UNKNOWN"
}

func ParserefundState(code string) (*RefundState, bool) {
	v, ok := refundStateCode2ValueMap[code]
	return &v, ok
}

func (t RefundState) IsStable(tradeState TradeState) bool {
	if !tradeState.IsStable() {
		return false
	}
	if tradeState != TRADE_STATE_REFUND {
		return true
	}
	return t != REFUND_STATE_NIL && t != REFUND_STATE_PROCESSING
}
