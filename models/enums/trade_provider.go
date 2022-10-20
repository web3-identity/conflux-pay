package enums

type TradeProvider uint

const (
	TRADE_PROVIDER_WECHAT = iota + 1
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

func (p *TradeProvider) String() string {
	v, ok := tradeProviderValue2StrMap[*p]
	if ok {
		return v.Name
	}
	return "unknown"
}

func (p *TradeProvider) Code() string {
	v, ok := tradeProviderValue2StrMap[*p]
	if ok {
		return v.Code
	}
	return "unknown"
}

func ParseTradeProviderByName(str string) (*TradeProvider, bool) {
	v, ok := tradeProviderName2ValueMap[str]
	return &v, ok
}

func ParseTradeProviderByCode(str string) (*TradeProvider, bool) {
	v, ok := tradeProviderCode2ValueMap[str]
	return &v, ok
}
