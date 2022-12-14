package enums

import (
	"encoding/json"
	"errors"
)

type TradeType uint

const (
	TRADE_TYPE_NATIVE = iota + 1
	TRADE_TYPE_H5
)

var (
	tradeTypeValue2StrMap map[TradeType]string
	tradeTypeStr2ValueMap map[string]TradeType
)

var (
	ErrUnkownTradeType = errors.New("unknown trade type")
)

func init() {
	tradeTypeValue2StrMap = map[TradeType]string{
		TRADE_TYPE_NATIVE: "native",
		TRADE_TYPE_H5:     "h5",
	}

	tradeTypeStr2ValueMap = make(map[string]TradeType)
	for k, v := range tradeTypeValue2StrMap {
		tradeTypeStr2ValueMap[v] = k
	}
}

func (t TradeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TradeType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, ok := ParseTradeType(str)
	if !ok {
		return errors.New("unkown trade_type")
	}
	*t = *val

	return nil
}

func (t *TradeType) String() string {
	v, ok := tradeTypeValue2StrMap[*t]
	if ok {
		return v
	}
	return "UNKNOWN"
}

func ParseTradeType(str string) (*TradeType, bool) {
	v, ok := tradeTypeStr2ValueMap[str]
	return &v, ok
}
