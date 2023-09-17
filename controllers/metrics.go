package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"

	"github.com/gin-gonic/gin"
)

// ============ For event_wallet_conn ===========

func GetEventWalletConn(c *gin.Context) {
	promotionId := c.Query("promotion-id")
	addr := c.Query("addr")
	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	var events []schema.EventwalletConnRow
	err := models.QueryEventWalletConn(&events, promotionId, addr, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, &events)
}

func GetEventWalletConnCount(c *gin.Context) {
	promotionId := c.Query("promotion-id")
	addr := c.Query("addr")
	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	var resp types.ResGetEventCount
	err := models.QueryEventWalletConnCount(&resp, promotionId, addr, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, &resp)
}

func PostEventWalletConn(c *gin.Context) {
	jsonBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqPostEvent
	err = json.Unmarshal(jsonBytes, &req)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	err = models.CreateEventWalletConn(&req)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, nil)
}

// ============ For event_flip_link ===========

func GetEventFlipLink(c *gin.Context) {
	promotionId := c.Query("promotion-id")
	addr := c.Query("addr")
	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	var events []schema.EventFlipLinkRow
	err := models.QueryEventFlipLink(&events, promotionId, addr, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, &events)
}

func GetEventFlipLinkCount(c *gin.Context) {
	promotionId := c.Query("promotion-id")
	addr := c.Query("addr")
	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	var resp types.ResGetEventCount
	err := models.QueryEventFlipLinkCount(&resp, promotionId, addr, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, &resp)
}

func PostEventFlipLinks(c *gin.Context) {
	jsonBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqPostEvent
	err = json.Unmarshal(jsonBytes, &req)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	err = models.CreateEventFlipLink(&req)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, nil)
}
