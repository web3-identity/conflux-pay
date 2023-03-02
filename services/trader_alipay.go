package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay/v3"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/web3-identity/conflux-pay/utils"
)

type AlipayTrader struct {
	app    config.App
	client *alipay.Client
}

func NewAlipayTrader(appName string) (*AlipayTrader, error) {
	app := config.MustGetApp(appName)
	client, err := alipay.New(app.AppIdAlipay, config.CompanyVal.Alipay.PrivateKey, true)
	if err != nil {
		return nil, err
	}
	err = client.LoadAliPayPublicKey(config.CompanyVal.Alipay.AlipayPublicKey)
	if err != nil {
		return nil, err
	}

	trader := &AlipayTrader{
		app:    app,
		client: client,
	}

	return trader, nil
}

// precreate
func (a *AlipayTrader) PreCreate(tradeNo string, req MakeOrderReq) (*models.OrderCore, error) {

	expire := time.Unix(req.TimeExpire, 0)
	descr := req.Description
	if descr != nil {
		tmp := utils.ReplaceEmoji(*req.Description, "[e]")
		descr = &tmp
	}

	var trade = alipay.Trade{}
	trade.Subject = *req.Description
	trade.OutTradeNo = tradeNo
	trade.TotalAmount = decimal.NewFromInt(req.Amount).Div(decimal.NewFromInt(100)).String()

	orderCore := &models.OrderCore{
		TradeType: req.TradeType,
		TradeNo:   tradeNo,
	}

	switch req.TradeType {
	case enums.TRADE_TYPE_NATIVE:
		resp, err := a.client.TradePreCreate(alipay.TradePreCreate{Trade: trade})
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("native pay complete")
		if err != nil {
			return nil, err
		}
		if !resp.IsSuccess() {
			return nil, fmt.Errorf("msg:%v, sub msg:%v", resp.Content.Msg, resp.Content.SubMsg)
		}
		orderCore.CodeUrl = &resp.Content.QRCode
	case enums.TRADE_TYPE_H5:
		p := alipay.TradePagePay{Trade: trade}
		p.ProductCode = "FAST_INSTANT_TRADE_PAY"
		p.QRPayMode = req.QrPayMode
		p.QRCodeWidth = req.QrCodeWidth
		p.ReturnURL = req.ReturnUrl
		p.TimeoutExpress = fmt.Sprintf("%vm", math.Round(time.Until(expire).Minutes()))
		if time.Until(expire) < time.Minute || time.Until(expire) > time.Hour*24*15 {
			return nil, errors.New("expire time must between 1 minute and 24 hours")
		}

		resp, err := a.client.TradePagePay(p)
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("h5 pay complete")
		if err != nil {
			return nil, err
		}
		tmp := resp.String()
		orderCore.H5Url = &tmp
	case enums.TRADE_TYPE_WAP:
		p := alipay.TradeWapPay{Trade: trade}
		// p.ProductCode = "FAST_INSTANT_TRADE_PAY"
		// p.QRPayMode = req.QrPayMode
		// p.QRCodeWidth = req.QrCodeWidth
		// p.ReturnURL = req.ReturnUrl
		p.TimeoutExpress = fmt.Sprintf("%vm", math.Round(time.Until(expire).Minutes()))
		if time.Until(expire) < time.Minute || time.Until(expire) > time.Hour*24*15 {
			return nil, errors.New("expire time must between 1 minute and 24 hours")
		}

		resp, err := a.client.TradeWapPay(p)
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("wap pay complete")
		if err != nil {
			return nil, err
		}
		tmp := resp.String()
		orderCore.WapUrl = &tmp

	default:
		return nil, enums.ErrUnkownTradeType
	}

	a.GetTradeState(tradeNo)

	return orderCore, nil
}

