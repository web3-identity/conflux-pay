package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	cns_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type TradeStateChangeHandler func(o *models.Order)
type RefundStateChangeHandler func(o *models.Order)

type MakeOrderReq struct {
	TradeType   enums.TradeType `json:"trade_type" binding:"required" swaggertype:"string"`
	Description *string         `json:"description" binding:"required"`
	TimeExpire  int64           `json:"time_expire,omitempty" binding:"required"`
	Amount      int64           `json:"amount" binding:"required"`
	NotifyUrl   *string         `json:"notify_url"`
}

type MakeOrderResp struct {
	TradeProvider enums.TradeProvider `json:"trade_provider" swaggertype:"string"`
	TradeType     enums.TradeType     `json:"trade_type" swaggertype:"string"`
	TradeNo       string              `json:"trade_no"`
	CodeUrl       *string             `json:"code_url,omitempty"`
	H5Url         *string             `json:"h5_url,omitempty"`
}

func NewMakeOrderRespFromRaw(raw *models.Order) *MakeOrderResp {
	return &MakeOrderResp{
		TradeProvider: raw.Provider,
		TradeType:     raw.TradeType,
		TradeNo:       raw.TradeNo,
		CodeUrl:       raw.CodeUrl,
		H5Url:         raw.H5Url,
	}
}

type PrepayRequest native.PrepayRequest

func (r *PrepayRequest) toH5() (val *h5.PrepayRequest) {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(j, &val); err != nil {
		panic(err)
	}
	return
}

type WechatOrderService struct {
	TradeStateChangeEvent  []TradeStateChangeHandler
	RefundStateChangeEvent []RefundStateChangeHandler
}

func (w *WechatOrderService) RegisterTradeStateChangeEvent(h TradeStateChangeHandler) {
	w.TradeStateChangeEvent = append(w.TradeStateChangeEvent, h)
	// w.TradeStateChangeEvent = append(w.TradeStateChangeEvent, w.GetOrderDetailAndSave)
}

func (w *WechatOrderService) RegisterRefundStateChangeEvent(h RefundStateChangeHandler) {
	w.RefundStateChangeEvent = append(w.RefundStateChangeEvent, h)
}

func (w *WechatOrderService) InvokeTradeStateChangedEvent(o *models.Order) {
	fmt.Println("invoke trade state changed")
	for _, h := range w.TradeStateChangeEvent {
		h(o)
	}
}

func (w *WechatOrderService) InvokeRefundStateChangedEvent(o *models.Order) {
	fmt.Println("invoke refund state changed")
	for _, h := range w.RefundStateChangeEvent {
		h(o)
	}
}

func NewWechatOrderService() *WechatOrderService {
	w := WechatOrderService{}
	w.RegisterTradeStateChangeEvent(func(o *models.Order) {
		go runPayNotifyTask(o)
		go runRefundNotifyTask(o) // 微信支付只通知成功的状态，在交易关闭后也需要处理refund notify标志
	})
	w.RegisterRefundStateChangeEvent(func(o *models.Order) {
		go runRefundNotifyTask(o)
	})
	return &w
}

// 统一wechat所有下单接口
// 生成 trade_no
// 返回 trade_no, pay_url

func (w *WechatOrderService) MakeOrder(appName string, req MakeOrderReq) (*models.Order, error) {
	app := config.MustGetApp(appName)
	no := genTradeNo(app.AppInternalID, enums.TRADE_PROVIDER_WECHAT)

	expire := time.Unix(req.TimeExpire, 0)
	orderResp, err := w.prePay(appName, no, req)
	if err != nil {
		return nil, err
	}

	// save db
	order := &models.Order{
		OrderCore: models.OrderCore{
			AppName:     appName,
			Provider:    enums.TRADE_PROVIDER_WECHAT,
			TradeNo:     no,
			TradeType:   req.TradeType,
			TradeState:  enums.TRADE_STATE_NOTPAY,
			Amount:      uint(req.Amount),
			Description: req.Description,
			TimeExpire:  &expire,
			CodeUrl:     orderResp.CodeUrl,
			H5Url:       orderResp.H5Url,
		},
		OrderNofity: models.OrderNofity{
			AppPayNotifyUrl: req.NotifyUrl,
		},
	}

	models.GetDB().Save(order)

	detail, err := w.getRemoteOrderDetail(order.TradeNo, order.TradeType)
	if err != nil {
		return nil, err
	}
	models.GetDB().Save(detail)

	go w.autoCloseOrder(order)

	return order, nil
}

