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

// 프로모션 조회
func GetPromotions(c *gin.Context) {
	promotions := make([]schema.PromotionRow, 0, 100)
	err := models.QueryTbPromotions(&promotions)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, promotions)
}


// 프로모션 생성
func CreatePromotion(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqCreatePromotion
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// data handling
	promotion := schema.PromotionRow{
		Title: req.Title,
		Desc: req.Desc,
		IsActive: req.IsActive,
		IsWhitelisted: req.IsWhitelisted,
		VoucherName: req.VoucherName,
		VoucherExchangeRatio0: req.VoucherExchangeRatio0,
		VoucherExchangeRatio1: req.VoucherExchangeRatio1,
		VoucherTotalSupply: req.VoucherTotalSupply,
		VoucherRemainingQty: req.VoucherTotalSupply,	// 초기값은 TotalSupply
		PromotionStartAt: req.PromotionStartAt,
		PromotionEndAt: req.PromotionEndAt,
		ClaimStartAt: req.ClaimStartAt,
		ClaimEndAt: req.ClaimEndAt,
	}
	err = models.CreatePromotion(&promotion)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, promotion)
	}
}

// 특정 프로모션 조회
func GetPromotion(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("promotion_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	promotion := schema.PromotionRow{
		PromotionId: reqId,
	}
	err = models.QueryTbPromotion(&promotion)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, promotion)
	}
}

// 프로모션 정보 수정
func UpdatePromotion(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("promotion_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request Id path parameter", err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
			services.BadRequest(c, "Bad Body Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
	}
	var req types.ReqUpdatePromotion
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// TotalSupply 가 변경된 경우 remainingQty 계산
	

	// handler data
	promotion := schema.PromotionRow{
		PromotionId: reqId,
		Title: req.Title,
		Desc: req.Desc,
		IsActive: req.IsActive,
		IsWhitelisted: req.IsWhitelisted,
		VoucherName: req.VoucherName,
		VoucherExchangeRatio0: req.VoucherExchangeRatio0,
		VoucherExchangeRatio1: req.VoucherExchangeRatio1,
		VoucherTotalSupply: req.VoucherTotalSupply,
		VoucherRemainingQty: req.VoucherRemainingQty,
		PromotionStartAt: req.PromotionStartAt,
		PromotionEndAt: req.PromotionEndAt,
		ClaimStartAt: req.ClaimStartAt,
		ClaimEndAt: req.ClaimEndAt,
		UpdatedAt: time.Now(),
	}
	err = models.UpdatePromotion(&promotion)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, promotion)
	}
}


// 프로모션 정보 삭제
func DeletePromotion(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("promotion_id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// handler data
	promotion := schema.PromotionRow{
		PromotionId: reqId,
	}
	err = models.DeletePromotion(&promotion)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, promotion)
	}
}