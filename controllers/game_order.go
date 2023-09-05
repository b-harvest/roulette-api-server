package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Orders 조회
func GetGameOrders(c *gin.Context) {
	orders := make([]schema.OrderRow, 0, 100)
	err := models.QueryOrders(&orders)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, orders)
}


// Order 생성
func CreateGameOrder(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqTbCreateOrder
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// UsedTicketQty, IsWin 계산 필요
	// ClaimedAt. ClaimFinishedAt, 

	// data handling
	order := schema.OrderRow{
		AccountId:       req.AccountId,
		Addr:            req.Addr,
		PromotionId:     req.PromotionId,
		GameId:          req.GameId,
		Status:          req.Status,
		UsedTicketQty:   req.UsedTicketQty,
		PrizeId:         req.PrizeId,
		StartedAt:       time.Now(),
	}
	err = models.CreateOrder(&order)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, order)
	}
}

// 특정 Order 조회
func GetGameOrder(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("order_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	order := schema.OrderRow{
		OrderId: reqId,
	}
	err = models.QueryOrder(&order)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, order)
	}
}

// Order 정보 수정
func UpdateGameOrder(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("order_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
			services.BadRequest(c, "Bad Body request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
	}
	var req types.ReqTbUpdateOrder
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	order := schema.OrderRow{
		OrderId: reqId,
		IsWin: req.IsWin,
		Status: req.Status,
		UsedTicketQty: req.UsedTicketQty,
		ClaimedAt: req.ClaimedAt,
		ClaimFinishedAt: req.ClaimFinishedAt,
		UpdatedAt: time.Now(),
	}
	err = models.UpdateOrder(&order)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, order)
	}
}


// Order 삭제
func DeleteGameOrder(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("order_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// handler data
	order := schema.OrderRow{
		OrderId: reqId,
	}
	err = models.DeleteOrder(&order)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, order)
	}
}




