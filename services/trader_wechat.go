package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/web3-identity/conflux-pay/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type WechatTrader struct {
}

// precreate
func (w *WechatTrader) PreCreate(appName string, tradeNo string, req MakeOrderReq) (*models.OrderCore, error) {
	app := config.MustGetApp(appName)

	expire := time.Unix(req.TimeExpire, 0)
	descr := req.Description
	if descr != nil {
		tmp := utils.ReplaceEmoji(*req.Description, "[e]")
		descr = &tmp
	}

	prepayReq := PrepayRequest{
		Appid:       &app.AppId,
		Mchid:       &config.CompanyVal.MchID,
		Description: descr,
		OutTradeNo:  &tradeNo,
		TimeExpire:  &expire,
		Amount:      &native.Amount{Total: &req.Amount},
		NotifyUrl:   config.GetWxPayNotifyUrl(tradeNo),
	}

	orderCore := &models.OrderCore{
		TradeType: req.TradeType,
		TradeNo:   tradeNo,
	}

	switch req.TradeType {
	case enums.TRADE_TYPE_NATIVE:
		resp, _, err := wxNativeService.Prepay(context.Background(), native.PrepayRequest(prepayReq))
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("native pay complete")
		if err != nil {
			return nil, err
		}
		orderCore.CodeUrl = resp.CodeUrl
	case enums.TRADE_TYPE_H5:
		resp, _, err := wxH5Service.Prepay(context.Background(), *prepayReq.toH5())
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("h5 pay complete")
		if err != nil {
			return nil, err
		}
		orderCore.H5Url = resp.H5Url
	default:
		return nil, enums.ErrUnkownTradeType
	}
	return orderCore, nil
}

// get order
func (w *WechatTrader) GetTradeState(tradeNo string) (enums.TradeState, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return enums.TRADE_STATE_NIL, err
	}

	var tx *payments.Transaction

	switch o.TradeType {
	case enums.TRADE_TYPE_NATIVE:
		resp, _, err := wxNativeService.QueryOrderByOutTradeNo(context.Background(), native.QueryOrderByOutTradeNoRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &tradeNo,
		})
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
		if err != nil {
			return enums.TRADE_STATE_NIL, err
		}
		tx = resp
	case enums.TRADE_TYPE_H5:
		resp, _, err := wxH5Service.QueryOrderByOutTradeNo(context.Background(), h5.QueryOrderByOutTradeNoRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &tradeNo,
		})
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
		if err != nil {
			return enums.TRADE_STATE_NIL, err
		}
		tx = resp

	default:
		return enums.TRADE_STATE_NIL, enums.ErrUnkownTradeType
	}

	v, ok := enums.ParseTradeState(*tx.TradeState)
	if !ok {
		return enums.TRADE_STATE_NIL, fmt.Errorf("failed to parse trade state %v", *tx.TradeState)
	}

	return *v, nil
}

// refund
func (w *WechatTrader) Refund(tradeNo string, req RefundReq) error {
	oSummary, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return err
	}

	oSummary.AppPayNotifyUrl = req.NotifyUrl
	if err = oSummary.Save(); err != nil {
		return err
	}

	order, err := models.FindWechatOrderDetailByTradeNo(tradeNo)
	if err != nil {
		return err
	}

	if order.Amount == 0 {
		return fmt.Errorf("nothing could be refund")
	}

	_, _, err = wxRefundService.Create(context.Background(),
		refunddomestic.CreateRequest{
			OutTradeNo:  order.TradeNo,
			OutRefundNo: order.TradeNo,
			Reason:      &req.Reason,
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(int64(order.Amount)),
				Total:    core.Int64(int64(order.Amount)),
			},
			NotifyUrl: config.GetWxRefundNotifyUrl(tradeNo),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// get refund
func (w *WechatTrader) GetRefundState(tradeNo string) (enums.RefundState, error) {
	req := refunddomestic.QueryByOutRefundNoRequest{OutRefundNo: &tradeNo}
	resp, _, err := wxRefundService.QueryByOutRefundNo(context.Background(), req)
	logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order refund detail complete remote")
	if err != nil {
		return enums.REFUND_STATE_NIL, err
	}

	v, ok := enums.ParserefundState(string(*resp.Status))
	if !ok {
		return enums.REFUND_STATE_NIL, fmt.Errorf("failed to parse refund state %v", *resp.Status)
	}

	return *v, nil
}

// close
func (w *WechatTrader) Close(tradeNo string) error {
	order, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return err
	}

	switch order.TradeType {
	case enums.TRADE_TYPE_NATIVE:
		result, err := wxNativeService.CloseOrder(context.Background(), native.CloseOrderRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &order.TradeNo,
		})
		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
		if err != nil {
			return err
		}

	case enums.TRADE_TYPE_H5:
		result, err := wxH5Service.CloseOrder(context.Background(), h5.CloseOrderRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &order.TradeNo,
		})
		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
		if err != nil {
			return err
		}
	default:
		return enums.ErrUnkownTradeType
	}
	return nil
}
