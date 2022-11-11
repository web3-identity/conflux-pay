package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models/enums"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	wxClient        *core.Client
	wxNativeService native.NativeApiService
	wxH5Service     h5.H5ApiService
	wxRefundService refunddomestic.RefundsApiService
	wxNotifyHandler *notify.Handler
)

func Init() {
	var err error
	if wxClient, wxNotifyHandler, err = newWechatClient(); err != nil {
		logrus.WithError(err).Panic("failed creat wechat client")
		panic(err)
	}

	wxNativeService = native.NativeApiService{Client: wxClient}
	wxH5Service = h5.H5ApiService{Client: wxClient}
	wxRefundService = refunddomestic.RefundsApiService{Client: wxClient}
}

func StartTasks() {
	InitCloseOrderTask()
	go RunNotifyTask()
}

// 交易号生成规则
// 时间戳 + 上游服务方ID
// TODO: 防止并发导致的重复
func genTradeNo(userID uint, provider enums.TradeProvider) string {
	return fmt.Sprintf("%s%13d%05d", provider.Code(), time.Now().UnixMilli(), userID)
}

func newWechatClient() (*core.Client, *notify.Handler, error) {
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	company := config.CompanyVal

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKey(company.MchPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(
			company.MchID,
			company.MchCertificateSerialNumber,
			mchPrivateKey,
			company.MchAPIv3Key),
	}

	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(company.MchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(company.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, nil, err
	}

	return client, handler, nil
}
