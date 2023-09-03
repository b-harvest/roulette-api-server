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

// Prizes 조회
func GetPrizes(c *gin.Context) {
	prizes := make([]schema.PrizeRow, 0, 100)
	err := models.QueryPrizes(&prizes)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, prizes)
}


// Prize 생성
func CreatePrize(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqCreatePrize
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// data handling
	prize := schema.PrizeRow{
		DistPoolId:       req.PromotionId,
		PromotionId:      req.PromotionId,
		PrizeDenomId:     req.PrizeDenomId,
		Amount:           req.Amount,
		Odds:             req.Odds,
		WinImageUrl:      req.WinImageUrl,
		MaxDailyWinLimit: req.MaxDailyWinLimit,
		MaxTotalWinLimit: req.MaxTotalWinLimit,
		IsActive:         req.IsActive,
	}
	err = models.CreatePrize(&prize)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, prize)
	}
}

// 특정 Prize 조회
func GetPrize(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("prize_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	prize := schema.PrizeRow{
		PrizeId: reqId,
	}
	err = models.QueryPrize(&prize)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, prize)
	}
}

// Prize 정보 수정
func UpdatePrize(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("prize_id")
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
	var req types.ReqUpdatePrize
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	prize := schema.PrizeRow{
		PrizeId: reqId,
		Odds: req.Odds,
		MaxDailyWinLimit: req.MaxDailyWinLimit,
		MaxTotalWinLimit: req.MaxTotalWinLimit,
		IsActive: req.IsActive,
		UpdatedAt: time.Now(),
	}
	err = models.UpdatePrize(&prize)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, prize)
	}
}


// DistPool 삭제
func DeletePrize(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("prize_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// handler data
	prize := schema.PrizeRow{
		PrizeId: reqId,
	}
	err = models.DeletePrize(&prize)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, prize)
	}
}




