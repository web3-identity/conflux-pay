package enums

type TradeProvider uint

const (
	TRADE_PROVIDER_WECHAT = iota + 1
	TRADE_PROVIDER_ALIPAY
)

func (p *TradeProvider) String() string {
	switch *p {
	case TRADE_PROVIDER_WECHAT:
		return "WECHAT"
	case TRADE_PROVIDER_ALIPAY:
		return "APLIPAY"
	}
	return "UNKNOWN"
}

func (p *TradeProvider) Code() string {
	switch *p {
	case TRADE_PROVIDER_WECHAT:
		return "WX"
	case TRADE_PROVIDER_ALIPAY:
		return "AL"
	}
	return "UNKNOWN"
}
