package services

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/utils"
	"gorm.io/gorm"
)

func StartCmbTransTracker() {
	cmbClient := utils.GetDefaultCmbClient()
	for {
		// 检查当前时间是否为整点或半点
		if time.Now().Minute() != 0 && time.Now().Minute() != 30 {
			trackRecentTransactionHistory(cmbClient)
		}
		time.Sleep(60 * time.Second)
	}
}

func trackRecentTransactionHistory(client *utils.CmbClient) {
	records, err := client.AutoGetRecentTransactionHistory()
	if err != nil {
		logrus.Error(err)
		return
	}
	go updateTableFromRecentRecords(records)
}

func updateTableFromRecentRecords(records *[]models.CmbRecord) {
	db := models.GetDB()
	for _, record := range *records {
		var existingRecord models.CmbRecord
		result := db.Where("trx_nbr = ?", record.TrxNbr).First(&existingRecord)
		// if record is found, don't do anything
		if result.Error == nil {
			continue
		}
		if result.Error != gorm.ErrRecordNotFound {
			logrus.WithError(result.Error).Error("unexpected error")
		}
		if err := db.Create(&record).Error; err != nil {
			logrus.Error(err)
		}
	}
}
