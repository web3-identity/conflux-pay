package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/conflux-pay/config"
	"github.com/wangdayong228/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	wxClient        *core.Client
	wxNativeService native.NativeApiService
	wxH5Service     h5.H5ApiService
)

func Init() {
	var err error
	if wxClient, err = newWechatClient(); err != nil {
		logrus.WithError(err).Panic("failed creat wechat client")
		panic(err)
	}

	wxNativeService = native.NativeApiService{Client: wxClient}
	wxH5Service = h5.H5ApiService{Client: wxClient}
}

// 交易号生成规则
// 时间戳 + 上游服务方ID
// TODO: 防止并发导致的重复
func genTradeNo(userID uint, provider enums.TradeProvider) string {
	return fmt.Sprintf("%s%13d%05d", provider.Code(), time.Now().UnixMilli(), userID)
}

func newWechatClient() (*core.Client, error) {
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	company := config.CompanyVal

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(company.MchPrivateKey)
	if err != nil {
		return nil, err
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(
			company.MchID,
			company.MchCertificateSerialNumber,
			mchPrivateKey,
			company.MchAPIv3Key),
	}
	return core.NewClient(ctx, opts...)
}
