package services

import (
	"encoding/json"

	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
)

type TradeStateChangeHandler func(o *models.Order)
type RefundStateChangeHandler func(o *models.Order)

type MakeOrderReq struct {
	TradeProvider enums.TradeProvider `json:"trade_provider" swaggertype:"string"`
	TradeType     enums.TradeType     `json:"trade_type" binding:"required" swaggertype:"string"`
	Description   *string             `json:"description" binding:"required"`
	TimeExpire    int64               `json:"time_expire,omitempty" binding:"required"`
	Amount        int64               `json:"amount" binding:"required"`
	NotifyUrl     *string             `json:"notify_url"`
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

type RefundReq struct {
	Reason    string  `json:"reason" binding:"required"`
	NotifyUrl *string `json:"notify_url"`
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
