package controllers

import (
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"

	"github.com/gin-gonic/gin"
)

func GetBalanceByAddr(c *gin.Context) {
	iAddr := c.Param("addr")

	var account schema.Account
	err := models.QueryBalanceByAddr(&account, iAddr)
	if err != nil {
		if err.Error() == "record not found" {	// if addr not exists
			err = models.InsertNewAddr(&account, iAddr)
			if err != nil {
				services.NotAcceptable(c, "Can not enroll user address", err)
			}
			services.Success(c, nil, account)
		} else {
			services.NotAcceptable(c, "Something went wrong! Can not query user balance", err)
		}
	} else {
		services.Success(c, nil, account)
	}
}
