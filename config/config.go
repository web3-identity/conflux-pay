package config

import (
	// "github.com/Conflux-Chain/go-conflux-util/viper"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/conflux-pay/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.conflux-pay") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	viper.AddConfigPath("..")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}

	CompanyVal = getCompany()
	Apps = getApps()
	WechatOrderConfig = getOrderConfig("wechat")

	logrus.WithField("WechatOrderConfig", WechatOrderConfig).Info("init config done")
}

var (
	CompanyVal        *Company
	Apps              map[string]App
	WechatOrderConfig *OrderItem
)

type Company struct {
	MchID                      string
	MchCertificateSerialNumber string
	MchAPIv3Key                string
	MchPrivateKey              string
}

type App struct {
	AppId         string
	AppSecretHash string
	AppInternalID uint
}

type OrderItem struct {
	PayNotifyUrlBase    string
	RefundNotifyUrlBase string
	UpdateUseNotify     bool
}

func getCompany() *Company {
	sub := viper.GetViper().Sub("company")
	return &Company{
		MchID:                      sub.GetString("mchid"),
		MchCertificateSerialNumber: sub.GetString("mchCertNo"),
		MchAPIv3Key:                sub.GetString("mchAPIv3Key"),
		MchPrivateKey:              sub.GetString("mchPrivateKey"),
	}
}

func getApps() map[string]App {
	var apps map[string]App
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		panic(err)
	}
	return apps
}

// providerName maybe wechat/alipay/bank
func getOrderConfig(providerName string) *OrderItem {
	order := viper.GetViper().Sub("order")
	var wx OrderItem
	order.UnmarshalKey("wechat", &wx)
	return &wx
}

func MustGetApp(appName string) App {
	v, ok := Apps[appName]
	if !ok {
		panic("not exists")
	}
	return v
}

func GetWxPayNotifyUrl(tradeNo string) *string {
	v := fmt.Sprintf("%v%v%v", WechatOrderConfig.PayNotifyUrlBase, "/v0/orders/wechat/notify-pay/", tradeNo)
	return &v
}

func GetWxRefundNotifyUrl(tradeNo string) *string {
	v := fmt.Sprintf("%v%v%v", WechatOrderConfig.RefundNotifyUrlBase, "/v0/orders/wechat/notify-refund/", tradeNo)
	return &v
}
