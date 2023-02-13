package enums

import (
	"encoding/json"
	"errors"
)

type TradeProvider uint

const (
	TRADE_PROVIDER_WECHAT TradeProvider = iota + 1
	TRADE_PROVIDER_ALIPAY
)

type TradeProviderDesc struct {
	Name string
	Code string
}

var (
	tradeProviderValue2StrMap  map[TradeProvider]TradeProviderDesc
	tradeProviderName2ValueMap map[string]TradeProvider
	tradeProviderCode2ValueMap map[string]TradeProvider
)

func init() {
	tradeProviderValue2StrMap = map[TradeProvider]TradeProviderDesc{
		TRADE_PROVIDER_WECHAT: {"wechat", "WX"},
		TRADE_PROVIDER_ALIPAY: {"alipay", "AL"},
	}

	tradeProviderName2ValueMap = make(map[string]TradeProvider)
	tradeProviderCode2ValueMap = make(map[string]TradeProvider)
	for k, v := range tradeProviderValue2StrMap {
		tradeProviderName2ValueMap[v.Name] = k
		tradeProviderCode2ValueMap[v.Code] = k
	}
}

func (t TradeProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TradeProvider) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	val, ok := ParseTradeProviderByName(str)
	if !ok {
		return errors.New("unkown trade_provider")
	}
	*t = *val

	return nil
}

func (p TradeProvider) String() string {
	v, ok := tradeProviderValue2StrMap[p]
	if ok {
		return v.Name
	}
	return "unkown"
}

func (p TradeProvider) Code() string {
	v, ok := tradeProviderValue2StrMap[p]
	if ok {
		return v.Code
	}
	return "UNKNOWN"
}

func ParseTradeProviderByName(str string) (*TradeProvider, bool) {
	v, ok := tradeProviderName2ValueMap[str]
	return &v, ok
}

func ParseTradeProviderByCode(str string) (*TradeProvider, bool) {
	v, ok := tradeProviderCode2ValueMap[str]
	return &v, ok
}
