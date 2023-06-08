package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/web3-identity/conflux-pay/models"
	pay_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/web3-identity/conflux-pay/utils"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

type AddUnitAccountReq struct {
	UnitAccountName string `json:"unit_account_name" binding:"required"`
	UnitAccountNbr  string `json:"unit_account_nbr" binding:"required"`
}

type AddUnitAccountResult struct {
	FullUnitAccountNbr string `json:"full_unit_account_nbr" binding:"required"`
}

type SetUnitAccountRelationReq struct {
	UnitAccountNbr string `json:"unit_account_nbr" binding:"required"`
	BankAccountNbr string `json:"bank_account_nbr" binding:"required"`
}

type SetUnitAccountRelationResult struct {
	Code int `json:"code" binding:"required"`
}

//	@Tags			Cmb
//	@ID				AddUnitAccount
//	@Summary		Add a unit account
//	@Description	Add a unit account
//	@Produce		json
//	@Param			add_unit_account_req	body	controllers.AddUnitAccountReq	true	"add_unit_account_req"
//	@Success		200 {object}    AddUnitAccountResult
//	@Failure		400	{object}	cns_errors.RainbowErrorDetailInfo	"Invalid request"
//	@Failure		500	{object}	cns_errors.RainbowErrorDetailInfo	"Internal Server error"
//	@Router			/cmb/unit-account [post]
func AddUnitAccount(c *gin.Context) {
	req := &AddUnitAccountReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	cmbClient := utils.GetDefaultCmbClient()
	_, err := cmbClient.AddUnitAccountV1Wrapper(req.UnitAccountName, req.UnitAccountNbr)
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INTERNAL_SERVER_COMMON)
		return
	}
	rtn := &AddUnitAccountResult{
		FullUnitAccountNbr: (viper.GetString("apps.cmb.accNbr") + req.UnitAccountNbr),
	}
	ginutils.RenderRespOK(c, rtn)
}

//	@Tags			Cmb
//	@ID				SetUnitAccountRelation
//	@Summary		Set a related bank account of a unit account
//	@Description	Set a related bank account of a unit account
//	@Produce		json
//	@Param			set_unit_account_relation_req	body	controllers.SetUnitAccountRelationReq	true	"set_unit_account_relation_req"
//	@Success		200 {object}    SetUnitAccountRelationResult
//	@Failure		400	{object}	cns_errors.RainbowErrorDetailInfo	"Invalid request"
//	@Failure		500	{object}	cns_errors.RainbowErrorDetailInfo	"Internal Server error"
//	@Router			/cmb/unit-account/relation [post]
func SetUnitAccountRelation(c *gin.Context) {
	req := &SetUnitAccountRelationReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	cmbClient := utils.GetDefaultCmbClient()
	_, err := cmbClient.SetUnitAccountRelationWrapper(req.UnitAccountNbr, req.BankAccountNbr)
	if err != nil {
		ginutils.RenderRespError(c, err, pay_errors.ERR_INTERNAL_SERVER_COMMON)
		return
	}
	ginutils.RenderRespOK(c, &SetUnitAccountRelationResult{
		Code: 200,
	})
}

//	@Tags			Cmb
//	@ID				QueryRecentCmbRecords
//	@Summary		查询昨天和今天汇入的交易
//	@Description	查询昨天和今天汇入的交易
//	@Produce		json
//	@Param			limit	query		int	true	"limit"
//	@Param			offset	query		int	true	"offset"
//	@Success		200		{array}		models.CmbRecord
//	@Failure		400		{object}	cns_errors.RainbowErrorDetailInfo	"Invalid request"
//	@Failure		500		{object}	cns_errors.RainbowErrorDetailInfo	"Internal Server error"
//	@Router			/cmb/history/recent [get]
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
		ginutils.RenderRespError(c, err, pay_errors.ERR_INTERNAL_SERVER_COMMON)
		return
	}
	ginutils.RenderRespOK(c, records)
}

//	@Tags			Cmb
//	@ID				QueryHistoryCmbRecords
//	@Summary		查询历史交易
//	@Description	查询历史交易
//	@Produce		json
//	@Param			limit					query		int		true	"limit"
//	@Param			offset					query		int		true	"offset"
//	@Param			unit_account_nbr		query		string	true	"specified unit account number"
//	@Param			transaction_date		query		string	true	"specified date, e.g. 20230523"
//	@Param			transaction_direction	query		string	true	"transaction direction, C for recieve and D for out"
//	@Success		200						{array}		models.CmbRecord
//	@Failure		400						{object}	cns_errors.RainbowErrorDetailInfo	"Invalid request"
//	@Failure		500						{object}	cns_errors.RainbowErrorDetailInfo	"Internal Server error"
//	@Router			/cmb/history [get]
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
		ginutils.RenderRespError(c, err, pay_errors.ERR_INTERNAL_SERVER_COMMON)
		return
	}
	ginutils.RenderRespOK(c, records)
}
