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

// var (
// 	appId        = "2021003175610294"
// 	privateKey   = `MIIEowIBAAKCAQEAwJSaJY1hzdKvI7ujbMHfZvZTbtsdFqeV6mRzplnsiwRd2b9hwfqDfBJSlG4+XWDKu9l84GQ0Yx9eSp0MIVQSQRWup2cmxiC8+htYNkbWuGEUf+/yYPx805MhYl65YRCh8A9XVq7oR7yEIAXLxDzxeVmLFi1EwMllqMgD99ah8pUVICrTHIqVzLnJ0+eJWXw1TDTwK/bhzIN9ZMFptQI8htiZQQtz+uPpCqrc4O4quVNsatgSK4U0md4dMPGGniqKDJjR7rUXRFjlkPS2wsKp9FZZ9f7dD9ULPuraCgtjHvdgRXwVr47rnWsQFytsSJ/zdf7OYwGlqNhts5Ww9yFCdQIDAQABAoIBACLXf+AFaUmEsZ0kaJfXp6SIMmYfDG852Ly0edwB7vLj0lr/7h7bRQighAJIw82/Ik7ENXyfhH7egP+81CH/hOHzm0q8Nd6os6gIZHhFbrmjDsNq1Q5JAAiDWQnkG2P9T18QV3veXzYDXGAyzD/vyrxqv+g+Pm8mwNa9gUJIuboagG8TZe0Y9LJ10dkyXslWjXwOu6lBDfjABYaG6ug9iLxdNDKq52tGiOsZ6tLUKpKabLgm+ZF8ACOCYLwedGgC6cJT3kdpr3FVAousl42JBeqDoOvDe9pAGBQYQIBzQ8KLiK3XWfvEFnsJeieSo6WR79L6guqUQRMvKrFKqgY4Qu0CgYEA6IhL55hVox2o0DUHNZSIKvCWjgG71vwS4ggq5NG/3brdtnBY3VEJNpS3/JboZoF8Fcz3DXGdI6qiZCblJFLgEyo9TzBDTubEfU7N4KlsVjXV92qFDyvTpkFFx6xoqBpvvypg1NVLrLhktd0AHvWEuIdEKE59990cpf2lY/NqFw8CgYEA1AQbspI+R5G2p6QCl4bqhdfTKDHFmyFRiDLuGlWojyl3/dPvGux743Uszi4bJ3QfSv0OmjAq6zns+f6k+aJha5DGSOsNLGqGOZCQ3Os6ksx6/94CDKOD8iVt34xBPo1HSNBQ3dtnyxRwUKlX/WnJF0nKhScJEr5+RWbCFw/i7jsCgYAea+ZyUC2z/2dchfOBgQMniv5HadanU6csxyDFeuN9ILts6NnXaoioCWDgvOV+s6YGPCB+M8T5K5O/Qo9r5yPFnhsTRx8nLW27bxnkMIYp6TUq/1aVG4i/EX8NlnLCu2KvQd4VOiqCWEVkvZsMcdaBRcEW/N3iFZ1v4fVHVEsm5QKBgQDD0JuIPRu6XDFX+dnO+3PFdEV4/ScmFQrJgUh6GB0bRFCnpcNTmZD+zm04bEr2EIEKcFi5Pb2WDaT6bB8Q1NGnWEpadIVxPV2E8ylocPVjOepsQS6hX7Bwx/MHofFshW2OKaBWl9rwLItjZFR5H+fzU1rxydDOeBQFo1elly2fmwKBgGZ8f6fEl6grDd8Van38G8D174GKiC0r7I1+rBgcy7Xv3sWdJ1wmsUuBh0/NzZ6S1q/NshdFPhKz+8o/vylxIcbWVI+jAJytg4L0yHeDU0Zv2tpZMTMmYZs9MvbnwWxopw6BpbKmX5tCpvU+xwarZ4sd6QVx/TI02HYdewL8C7JL`
// 	alipayPubkey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArFZVasOUPu7ny4YU3/0cW78UZM6Gn5js+Y2VbE+QBJYZL0Ve2XQWLQOiDMXSnzwYaJp7qGl6rwcVir6JMqy2mmXT0DsILGXo0c9nAAgn/Ve6GRHtdRTHydVGKN/qPtUXP280EVCYBVZ9MKWyQ630xHSXxyl0H9JMTBfCjKPlemjRIJ6YAX57iGFXelqLf691gx7gD73i1aAB6XZIA1mCgTGFao0dzKgDD7o+EKBStNKMCGbvKsm6RFCfockmq/5HkaULyQvPFflTwvJzBbZDytU7udtH0veMGNN/KOdmRFfpBBtR0J1ircIMsLgpb52YTLflGtOQ5TkT2L555j6oLwIDAQAB"
// )

type AlipayTrader struct {
	app    config.App
	client *alipay.Client
	// clients          map[string]*alipay.Client
	// clientCreateOnce sync.Once
}

// func (a *AlipayTrader) getClient(appId string) *alipay.Client {
// 	if _, ok := a.clients[appId]; !ok {
// 		c, err := alipay.New(appId, config.CompanyVal.Alipay.PrivateKey, true)
// 		if err != nil {
// 			return nil
// 		}

// 		a.clientCreateOnce.Do(func() {
// 			a.clients[appId] = c
// 		})
// 	}
// 	return a.clients[appId]
// }

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

	// client, err := alipay.New(appId, privateKey, true)
	// client.LoadAliPayPublicKey(alipayPubkey)
	// if err != nil {
	// 	return nil, err
	// }

	trader := &AlipayTrader{
		app:    app,
		client: client,
	}

	// tradeNo := "20230320063232117"
	// p := alipay.TradeQuery{OutTradeNo: tradeNo}
	// // p.OutTradeNo = "20230320063232117"
	// res, err := client.TradeQuery(p)
	// fmt.Printf("res %v\nerr %v\n", res, err)

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

	_, err = a.GetRefundState(tradeNo)
	if err != nil {
		return convertAlTradeState(res.Content.TradeStatus, false), nil
	}

	return convertAlTradeState(res.Content.TradeStatus, true), nil
}

// refund
func (a *AlipayTrader) Refund(tradeNo string, req RefundReq) error {
	tq := alipay.TradeQuery{OutTradeNo: tradeNo}
	orderRes, err := a.client.TradeQuery(tq)
	if err != nil {
		return err
	}
	if !orderRes.IsSuccess() {
		return fmt.Errorf("msg:%v, sub msg:%v", orderRes.Content.Msg, orderRes.Content.SubMsg)
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
		return fmt.Errorf("msg:%v, sub msg:%v", refundRes.Content.Msg, refundRes.Content.SubMsg)
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
		return enums.REFUND_STATE_NIL, fmt.Errorf("msg:%v, sub msg:%v", res.Content.Msg, res.Content.SubMsg)
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
	switch status {
	case alipay.TradeStatusWaitBuyerPay:
		return enums.TRADE_STATE_NOTPAY
	case alipay.TradeStatusClosed:
		if !isRefund {
			return enums.TRADE_STATE_CLOSED
		}
		return enums.TRADE_STATE_REFUND
	case alipay.TradeStatusSuccess:
		return enums.TRADE_STATE_SUCCESSS
	// TODO: 交易结束不可退款，微信没有对应状态，暂用SUCCESS
	case alipay.TradeStatusFinished:
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

func IsNotExistErr(err error) bool {
	return strings.Contains(err.Error(), "查询的交易不存在")
}