func (w *WechatOrderService) RefreshUrl(tradeNo string) (*MakeOrderResp, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}
	// return error if order is complete
	if o.TradeState.IsStable() {
		return nil, cns_errors.ERR_ORDER_COMPLETED
	}

	resp, err := w.prePay(o.AppName, tradeNo, MakeOrderReq{
		TradeType:   o.TradeType,
		Description: o.Description,
		TimeExpire:  o.TimeExpire.Unix(),
		Amount:      int64(o.Amount),
	})
	if err != nil {
		return nil, err
	}

	o.CodeUrl = resp.CodeUrl
	o.H5Url = resp.H5Url
	models.GetDB().Save(o)
	return resp, nil
}

func (w *WechatOrderService) prePay(appName string, tradeNo string, req MakeOrderReq) (*MakeOrderResp, error) {
	app := config.MustGetApp(appName)

	expire := time.Unix(req.TimeExpire, 0)
	prepayReq := PrepayRequest{
		Appid:       &app.AppId,
		Mchid:       &config.CompanyVal.MchID,
		Description: req.Description,
		OutTradeNo:  &tradeNo,
		TimeExpire:  &expire,
		Amount:      &native.Amount{Total: &req.Amount},
		NotifyUrl:   config.GetWxPayNotifyUrl(tradeNo),
	}

	orderResp := &MakeOrderResp{
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
		orderResp.CodeUrl = resp.CodeUrl
	case enums.TRADE_TYPE_H5:
		resp, _, err := wxH5Service.Prepay(context.Background(), *prepayReq.toH5())
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("h5 pay complete")
		if err != nil {
			return nil, err
		}
		orderResp.H5Url = resp.H5Url
	default:
		return nil, enums.ErrUnkownTradeType
	}
	return orderResp, nil
}

func (w *WechatOrderService) GetOrderDetail(tradeNo string) (*models.WechatOrderDetail, error) {
	if config.WechatOrderConfig.UpdateUseNotify {
		return models.FindWechatOrderDetailByTradeNo(tradeNo)
	}
	return w.GetOrderDetailAndSave(tradeNo)
}

func (w *WechatOrderService) GetOrderDetailAndSave(tradeNo string) (*models.WechatOrderDetail, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}

	logrus.WithField("current order", o).Info("will get order detail and save")

	if o.IsStable() {
		return models.FindWechatOrderDetailByTradeNo(tradeNo)
	}

	// if o.TradeState.IsStable() {

	// 	fmt.Println("aaaaa")
	// 	oDetail, err := models.FindWechatOrderDetailByTradeNo(tradeNo)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("bbbbb")
	// 	if o.RefundState.IsStable(o.TradeState) {
	// 		return oDetail, nil
	// 	}
	// 	fmt.Println("ccccc")
	// 	// refresh refund status
	// 	refundDetial, err := w.getRemoteRefundDetail(tradeNo)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("ddddd")
	// 	models.UpdateRefundDetail(refundDetial)

	// 	if v, ok := enums.ParserefundState(*refundDetial.Status); ok && *v != o.RefundState {
	// 		o.RefundState = *v
	// 		models.GetDB().Save(o)
	// 	}
	// 	return oDetail, nil
	// } else {
	// fmt.Println("eeeee")
	detail, err := w.getRemoteOrderDetail(tradeNo, o.TradeType)
	if err != nil {
		return nil, err
	}
	// fmt.Println("fffff")
	v, ok := enums.ParseTradeState(*detail.TradeState)
	if !ok {
		return nil, fmt.Errorf("unknown trade state %v", *detail.TradeState)
	}

	// fmt.Println("ggggg")
	if *v != o.TradeState {
		o.TradeState = *v
		models.UpdateWechatOrderDetail(detail)
		models.GetDB().Save(o)
		logrus.WithField("trade_no", o.TradeNo).WithField("trade_state", o.TradeState).Info("update order and detail")
	}

	if o.TradeState != enums.TRADE_STATE_REFUND {
		return detail, nil
	}

	refundDetial, err := w.getRemoteRefundDetail(tradeNo)
	if err != nil {
		return nil, err
	}
	// fmt.Println("ddddd")
	models.UpdateRefundDetail(refundDetial)

	if v, ok := enums.ParserefundState(*refundDetial.Status); ok && *v != o.RefundState {
		o.RefundState = *v
		models.GetDB().Save(o)
	}

	return detail, nil
	// }
}

