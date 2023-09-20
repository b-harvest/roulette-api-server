package controllers

import (
	// "encoding/json"
	"fmt"
	"strconv"

	// "io"
	"roulette-api-server/models"
	"roulette-api-server/services"
	"roulette-api-server/types"

	"github.com/gin-gonic/gin"
)

// /stats/accounts[?start-date=2023-09-01&end-date=2023-09-30]
// 전체 계정 수, 블랙리스트 수
// [특정 기간에] 등록한 사용자 수, 로그인한 사용자 수
func GetAccountStat(c *gin.Context) {
	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	stat, err := models.QueryAccountStat(startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

// /stats/promotions[?start-date=2023-09-01&end-date=2023-09-30]
// 진행중, 종료, 미시작 프로모션 수
func GetPromotionStat(c *gin.Context) {

	stat, err := models.QueryPromotionStat()

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

func GetPromotionSummary(c *gin.Context) {

	strId := c.Param("promotion_id")
	reqId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 프로모션 Summary 조회
	pSummary, err := models.QueryPromotionSummary(reqId)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, pSummary)
}

// /stats/flip-links/{promotion_id}[?start-date=2023-09-01&end-date=2023-09-30]
func GetFlipLinkStat(c *gin.Context) {

	strId := c.Param("promotion_id")

	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	if startDate == "" || endDate == "" {
		reqId, err := strconv.ParseUint(strId, 10, 64)
		if err != nil {
			services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		p := types.ResGetPromotion{
			PromotionId: reqId,
		}

		// 프로모션 조회
		err = models.QueryPromotion(&p)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		startDate = p.PromotionStartAt.Format("2006-01-02")
		endDate = p.PromotionEndAt.Format("2006-01-02")
	}

	stat, err := models.QueryFlipLinkStat(strId, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

// /stats/wallet-connects/{promotion_id}[?start-date=2023-09-01&end-date=2023-09-30]
func GetWalletConnectStat(c *gin.Context) {

	strId := c.Param("promotion_id")

	startDate := c.Query("start-date")
	endDate := c.Query("end-date")

	if startDate == "" || endDate == "" {
		reqId, err := strconv.ParseUint(strId, 10, 64)
		if err != nil {
			services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		p := types.ResGetPromotion{
			PromotionId: reqId,
		}

		// 프로모션 조회
		err = models.QueryPromotion(&p)
		if err != nil {
			fmt.Printf("%+v\n", err.Error())
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		startDate = p.PromotionStartAt.Format("2006-01-02")
		endDate = p.PromotionEndAt.Format("2006-01-02")
	}

	stat, err := models.QueryWalletConnectStat(strId, startDate, endDate)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

// /stats/vouchers/{promotion_id}
func GetVoucherStat(c *gin.Context) {
	strId := c.Param("promotion_id")

	stat, err := models.QueryVoucherStat(strId)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

// /stats/tickets/{promotion_id}
func GetTicketStat(c *gin.Context) {
	strId := c.Param("promotion_id")

	stat, err := models.QueryTicketStat(strId)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}

// /stats/prizes/{promotion_id}
func GetPrizeStat(c *gin.Context) {
	strId := c.Param("promotion_id")

	stat, err := models.QueryPrizeStat(strId)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, stat)
}
