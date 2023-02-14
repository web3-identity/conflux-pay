package services

import "testing"

func TestWechatTraderIsTrader(t *testing.T) {
	var _ Trader = &WechatTrader{}
}

func TestAlipayTraderIsTrader(t *testing.T) {
	var _ Trader = &AlipayTrader{}
}
