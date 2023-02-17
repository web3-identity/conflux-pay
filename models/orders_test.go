package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/web3-identity/conflux-pay/config"
)

func TestGetExpireOrder(t *testing.T) {
	var orders []*Order
	err := GetDB().Debug().
		Where("time_expire < ?", time.Now()).
		Find(&orders).Error
	assert.NoError(t, err)
}

func init() {
	config.Init()
	ConnectDB()
}
