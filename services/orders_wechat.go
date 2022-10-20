package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/conflux-pay/config"
	"github.com/wangdayong228/conflux-pay/models"
	"github.com/wangdayong228/conflux-pay/models/enums"
	cns_errors "github.com/wangdayong228/conflux-pay/pay_errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
)

type WechatOrderService struct {
}

type MakeOrderReq struct {
	TradeType   enums.TradeType `json:"trade_type" binding:"required"`
	Description *string         `json:"description" binding:"required"`
	TimeExpire  int64           `json:"time_expire,omitempty" binding:"required"`
	Amount      int64           `json:"amount" binding:"required"`
}

type MakeOrderResp struct {
	TradeProvider enums.TradeProvider `json:"trade_provider"`
	TradeType     enums.TradeType     `json:"trade_type"`
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
	}
	models.GetDB().Save(order)
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
	notifyUrl := "https://www.baidu.com"
	prepayReq := PrepayRequest{
		Appid:       &app.AppId,
		Mchid:       &config.CompanyVal.MchID,
		Description: req.Description,
		OutTradeNo:  &tradeNo,
		TimeExpire:  &expire,
		Amount:      &native.Amount{Total: &req.Amount},
		NotifyUrl:   &notifyUrl,
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

func (w *WechatOrderService) GetOrderDetailAndSave(tradeNo string) (*models.WechatOrderDetail, error) {
	o, err := models.FindOrderByTradeNo(tradeNo)
	if err != nil {
		return nil, err
	}
	if o.TradeState.IsStable() {
		return models.FindWechatOrderDetailByTradeNo(tradeNo)
	} else {
		detail, err := w.getRemoteOrderDetail(tradeNo, o.TradeType)
		if err != nil {
			return nil, err
		}
		if v, ok := enums.ParseTradeState(*detail.TradeState); ok && v.IsStable() {
			models.GetDB().Save(detail)
			o.TradeState = *v
			models.GetDB().Save(o)
		}

		return detail, nil
	}
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

func (w *WechatOrderService) NotifyHandler() {

}
