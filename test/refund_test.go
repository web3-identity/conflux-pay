package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/logger"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/web3-identity/conflux-pay/services"
)

func init() {
	config.Init()
	logger.Init()
	services.Init()
	models.ConnectDB()
}

func TestTradeStateAndRefundStateRight(t *testing.T) {
	var o *models.Order
	err := models.GetDB().Model(&models.Order{}).Where("id=?", 113).First(&o).Error
	assert.NoError(t, err)

	oCore, err := services.NewOrderService().Refund(o.TradeNo, services.RefundReq{Reason: "test"})
	assert.NoError(t, err)

	assert.Equal(t, oCore.TradeState, enums.TRADE_STATE_REFUND)
}

func TestAlpayRefund(t *testing.T) {
	// 已经退过款的正常返回
	oCore, err := services.NewOrderService().Refund("AL167662031608000001", services.RefundReq{Reason: "test"})
	assert.NoError(t, err)
	assert.Equal(t, oCore.TradeState, enums.TRADE_STATE_REFUND)

	// 不存在的订单应该返回错误
	_, err = services.NewOrderService().Refund("AL167757297410200001", services.RefundReq{Reason: "test"})
	assert.Error(t, err)
}
