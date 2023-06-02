package utils

import (
	"errors"
	"time"

	cmb_config "github.com/ahKevinXy/go-cmb/config"
	"github.com/ahKevinXy/go-cmb/handler/unit_manager"
	cmb_helper "github.com/ahKevinXy/go-cmb/help"
	cmb_models "github.com/ahKevinXy/go-cmb/models"
	"github.com/spf13/viper"
	"github.com/web3-identity/conflux-pay/config"
	"github.com/web3-identity/conflux-pay/models"
)

const (
	SUCCESSCODE = "SUC0000"
)

type CmbClient struct {
	Sm4Key            string
	Sm2PrivateKey     string
	UserId            string
	AccNbr            string
	BbkNbr            string
	UnitManagerBusmod string
}

// example output: "20230522"
func yesterdayString() string {
	yesterday := time.Now().AddDate(0, 0, -1)
	return yesterday.Format("20060102")
}

// example output: "20230523"
func todayString() string {
	return time.Now().Format("20060102")

}

func GetDefaultCmbClient() *CmbClient {
	return &CmbClient{
		UserId:            viper.GetString("company.cmb.userId"),
		Sm4Key:            viper.GetString("company.cmb.sm4Key"),
		Sm2PrivateKey:     viper.GetString("company.cmb.sm2PrivateKey"),
		AccNbr:            viper.GetString("company.cmb.accNbr"),
		BbkNbr:            viper.GetString("company.cmb.bbkNbr"),
		UnitManagerBusmod: viper.GetString("company.cmb.unitManagerBusmod"),
	}
}

func init() {
	config.Init()
	cmb_config.InitConfig(
		"",
		viper.GetString("company.cmb.apiServer"),
		"",
		"",
	)
}

func (client *CmbClient) GetUnitAccountTransHistoryListWrapper(dmanbr, begdat, enddat, ctnkey string) (*cmb_models.UnitAccountTransHistoryResponse, error) {
	res, err := unit_manager.GetUnitAccountTransHistoryList(client.UserId, client.Sm4Key, client.Sm2PrivateKey, client.AccNbr, dmanbr, begdat, enddat, ctnkey)
	if err != nil {
		return nil, err
	}
	if res.Response.Head.Resultcode != SUCCESSCODE {
		return nil, errors.New(res.Response.Head.Resultmsg)
	}
	return res, err
}

func (client *CmbClient) GetUnitAccountTransDailyListWrapper(dmanbr, ctnkey string) (*cmb_models.UnitAccountTransDailyResponse, error) {
	res, err := unit_manager.GetUnitAccountTransList(client.UserId, client.Sm4Key, client.Sm2PrivateKey, client.AccNbr, dmanbr, ctnkey)
	if err != nil {
		return nil, err
	}
	if res.Response.Head.Resultcode != SUCCESSCODE {
		return nil, errors.New(res.Response.Head.Resultmsg)
	}
	return res, err
}

func (client *CmbClient) AutoGetAllUnitAccountsTransHistoryList(begdat, enddat string) (*[]models.CmbRecord, error) {
	initRes, err := client.GetUnitAccountTransHistoryListWrapper("", begdat, enddat, "")
	if err != nil {
		return nil, err
	}
	records := models.ConvertUnitAccountTransHistoryResponseToCmbRecords(*initRes)
	nextRes := initRes
	var ctnKey string = ""

	for {
		if len(nextRes.Response.Body.Ntdmthlsy) != 0 {
			ctnKey = (nextRes.Response.Body.Ntdmthlsy[0]).Ctnkey
		} else {
			break
		}
		if ctnKey == "" {
			break
		}
		nextRes, err = client.GetUnitAccountTransHistoryListWrapper("", begdat, enddat, ctnKey)
		if err != nil {
			return nil, err
		}
		nextRecords := models.ConvertUnitAccountTransHistoryResponseToCmbRecords(*nextRes)
		records = append(records, nextRecords...)
	}

	return &records, nil
}

func (client *CmbClient) AutoGetAllUnitAccountsTransDailyList() (*[]models.CmbRecord, error) {
	initRes, err := client.GetUnitAccountTransDailyListWrapper("", "")
	if err != nil {
		return nil, err
	}
	records := models.ConvertUnitAccountTransDailyResponseToCmbRecords(*initRes)
	nextRes := initRes
	var ctnKey string = ""

	for {
		if len(nextRes.Response.Body.Ntdmtlsty) != 0 {
			ctnKey = (nextRes.Response.Body.Ntdmtlsty[0]).Ctnkey
		} else {
			break
		}
		if ctnKey == "" {
			break
		}
		nextRes, err = client.GetUnitAccountTransDailyListWrapper("", ctnKey)
		if err != nil {
			return nil, err
		}
		nextRecords := models.ConvertUnitAccountTransDailyResponseToCmbRecords(*nextRes)
		records = append(records, nextRecords...)
	}

	return &records, nil
}

// From yesterday to today
func (client *CmbClient) AutoGetRecentTransactionHistory() (*[]models.CmbRecord, error) {
	resYesterday, err := client.AutoGetAllUnitAccountsTransHistoryList(yesterdayString(), yesterdayString())
	if err != nil {
		return nil, err
	}
	resToday, err := client.AutoGetAllUnitAccountsTransDailyList()
	if err != nil {
		return nil, err
	}
	resAll := append(*resYesterday, *resToday...)
	return &resAll, nil
}

// func (client *CmbClient) QueryUnitAccountInfo(dmanbr string) {
// 	return unit_manager.QueryUnitAccountInfo()
// }

func (client *CmbClient) AddUnitAccountV1Wrapper(unitAccountName, unitAccountNbr string) (*cmb_models.AddUnitAccountV1Response, error) {
	res, err := unit_manager.AddUnitAccountV1(client.UserId, client.Sm4Key, client.Sm2PrivateKey, client.AccNbr, unitAccountName, unitAccountNbr)
	if err != nil {
		return nil, err
	}
	if res.Response.Head.Resultcode != SUCCESSCODE {
		return nil, errors.New(res.Response.Head.Resultmsg)
	}
	return res, nil
}

// useless
func (client *CmbClient) CloseUnitAccountWrapper(unitAccountNbr string) (*cmb_models.CloseUnitAccountResponse, error) {
	yurref := cmb_helper.GenYurref()

	res, err := unit_manager.CloseUnitAccount(client.UserId, client.Sm4Key, client.Sm2PrivateKey, client.AccNbr, client.BbkNbr, unitAccountNbr, client.UnitManagerBusmod, yurref)
	if err != nil {
		return nil, err
	}
	if res.Response.Head.Resultcode != SUCCESSCODE {
		return nil, errors.New(res.Response.Head.Resultmsg)
	}
	return res, nil
}

func (client *CmbClient) SetUnitAccountRelationWrapper(unitAccountNbr, relatedBankAccount string) (*cmb_models.SetUnitAccountRelationResponse, error) {
	yurref := cmb_helper.GenYurref()

	// Y: refuse
	res, err := unit_manager.SetUnitAccountRelation(client.UserId, client.Sm4Key, client.Sm2PrivateKey, client.UnitManagerBusmod, client.BbkNbr, client.AccNbr, unitAccountNbr, "R", relatedBankAccount, yurref)
	if err != nil {
		return nil, err
	}
	if res.Response.Head.Resultcode != SUCCESSCODE {
		return nil, errors.New(res.Response.Head.Resultmsg)
	}
	return res, nil
}
