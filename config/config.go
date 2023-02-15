package config

import (
	// "github.com/Conflux-Chain/go-conflux-util/viper"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/web3-identity/conflux-pay/models/enums"
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

	NotifyConfig = make(map[enums.TradeProvider]*InNotifyItem)
	NotifyConfig[enums.TRADE_PROVIDER_WECHAT] = getOrderConfig("wechat")
	NotifyConfig[enums.TRADE_PROVIDER_ALIPAY] = getOrderConfig("alipay")

	fmt.Printf("init config done, notify config %v\n", NotifyConfig)
}

var (
	CompanyVal   *Company
	Apps         map[string]App
	NotifyConfig map[enums.TradeProvider]*InNotifyItem
)

type Company struct {
	Wechat CompanyWechat
	Alipay CompanyAlipay
}

type CompanyWechat struct {
	MchID         string
	MchCertNo     string
	MchApiV3Key   string
	MchPrivateKey string
}

type CompanyAlipay struct {
	PrivateKey      string
	AlipayPublicKey string
}

type App struct {
	AppIdAlipay   string
	AppIdWechat   string
	AppSecretHash string
	AppInternalID uint
}

type InNotifyItem struct {
	PayNotifyUrlBase    string
	RefundNotifyUrlBase string
	Enable              bool
}

func getCompany() *Company {
	var v Company
	viper.UnmarshalKey("company", &v)
	return &v
}

func getApps() map[string]App {
	var apps map[string]App
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		panic(err)
	}
	return apps
}

// providerName maybe wechat/alipay/bank
func getOrderConfig(providerName string) *InNotifyItem {
	order := viper.GetViper().Sub("inNotify")
	var notify InNotifyItem
	if err := order.UnmarshalKey(providerName, &notify); err != nil {
		panic(err)
	}
	return &notify
}

func MustGetApp(appName string) App {
	v, ok := Apps[appName]
	if !ok {
		panic("not exists")
	}
	return v
}

func GetPayNotifyUrl(provider enums.TradeProvider, tradeNo string) *string {
	if !NotifyConfig[provider].Enable {
		getDefaultNotifyUrl(provider)
	}
	v := fmt.Sprintf("%v%v%v", NotifyConfig[provider].PayNotifyUrlBase, "/v0/orders/wechat/notify-pay/", tradeNo)
	return &v
}

func GetRefundNotifyUrl(provider enums.TradeProvider, tradeNo string) *string {
	if !NotifyConfig[provider].Enable {
		getDefaultNotifyUrl(provider)
	}
	v := fmt.Sprintf("%v%v%v", NotifyConfig[provider].RefundNotifyUrlBase, "/v0/orders/wechat/notify-refund/", tradeNo)
	return &v
}

func getDefaultNotifyUrl(provider enums.TradeProvider) *string {
	invalid := "https://a.com"
	if provider == enums.TRADE_PROVIDER_WECHAT {
		return &invalid
	}
	return nil
}
