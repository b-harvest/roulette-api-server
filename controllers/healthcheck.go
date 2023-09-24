package controllers

import (
	// "encoding/json"

	// "io"

	"roulette-api-server/services"
	"roulette-api-server/types"

	"github.com/gin-gonic/gin"
)

func GetHealthcheck(c *gin.Context) {
	response := types.ResHealthcheck{
		Status: "OK",
	}
	services.Success(c, nil, response)
}
