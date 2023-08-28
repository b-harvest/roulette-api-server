package controllers

import (
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SwapVoucherToTicket(c *gin.Context) {
	iAddr := c.Param("addr")
	numOfVoucher := c.Param("voucher_num")
	iNumOfVoucher, err := strconv.ParseInt(numOfVoucher, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "num of voucher is not an integer", err)
		return
	}

	var account schema.Account
	err = models.Swap(&account, iAddr, iNumOfVoucher)
	if err != nil {
		services.NotAcceptable(c, "Can not swap", err)
	} else {
		services.Success(c, nil, account)
	}
}


func SendVoucher(c *gin.Context) {
	iAddr := c.Param("addr")
	numOfVoucher := c.Param("voucher_num")
	iNumOfVoucher, err := strconv.ParseInt(numOfVoucher, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "num of voucher is not an integer", err)
		return
	}

	// check whether account exists
	var account schema.Account
	err = models.QueryBalanceByAddr(&account, iAddr)
	if err != nil {
		if err.Error() == "record not found" {	// if addr not exists
			err = models.InsertNewAddr(&account, iAddr)
			if err != nil {
				services.NotAcceptable(c, "Can not enroll user address", err)
			}
		} else {
			services.NotAcceptable(c, "Something went wrong! Can not query user balance", err)
		}
	}

	// send voucher
	err = models.SendVoucher(&account, iAddr, iNumOfVoucher)
	if err != nil {
		services.NotAcceptable(c, "Can not send voucher", err)
	} else {
		services.Success(c, nil, account)
	}
}