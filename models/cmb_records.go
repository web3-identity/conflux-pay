package models

import (
	"errors"
	"strconv"
	"time"

	cmb_models "github.com/ahKevinXy/go-cmb/models"
)

type CmbRecord struct {
	BaseModel
	AccNbr string  `gorm:"type:varchar(35);not null"`
	DmaNbr string  `gorm:"type:varchar(20);not null"`             // sub unit number
	DmaNam string  `gorm:"type:varchar(82)"`                      // sub unit name
	TrxNbr string  `gorm:"type:varchar(15);uniqueIndex;not null"` // unique
	CcyNbr string  `gorm:"type:varchar(2);not null"`              // expected to be rmb
	TrxAmt float64 `gorm:"type:decimal(15,2);not null"`
	TrxDir string  `gorm:"type:varchar(1);not null"` // tx direction
	TrxDat string  `gorm:"type:date;not null"`
	TrxTim string  `gorm:"type:varchar(6);not null"`
	AutFlg string  `gorm:"type:varchar(1);not null"` // useless
	RpyAcc string  `gorm:"type:varchar(35)"`         // useless
	RpyNam string  `gorm:"type:varchar(62)"`         // useless
	TrxTxt string  `gorm:"type:varchar(42)"`         // txt that sender appended
	NarInn string  `gorm:"type:varchar(20)"`         // useless
}

func ConvertUnitAccountTransHistoryResponseToCmbRecords(res cmb_models.UnitAccountTransHistoryResponse) []CmbRecord {
	var records []CmbRecord
	for _, n := range res.Response.Body.Ntdmthlsz {
		trxAmt, _ := strconv.ParseFloat(n.Trxamt, 64)
		records = append(records, CmbRecord{
			AccNbr: n.Accnbr,
			DmaNbr: n.Dmanbr,
			DmaNam: n.Dmanam,
			TrxNbr: n.Trxnbr,
			CcyNbr: n.Ccynbr,
			TrxAmt: trxAmt,
			TrxDir: n.Trxdir,
			TrxDat: n.Trxdat,
			TrxTim: n.Trxtim,
			AutFlg: n.Autflg,
			RpyAcc: n.Rpyacc,
			RpyNam: n.Rpynam,
			TrxTxt: n.Trxtxt,
			NarInn: n.Narinn,
		})
	}
	return records
}

func ConvertUnitAccountTransDailyResponseToCmbRecords(res cmb_models.UnitAccountTransDailyResponse) []CmbRecord {
	var records []CmbRecord
	for _, n := range res.Response.Body.Ntdmtlstz {
		trxAmt, _ := strconv.ParseFloat(n.Trxamt, 64)
		records = append(records, CmbRecord{
			AccNbr: n.Accnbr,
			DmaNbr: n.Dmanbr,
			DmaNam: n.Dmanam,
			TrxNbr: n.Trxnbr,
			CcyNbr: n.Ccynbr,
			TrxAmt: trxAmt,
			TrxDir: n.Trxdir,
			TrxDat: time.Now().Format("20060102"),
			TrxTim: n.Trxtim,
			AutFlg: n.Autflg,
			RpyAcc: n.Rpyacc,
			RpyNam: n.Rpynam,
			TrxTxt: n.Trxtxt,
			NarInn: n.Narinn,
		})
	}
	return records
}

func isValidDateFormat(input string) bool {
	_, err := time.Parse("20060102", input)
	return err == nil
}

func GetCmbRecords(unitAccountNbr string, transactionDate string, transactionDirection string, limit, offset int) (*[]CmbRecord, error) {

	tmp := db
	// do Date Filter
	if transactionDate != "" {
		if !isValidDateFormat(transactionDate) {
			return nil, errors.New("invalid transaction date, requires 20060102")
		}
		tmp = tmp.Where("trx_dat = ?", transactionDate)
	}
	if transactionDirection != "" {
		if transactionDirection != "C" && transactionDirection != "D" {
			return nil, errors.New("invalid transaction direction, requires C or D")
		}
		tmp = tmp.Where("trx_dir = ?", transactionDirection)
	}
	if unitAccountNbr != "" {
		tmp = tmp.Where("dma_nbr = ?", unitAccountNbr)
	}

	var records *[]CmbRecord
	err := tmp.
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error
	return records, err
}

func GetTodayAndYesterdayRecords(unitAccountNbr string, transactionDirection string, limit, offset int) (*[]CmbRecord, error) {

	tmp := db.Table("cmb_records")
	// do Date Filter
	tmp = tmp.Where("trx_dat = ? or trx_dat = ?", time.Now().Format("20060102"), (time.Now().AddDate(0, 0, -1)).Format("20060102"))
	if transactionDirection != "" {
		if transactionDirection != "C" && transactionDirection != "D" {
			return nil, errors.New("invalid transaction direction, requires C or D")
		}
		tmp = tmp.Where("trx_dir = ?", transactionDirection)
	}
	if unitAccountNbr != "" {
		tmp = tmp.Where("dma_nbr = ?", unitAccountNbr)
	}

	var records *[]CmbRecord
	err := tmp.
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error
	return records, err
}