// get trade state
func (a *AlipayTrader) GetTradeState(tradeNo string) (enums.TradeState, error) {
	// tradeNo = "20230320063232117"
	p := alipay.TradeQuery{OutTradeNo: tradeNo}
	res, err := a.client.TradeQuery(p)
	if err != nil {
		return enums.TRADE_STATE_NIL, err
	}
	if !res.IsSuccess() {
		return enums.TRADE_STATE_NIL, fmt.Errorf("msg:%v, sub msg:%v", res.Content.Msg, res.Content.SubMsg)
	}

	refundState, err := a.GetRefundState(tradeNo)
	if err != nil {
		return enums.TRADE_STATE_NIL, err
	}

	return convertAlTradeState(res.Content.TradeStatus, refundState != enums.REFUND_STATE_NIL), nil
}

// refund
func (a *AlipayTrader) Refund(tradeNo string, req RefundReq) error {
	tq := alipay.TradeQuery{OutTradeNo: tradeNo}
	orderRes, err := a.client.TradeQuery(tq)
	if err != nil {
		return err
	}
	if !orderRes.IsSuccess() {
		return fmt.Errorf("failed to query trade. msg:%v, sub msg:%v", orderRes.Content.Msg, orderRes.Content.SubMsg)
	}

	var tr = alipay.TradeRefund{}
	tr.OutTradeNo = tradeNo
	tr.RefundAmount = orderRes.Content.TotalAmount
	tr.RefundReason = req.Reason

	refundRes, err := a.client.TradeRefund(tr)
	if err != nil {
		return err
	}
	if !orderRes.IsSuccess() {
		return fmt.Errorf("failed to refund. msg:%v, sub msg:%v", refundRes.Content.Msg, refundRes.Content.SubMsg)
	}

	return nil
}

// get refund state
func (a *AlipayTrader) GetRefundState(tradeNo string) (enums.RefundState, error) {
	var p = alipay.TradeFastPayRefundQuery{}
	p.OutTradeNo = tradeNo
	p.OutRequestNo = p.OutTradeNo

	res, err := a.client.TradeFastPayRefundQuery(p)
	if err != nil {
		return enums.REFUND_STATE_NIL, err
	}
	if !res.IsSuccess() {
		// 外部订单号不存在
		// if strings.ToUpper(res.Content.SubCode) == "ACQ.TRADE_NOT_EXIST" {
		// 	return enums.REFUND_STATE_NIL, nil
		// }
		// return enums.REFUND_STATE_NIL, fmt.Errorf("msg:%v, sub msg:%v, sub code:%v", res.Content.Msg, res.Content.SubMsg, res.Content.SubCode)
		return enums.REFUND_STATE_NIL, nil
	}

	return convertAlRefundState(res.Content.RefundStatus), nil
}

// close
func (a *AlipayTrader) Close(tradeNo string) error {
	var p = alipay.TradeClose{}
	p.OutTradeNo = tradeNo

	res, err := a.client.TradeClose(p)
	if err != nil {
		return err
	}
	if !IsAlResSuccess(res.Content.Code) {
		return fmt.Errorf("msg:%v, sub msg:%v", res.Content.Msg, res.Content.SubMsg)
	}
	return nil
}

func convertAlTradeState(status alipay.TradeStatus, isRefund bool) enums.TradeState {
	if isRefund {
		return enums.TRADE_STATE_REFUND
	}

	switch status {
	case alipay.TradeStatusWaitBuyerPay:
		return enums.TRADE_STATE_NOTPAY
	case alipay.TradeStatusClosed:
		return enums.TRADE_STATE_CLOSED
		// TODO: 交易结束不可退款，微信没有对应状态，暂用SUCCESS
	case alipay.TradeStatusFinished:
		fallthrough
	case alipay.TradeStatusSuccess:
		return enums.TRADE_STATE_SUCCESSS
	}
	return enums.TRADE_STATE_NIL
}

func convertAlRefundState(refundState string) enums.RefundState {
	if refundState == "REFUND_SUCCESS" {
		return enums.REFUND_STATE_SUCCESSS
	}
	return enums.REFUND_STATE_PROCESSING
}

func IsAlResSuccess(alipayCode alipay.Code) bool {
	return alipayCode == alipay.CodeSuccess
}

func IsAlNotExistErr(err error) bool {
	return strings.Contains(err.Error(), "交易不存在")
}
