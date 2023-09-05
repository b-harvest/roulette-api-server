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

// TODO
/*
	- /promotions/live	(프론트에서 사용할 프로모션 정보)
	- 유저용/어드민용 따로 /promotions 분리
*/

/*
	1. promotion 테이블
		- not started / in progress 여부
		- 참여자 수
	2. 프로모션에 속하는 distribution_pool 리스트
	- 풀의 prize_denom 정보
	- 풀에 속하는 prize 리스트
	# query
		- by 프로모션 title
		- by 진행 중인지(기간)
		- by
	- order by promotion_start_at desc
*/
// 프로모션 조회
func GetPromotions(c *gin.Context) {
	promotions := make([]*types.ResGetPromotions, 0, 100)
	err := models.QueryPromotions(&promotions)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// 각 프로모션의 풀 조회
	for _, v := range promotions {
		v.DistributionPools, err = models.QueryDistPoolsByPromId(v.PromotionId)
		if err != nil {
			fmt.Println(err)
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
		}
	}

	services.Success(c, nil, promotions)
}

// 특정 프로모션 조회
func GetPromotion(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("promotion_id")
	reqId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	p := types.ResGetPromotion{
		PromotionId: reqId,
	}

	// 프로모션 조회
	err = models.QueryPromotion(&p)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// 프로모션 Summary 조회
	pSummary, err := models.QueryPromotionSummary(reqId)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	p.Summary = pSummary

	p.DistributionPools, err = models.QueryDistPoolsDetailByPromId(reqId)
	if err != nil {
		fmt.Println(err)
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}	

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, p)
	}
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