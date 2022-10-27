package models

import (
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type WechatRefundDetail struct {
	BaseModel
	RefundId *string `json:"refund_id"`
	// 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundNo string `json:"out_refund_no" gorm:"uniqueIndex"`
	// 微信支付交易订单号
	TransactionId *string `json:"transaction_id"`
	// 原支付交易对应的商户订单号
	TradeNo string `json:"out_trade_no"`
	// 枚举值： - ORIGINAL—原路退款 - BALANCE—退回到余额 - OTHER_BALANCE—原账户异常退到其他余额账户 - OTHER_BANKCARD—原银行卡异常退到其他银行卡 * `ORIGINAL` - 原路退款 * `BALANCE` - 退回到余额 * `OTHER_BALANCE` - 原账户异常退到其他余额账户 * `OTHER_BANKCARD` - 原银行卡异常退到其他银行卡
	Channel *string `json:"channel"`
	// 取当前退款单的退款入账方，有以下几种情况： 1）退回银行卡：{银行名称}{卡类型}{卡尾号} 2）退回支付用户零钱:支付用户零钱 3）退还商户:商户基本账户商户结算银行账户 4）退回支付用户零钱通:支付用户零钱通
	UserReceivedAccount *string `json:"user_received_account"`
	// 退款成功时间，退款状态status为SUCCESS（退款成功）时，返回该字段。遵循rfc3339标准格式，格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，T出现在字符串中，表示time元素的开头，HH:mm:ss表示时分秒，TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。
	SuccessTime *time.Time `json:"success_time,omitempty"`
	// 退款受理时间，遵循rfc3339标准格式，格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，T出现在字符串中，表示time元素的开头，HH:mm:ss表示时分秒，TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。
	CreateTime *time.Time `json:"create_time"`
	// 退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，手动处理此笔退款。 枚举值： - SUCCESS—退款成功 - CLOSED—退款关闭 - PROCESSING—退款处理中 - ABNORMAL—退款异常 * `SUCCESS` - 退款成功 * `CLOSED` - 退款关闭 * `PROCESSING` - 退款处理中 * `ABNORMAL` - 退款异常
	Status *string `json:"status"`
	// 退款所使用资金对应的资金账户类型 枚举值： - UNSETTLED : 未结算资金 - AVAILABLE : 可用余额 - UNAVAILABLE : 不可用余额 - OPERATION : 运营户 - BASIC : 基本账户（含可用余额和不可用余额） * `UNSETTLED` - 未结算资金 * `AVAILABLE` - 可用余额 * `UNAVAILABLE` - 不可用余额 * `OPERATION` - 运营户 * `BASIC` - 基本账户（含可用余额和不可用余额）
	FundsAccount *string `json:"funds_account,omitempty"`
	// 金额详细信息
	Amount *int64 `json:"amount"`
}

func NewWechatRefundDetailByRaw(raw *refunddomestic.Refund) *WechatRefundDetail {
	return &WechatRefundDetail{
		RefundId:            raw.RefundId,
		RefundNo:            *raw.OutRefundNo,
		TransactionId:       raw.TransactionId,
		TradeNo:             *raw.OutTradeNo,
		Channel:             (*string)(raw.Channel),
		UserReceivedAccount: raw.UserReceivedAccount,
		SuccessTime:         raw.SuccessTime,
		CreateTime:          raw.CreateTime,
		Status:              (*string)(raw.Status),
		FundsAccount:        (*string)(raw.FundsAccount),
		Amount:              raw.Amount.Total,
	}
}

func FindRefundDetailByTradeNo(tradeNo string) (*WechatRefundDetail, error) {
	o := WechatRefundDetail{}
	o.TradeNo = tradeNo
	return &o, GetDB().Where(&o).First(&o).Error
}

func UpdateRefundDetail(val *WechatRefundDetail) error {
	refundInDb, err := FindRefundDetailByTradeNo(val.TradeNo)
	if err != nil {
		return err
	}

	val.BaseModel = refundInDb.BaseModel
	return GetDB().Save(val).Error
}