func (w *WechatOrderService) getRemoteOrderDetail(tradeNo string, tradeType enums.TradeType) (*models.WechatOrderDetail, error) {
	switch tradeType {
	case enums.TRADE_TYPE_NATIVE:
		resp, _, err := wxNativeService.QueryOrderByOutTradeNo(context.Background(), native.QueryOrderByOutTradeNoRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &tradeNo,
		})
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
		if err != nil {
			return nil, err
		}
		return models.NewWechatOrderDetailByRaw(resp), nil
	case enums.TRADE_TYPE_H5:
		resp, _, err := wxH5Service.QueryOrderByOutTradeNo(context.Background(), h5.QueryOrderByOutTradeNoRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &tradeNo,
		})
		logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order detail complete remote")
		if err != nil {
			return nil, err
		}
		return models.NewWechatOrderDetailByRaw(resp), nil

	default:
		return nil, enums.ErrUnkownTradeType
	}
}

func (w *WechatOrderService) getRemoteRefundDetail(tradeNo string) (*models.WechatRefundDetail, error) {
	req := refunddomestic.QueryByOutRefundNoRequest{OutRefundNo: &tradeNo}
	resp, _, err := wxRefundService.QueryByOutRefundNo(context.Background(), req)
	logrus.WithField("trade_no", tradeNo).WithField("response", resp).WithError(err).Info("query order refund detail complete remote")
	if err != nil {
		return nil, err
	}
	return models.NewWechatRefundDetailByRaw(resp), nil
}

func (w *WechatOrderService) Close(tradeNo string) (*models.WechatOrderDetail, error) {
	order, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}

	switch order.TradeType {
	case enums.TRADE_TYPE_NATIVE:
		result, err := wxNativeService.CloseOrder(context.Background(), native.CloseOrderRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &order.TradeNo,
		})
		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
		if err != nil {
			return nil, err
		}

	case enums.TRADE_TYPE_H5:
		result, err := wxH5Service.CloseOrder(context.Background(), h5.CloseOrderRequest{
			Mchid:      &config.CompanyVal.MchID,
			OutTradeNo: &order.TradeNo,
		})
		logrus.WithField("trade_no", order.TradeNo).WithField("result", result).WithError(err).Info("close order complete remote")
		if err != nil {
			return nil, err
		}
	default:
		return nil, enums.ErrUnkownTradeType
	}
	return w.GetOrderDetailAndSave(tradeNo)
}

type RefundReq struct {
	Reason    string  `json:"reason" binding:"required"`
	NotifyUrl *string `json:"notify_url"`
}

func (w *WechatOrderService) Refund(tradeNo string, req RefundReq) (*models.WechatRefundDetail, error) {
	oSummary, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}

	oSummary.AppPayNotifyUrl = req.NotifyUrl
	if err = oSummary.Save(); err != nil {
		return nil, err
	}

	order, err := models.FindWechatOrderDetailByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}

	if order.Amount == 0 {
		return nil, fmt.Errorf("nothing could be refund")
	}

	resp, _, err := wxRefundService.Create(context.Background(),
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
		return nil, err
	}

	refundDetail := models.NewWechatRefundDetailByRaw(resp)
	return refundDetail, models.UpdateRefundDetail(refundDetail)
}

func (w *WechatOrderService) autoCloseOrder(order *models.Order) {
	timer := time.NewTimer(time.Until(*order.TimeExpire))
	<-timer.C
	if _, err := w.Close(order.TradeNo); err != nil {
		logrus.WithError(err).WithField("order id", order).Error("failed to close order")
		return
	}
	logrus.WithField("order id", order).Info("close order successed")
	w.GetOrderDetailAndSave(order.TradeNo)
	w.InvokeTradeStateChangedEvent(order)
}

// ==================== Notify ============================

func (w *WechatOrderService) PayNotifyHandler(tradeNo string, request *http.Request) error {
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
	if _, err := w.GetOrderDetailAndSave(tradeNo); err != nil {
		return err
	}

	o, err := models.FindOrderByTradeNo(*transaction.OutTradeNo)
	if err != nil {
		return err
	}
	// o.TradeState = *tradeState

	w.InvokeTradeStateChangedEvent(o)

	models.UpdateWechatOrderDetail(models.NewWechatOrderDetailByRaw(transaction))
	return models.GetDB().Save(o).Error
}

type refundWithRefundStatus struct {
	refunddomestic.Refund
	RefundStatus *string `gorm:"-" json:"refund_status"`
}

func (w *WechatOrderService) RefundNotifyHandler(tradeNo string, request *http.Request) error {
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
	if _, err := w.GetOrderDetailAndSave(tradeNo); err != nil {
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
