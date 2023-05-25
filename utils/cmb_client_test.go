package utils

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCmbClient(t *testing.T) {
	client := GetDefaultCmbClient()
	_, err := client.GetUnitAccountTransHistoryListWrapper("", "20230514", "20230515", "")
	assert.NoError(t, err)
	time.Sleep(time.Second * 11)
	_, err = client.AutoGetAllUnitAccountsTransHistoryList("20230514", "20230520")
	assert.NoError(t, err)
	time.Sleep(time.Second * 11)
	_, err = client.AutoGetRecentTransactionHistory()
	assert.NoError(t, err)
}

func genRandUnitAccountNbr() string {
	rand.Seed(time.Now().UnixNano())          // 设置随机种子
	num := rand.Intn(9000000000) + 1000000000 // 生成10位随机数
	return strconv.Itoa(num)
}

func TestSetCmbUnitAccountRelation(t *testing.T) {
	client := GetDefaultCmbClient()
	dmaNbr := genRandUnitAccountNbr()
	_, err := client.AddUnitAccountV1Wrapper("子账户名", dmaNbr)
	print(client.AccNbr)
	print(dmaNbr)
	assert.NoError(t, err)

	_, err = client.SetUnitAccountRelationWrapper(dmaNbr, "6214835982629402")
	assert.NoError(t, err)

	_, err = client.CloseUnitAccountWrapper(dmaNbr)
	assert.NoError(t, err)
}
