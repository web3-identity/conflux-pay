package services

import (
	"context"

	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	cns_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type OrderService struct {
	TradeStateChangeEvent  []TradeStateChangeHandler
	RefundStateChangeEvent []RefundStateChangeHandler
	wxTraders              map[string]*WechatTrader
	alTraders              map[string]*AlipayTrader
	traderLock             sync.Mutex
}

func NewOrderService() *OrderService {
	w := OrderService{
		wxTraders: make(map[string]*WechatTrader),
		alTraders: make(map[string]*AlipayTrader),
	}

	w.RegisterTradeStateChangeEvent(func(o *models.Order) {
		go runPayNotifyTask(o)
		go runRefundNotifyTask(o) // 微信支付只通知成功的状态，在交易关闭后也需要处理refund notify标志
	})
	w.RegisterRefundStateChangeEvent(func(o *models.Order) {
		go runRefundNotifyTask(o)
	})
	return &w
}

func (w *OrderService) GetTrader(appName string, provider enums.TradeProvider) (Trader, error) {
	app := config.MustGetApp(appName)
	appId := app.AppIdAlipay
	switch provider {
	case enums.TRADE_PROVIDER_ALIPAY:
		if _, ok := w.alTraders[appId]; !ok {
			w.traderLock.Lock()
			t, err := NewAlipayTrader(appName)
			if err != nil {
				w.traderLock.Unlock()
				return nil, err
			}
			w.alTraders[appName] = t
			w.traderLock.Unlock()
		}
		return w.alTraders[appName], nil
	case enums.TRADE_PROVIDER_WECHAT:
		if _, ok := w.wxTraders[appId]; !ok {
			w.traderLock.Lock()
			t := NewWechatTrader(appName)
			w.wxTraders[appName] = t
			w.traderLock.Unlock()
		}
		return w.wxTraders[appName], nil
	}
	return nil, errors.New("unkown provider")
}

func (w *OrderService) MustGetTrader(appName string, provider enums.TradeProvider) Trader {
	t, err := w.GetTrader(appName, provider)
	if err != nil {
		panic(err)
	}
	return t
}

func (w *OrderService) RegisterTradeStateChangeEvent(h TradeStateChangeHandler) {
	w.TradeStateChangeEvent = append(w.TradeStateChangeEvent, h)
	// w.TradeStateChangeEvent = append(w.TradeStateChangeEvent, w.GetOrderDetailAndSave)
}

func (w *OrderService) RegisterRefundStateChangeEvent(h RefundStateChangeHandler) {
	w.RefundStateChangeEvent = append(w.RefundStateChangeEvent, h)
}

func (w *OrderService) InvokeTradeStateChangedEvent(o *models.Order) {
	fmt.Println("invoke trade state changed")
	for _, h := range w.TradeStateChangeEvent {
		h(o)
	}
}

func (w *OrderService) InvokeRefundStateChangedEvent(o *models.Order) {
	fmt.Println("invoke refund state changed")
	for _, h := range w.RefundStateChangeEvent {
		h(o)
	}
}

// 统一wechat所有下单接口
// 生成 trade_no
// 返回 trade_no, pay_url

func (w *OrderService) MakeOrder(req MakeOrderReq) (*models.Order, error) {
	appName := req.AppName
	logrus.WithField("app name", appName).WithField("req", req).Info("make order")
	app := config.MustGetApp(appName)
	no := genTradeNo(app.AppInternalID, req.MustGetTradeProvider())
	orderResp, err := w.MustGetTrader(appName, req.MustGetTradeProvider()).PreCreate(no, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// save db
	order := &models.Order{
		OrderCore: models.OrderCore{
			AppName:    appName,
			TradeNo:    no,
			TradeState: enums.TRADE_STATE_NOTPAY,
			CodeUrl:    orderResp.CodeUrl,
			H5Url:      orderResp.H5Url,
			WapUrl:     orderResp.WapUrl,
		},
	}
	req.FillToOrder(order)

	if err = order.Save(); err != nil {
		return nil, err
	}
	go w.autoCloseOrder(order)

	return order, nil
}

func (w *OrderService) GetOrder(tradeNo string) (*models.Order, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if config.NotifyConfig[o.TradeProvider].Enable {
		return o, nil
	}

	if !o.IsStable() {
		tradeState, err := w.MustGetTrader(o.AppName, o.TradeProvider).GetTradeState(tradeNo)
		if err != nil {
			// 支付宝未付款的交易都会返回不存在
			if o.TradeProvider == enums.TRADE_PROVIDER_ALIPAY && IsNotExistErr(err) {
				return o, nil
			}
			return nil, errors.WithStack(err)
		}

		refundState := o.RefundState
		if tradeState == enums.TRADE_STATE_REFUND {
			_refundState, err := w.MustGetTrader(o.AppName, o.TradeProvider).GetRefundState(tradeNo)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			refundState = _refundState
		}

		o.UpdateStates(tradeState, refundState)
	}
	return o, nil
}

func (w *OrderService) getOrderCore(tradeNo string) (*models.OrderCore, error) {
	o, err := w.GetOrder(tradeNo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &o.OrderCore, nil
}

func (w *OrderService) RefreshUrl(tradeNo string) (*MakeOrderResp, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// return error if order is complete
	if o.TradeState.IsStable() {
		return nil, cns_errors.ERR_ORDER_COMPLETED
	}

	resp, err := w.MustGetTrader(o.AppName, o.TradeProvider).PreCreate(tradeNo, *NewMakeOrderReqFromOrder(o))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	o.CodeUrl = resp.CodeUrl
	o.H5Url = resp.H5Url
	o.Save()
	return NewMakeOrderRespFromRaw(o), nil
}

func (w *OrderService) Refund(tradeNo string, req RefundReq) (*models.OrderCore, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = w.MustGetTrader(o.AppName, o.TradeProvider).Refund(o.TradeNo, req); err != nil {
		return nil, errors.WithStack(err)
	}
	return w.getOrderCore(tradeNo)
}

func (w *OrderService) Close(tradeNo string) (*models.OrderCore, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = w.MustGetTrader(o.AppName, o.TradeProvider).Close(o.TradeNo); err != nil {
		return nil, errors.WithStack(err)
	}
	return w.getOrderCore(tradeNo)
}

// func (w *WechatOrderService) prePay(appName string, tradeNo string, req MakeOrderReq) (*MakeOrderResp, error) {
// 	app := config.MustGetApp(appName)

// 	expire := time.Unix(req.TimeExpire, 0)
// 	descr := req.Description
// 	if descr != nil {
// 		tmp := utils.ReplaceEmoji(*req.Description, "[e]")
// 		descr = &tmp
// 	}

// 	prepayReq := PrepayRequest{
// 		Appid:       &app.AppId,
// 		Mchid:       &config.CompanyVal.MchID,
// 		Description: descr,
// 		OutTradeNo:  &tradeNo,
// 		TimeExpire:  &expire,
// 		Amount:      &native.Amount{Total: &req.Amount},
// 		NotifyUrl:   config.GetWxPayNotifyUrl(tradeNo),
// 	}

// 	orderResp := &MakeOrderResp{
// 		TradeType: req.TradeType,
// 		TradeNo:   tradeNo,
// 	}

// 	switch req.TradeType {
// 	case enums.TRADE_TYPE_NATIVE:
// 		resp, _, err := wxNativeService.Prepay(context.Background(), native.PrepayRequest(prepayReq))
// 		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("native pay complete")
// 		if err != nil {
// 			return nil, err
// 		}
// 		orderResp.CodeUrl = resp.CodeUrl
// 	case enums.TRADE_TYPE_H5:
// 		resp, _, err := wxH5Service.Prepay(context.Background(), *prepayReq.toH5())
// 		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("h5 pay complete")
// 		if err != nil {
// 			return nil, err
// 		}
// 		orderResp.H5Url = resp.H5Url
// 	default:
// 		return nil, enums.ErrUnkownTradeType
// 	}
// 	return orderResp, nil
// }

// func (w *WechatOrderService) GetOrderDetail(tradeNo string) (*models.WechatOrderDetail, error) {
// 	if config.WechatOrderConfig.Enable {
// 		return models.FindWechatOrderDetailByTradeNo(tradeNo)
// 	}
// 	return w.GetOrderDetailAndSave(tradeNo)
// }

// func (w *WechatOrderService) GetOrderDetailAndSave(tradeNo string) (*models.WechatOrderDetail, error) {
// 	o, err := models.FindOrderByTradeNo(tradeNo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logrus.WithField("current order", o).Info("will get order detail and save")

// 	if o.IsStable() {
// 		return models.FindWechatOrderDetailByTradeNo(tradeNo)
// 	}

// 	detail, err := w.getRemoteOrderDetail(tradeNo, o.TradeType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	v, ok := enums.ParseTradeState(*detail.TradeState)
// 	if !ok {
// 		return nil, fmt.Errorf("unknown trade state %v", *detail.TradeState)
// 	}

// 	if *v != o.TradeState {
// 		o.TradeState = *v
// 		models.UpdateWechatOrderDetail(detail)
// 		models.GetDB().Save(o)
// 		logrus.WithField("trade_no", o.TradeNo).WithField("trade_state", o.TradeState).Info("update order and detail")
// 	}

// 	if o.TradeState != enums.TRADE_STATE_REFUND {
// 		return detail, nil
// 	}

// 	refundDetial, err := w.getRemoteRefundDetail(tradeNo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	models.UpdateRefundDetail(refundDetial)

// 	if v, ok := enums.ParserefundState(*refundDetial.Status); ok && *v != o.RefundState {
// 		o.RefundState = *v
// 		models.GetDB().Save(o)
// 	}

// 	return detail, nil
// }

// func (w *WechatOrderService) getRemoteOrderDetail(tradeNo string, tradeType enums.TradeType) (*models.WechatOrderDetail, error) {
// 	switch tradeType {
// 	case enums.TRADE_TYPE_NATIVE:
// 		resp, _, err := wxNativeService.QueryOrderByOutTradeNo(context.Background(), native.QueryOrderByOutTradeNoRequest{
// 			Mchid:      &config.CompanyVal.MchID,
// 			OutTradeNo: &tradeNo,
// 		})
// 		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
// 		if err != nil {
// 			return nil, err
// 		}
// 		return models.NewWechatOrderDetailByRaw(resp), nil
// 	case enums.TRADE_TYPE_H5:
// 		resp, _, err := wxH5Service.QueryOrderByOutTradeNo(context.Background(), h5.QueryOrderByOutTradeNoRequest{
// 			Mchid:      &config.CompanyVal.MchID,
// 			OutTradeNo: &tradeNo,
// 		})
// 		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
// 		if err != nil {
// 			return nil, err
// 		}
// 		return models.NewWechatOrderDetailByRaw(resp), nil

// 	default:
// 		return nil, enums.ErrUnkownTradeType
// 	}
// }

// func (w *WechatOrderService) getRemoteRefundDetail(tradeNo string) (*models.WechatRefundDetail, error) {
// 	req := refunddomestic.QueryByOutRefundNoRequest{OutRefundNo: &tradeNo}
// 	resp, _, err := wxRefundService.QueryByOutRefundNo(context.Background(), req)
// 	logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order refund detail complete remote")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return models.NewWechatRefundDetailByRaw(resp), nil
// }

// func (w *WechatOrderService) Close(tradeNo string) (*models.WechatOrderDetail, error) {
// 	order, err := models.FindOrderByTradeNo(tradeNo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch order.TradeType {
// 	case enums.TRADE_TYPE_NATIVE:
// 		result, err := wxNativeService.CloseOrder(context.Background(), native.CloseOrderRequest{
// 			Mchid:      &config.CompanyVal.MchID,
// 			OutTradeNo: &order.TradeNo,
// 		})
// 		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
// 		if err != nil {
// 			return nil, err
// 		}

// 	case enums.TRADE_TYPE_H5:
// 		result, err := wxH5Service.CloseOrder(context.Background(), h5.CloseOrderRequest{
// 			Mchid:      &config.CompanyVal.MchID,
// 			OutTradeNo: &order.TradeNo,
// 		})
// 		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		return nil, enums.ErrUnkownTradeType
// 	}
// 	return w.GetOrderDetailAndSave(tradeNo)
// }

// func (w *WechatOrderService) Refund(tradeNo string, req RefundReq) (*models.WechatRefundDetail, error) {
// 	oSummary, err := models.FindOrderByTradeNo(tradeNo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	oSummary.AppPayNotifyUrl = req.NotifyUrl
// 	if err = oSummary.Save(); err != nil {
// 		return nil, err
// 	}

// 	order, err := models.FindWechatOrderDetailByTradeNo(tradeNo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if order.Amount == 0 {
// 		return nil, fmt.Errorf("nothing could be refund")
// 	}

// 	resp, _, err := wxRefundService.Create(context.Background(),
// 		refunddomestic.CreateRequest{
// 			OutTradeNo:  order.TradeNo,
// 			OutRefundNo: order.TradeNo,
// 			Reason:      &req.Reason,
// 			Amount: &refunddomestic.AmountReq{
// 				Currency: core.String("CNY"),
// 				Refund:   core.Int64(int64(order.Amount)),
// 				Total:    core.Int64(int64(order.Amount)),
// 			},
// 			NotifyUrl: config.GetWxRefundNotifyUrl(tradeNo),
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	refundDetail := models.NewWechatRefundDetailByRaw(resp)
// 	return refundDetail, models.UpdateRefundDetail(refundDetail)
// }

func (w *OrderService) autoCloseOrder(order *models.Order) {
	timer := time.NewTimer(time.Until(*order.TimeExpire))
	<-timer.C
	if err := w.MustGetTrader(order.AppName, order.TradeProvider).Close(order.TradeNo); err != nil {
		logrus.WithError(err).WithField("order id", order).Error("failed to close order")
		return
	}

	order, err := w.GetOrder(order.TradeNo)
	logrus.WithError(err).WithField("order", order).Info("close order done")
	if err != nil {
		w.InvokeTradeStateChangedEvent(order)
	}
}

// ==================== Notify ============================

func (w *OrderService) PayNotifyHandler(tradeNo string, request *http.Request) error {
	transaction := new(payments.Transaction)
	notifyReq, err := wxNotifyHandler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		return err
	}

	// 处理通知内容
	logrus.WithFields(logrus.Fields{
		"summary":  notifyReq.Summary,
		"trade_no": transaction.OutTradeNo,
	}).Info("received pay notifiy")

	_, ok := enums.ParseTradeState(*transaction.TradeState)
	if !ok {
		return enums.ErrUnkownTradeState
	}

	// save order
	o, err := w.GetOrder(tradeNo)
	if err != nil {
		return err
	}

	// o, err := models.FindOrderByTradeNo(*transaction.OutTradeNo)
	// if err != nil {
	// 	return err
	// }
	// o.TradeState = *tradeState

	w.InvokeTradeStateChangedEvent(o)

	models.UpdateWechatOrderDetail(models.NewWechatOrderDetailByRaw(transaction))
	return models.GetDB().Save(o).Error
}

func (w *OrderService) RefundNotifyHandler(tradeNo string, request *http.Request) error {
	refundResp := new(refundWithRefundStatus)
	notifyReq, err := wxNotifyHandler.ParseNotifyRequest(context.Background(), request, refundResp)
	// 如果验签未通过，或者解密失败
	if err != nil {
		return err
	}
	refundResp.Status = (*refunddomestic.Status)(refundResp.RefundStatus)

	// 处理通知内容
	logrus.WithFields(logrus.Fields{
		"summary":  notifyReq.Summary,
		"trade_no": refundResp.OutTradeNo,
	}).Info("received pay notifiy")

	_, ok := enums.ParserefundState(string(*refundResp.Status))
	if !ok {
		return enums.ErrUnkownTradeState
	}

	// save order
	logrus.WithField("trade no", tradeNo).Info("get order detail and save after recieve refund notify")
	if _, err := w.GetOrder(tradeNo); err != nil {
		return err
	}

	o, err := models.FindOrderByTradeNo(*refundResp.OutTradeNo)
	if err != nil {
		return err
	}
	// o.RefundState = *refundState

	w.InvokeRefundStateChangedEvent(o)

	models.UpdateRefundDetail(models.NewWechatRefundDetailByRaw(&refundResp.Refund))
	return models.GetDB().Save(o).Error
}
