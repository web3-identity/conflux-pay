package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/web3-identity/conflux-pay/models"
	pay_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/web3-identity/conflux-pay/utils"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

type AddUnitAccountReq struct {
	UnitAccountName string `json:"unit_account_name" binding:"required"`
	UnitAccountNbr  string `json:"unit_account_nbr" binding:"required"`
}

type SetUnitAccountRelationReq struct {
	UnitAccountNbr string `json:"unit_account_nbr" binding:"required"`
	BankAccountNbr string `json:"bank_account_nbr" binding:"required"`
}

func AddUnitAccount(c *gin.Context) {
	req := &AddUnitAccountReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	cmbClient := utils.GetDefaultCmbClient()
	_, err := cmbClient.AddUnitAccountV1Wrapper(req.UnitAccountName, req.UnitAccountNbr)
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_BUSINESS_COMMON)
		return
	}
	ginutils.RenderRespOK(c, nil)
}

func SetUnitAccountRelation(c *gin.Context) {
	req := &SetUnitAccountRelationReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	cmbClient := utils.GetDefaultCmbClient()
	_, err := cmbClient.SetUnitAccountRelationWrapper(req.UnitAccountNbr, req.BankAccountNbr)
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	ginutils.RenderRespOK(c, nil)
}

func QueryRecentCmbRecords(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	records, err := models.GetTodayAndYesterdayRecords("", "C", limit, offset)

	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	ginutils.RenderRespOK(c, records)
}

func QueryCmbRecords(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	records, err := models.GetCmbRecords(c.Query("unit_account_nbr"), c.Query("transaction_date"), c.Query("transaction_direction"), limit, offset)
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	ginutils.RenderRespOK(c, records)
}
