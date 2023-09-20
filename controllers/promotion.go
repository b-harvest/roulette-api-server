package controllers

import (
	"encoding/json"
	"errors"
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
	- optional: /promotions/live	(프론트에서 사용할 프로모션 정보)
	- optional: 유저용/어드민용 따로 /promotions 분리
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
	- order by promotion_start_at desc
*/

func validateQuery(c *gin.Context, key string, col string, qMap types.QueryFilterMap) (rMap types.QueryFilterMap, err error) {
	if c.Query(key) != "" {
		qMap[col] = c.Query(key)
	}
	return qMap, nil
}

// 프로모션 조회
func GetPromotions(c *gin.Context) {
	var err error
	// Query filter 조회
	// promotion-id, title, url, status, is-whitelisted, is-active
	// status: "not-started", "in-progress", "finished"
	qMap := make(types.QueryFilterMap, 100)
	qMap, err = validateQuery(c, "promotion-id", "promotion_id", qMap)
	qMap, err = validateQuery(c, "title", "title", qMap)
	qMap, err = validateQuery(c, "url", "url", qMap)
	qMap, err = validateQuery(c, "is-whitelisted", "is_whitelisted", qMap)
	qMap, err = validateQuery(c, "is-active", "is_active", qMap)
	qMap, err = validateQuery(c, "status", "status", qMap)
	
	promotions := make([]*types.ResGetPromotions, 0, 100)
	promotions, err = models.QueryPromotions(&promotions, qMap)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 각 프로모션의 풀 조회
	for _, v := range promotions {
		v.DistributionPools, err = models.QueryDistPoolsByPromId(v.PromotionId)
		if err != nil {
			fmt.Println(err)
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
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

	// 프로모션 Summary 조회
	pSummary, err := models.QueryPromotionSummary(reqId)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	p.Summary = pSummary

	p.DistributionPools, err = models.QueryDistPoolsDetailByPromId(reqId)
	if err != nil {
		fmt.Println(err)
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, p)
	}
}

// 프로모션 생성
func CreatePromotion(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqCreatePromotion
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// create transaction
	tx, err := models.CreateTxInstance()
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "tx error : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// create promotion
	promotion := schema.PromotionRow{
		Title:                 req.Title,
		Desc:                  req.Desc,
		Url:                   req.Url,
		IsActive:              req.IsActive,
		IsWhitelisted:         req.IsWhitelisted,
		VoucherName:           req.VoucherName,
		VoucherExchangeRatio0: req.VoucherExchangeRatio0,
		VoucherExchangeRatio1: req.VoucherExchangeRatio1,
		VoucherTotalSupply:    req.VoucherTotalSupply,
		VoucherRemainingQty:   req.VoucherTotalSupply, // 초기값은 TotalSupply
		PromotionStartAt:      req.PromotionStartAt,
		PromotionEndAt:        req.PromotionEndAt,
		ClaimStartAt:          req.ClaimStartAt,
		ClaimEndAt:            req.ClaimEndAt,
	}
	err = models.CreatePromotionWithTx(tx, &promotion)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
			return
		} else {
			services.NotAcceptable(c, "fail CreatePromotionWithTx "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}

	// create dist pools
	for _, v := range req.DistributionPools {
		creatingPool := schema.PrizeDistPoolInsertRow{
			PromotionId:  promotion.PromotionId,
			PrizeDenomId: v.PrizeDenomId,
			TotalSupply:  v.TotalSupply,
			RemainingQty:  v.TotalSupply,	// default
			IsActive:     true,	// default
		}
		err = models.CreateDistPoolWithTx(tx, &creatingPool)
	
		// result
		if err != nil {
			tx.Rollback()
			fmt.Printf("%+v\n",err.Error())
			if strings.Contains(err.Error(),"1062") {
				services.NotAcceptable(c, "data already exists", err)
				return
			} else {
				services.NotAcceptable(c, "fail CreateDistPoolWithTx " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
				return
			}
		}

		// fmt.Println(creatingPool.ID)
		// 생성한 pool 조회
		pool := creatingPool

		// create prizes
		for _, reqPrize := range v.Prizes {
			// prize 생성
			prize := schema.PrizeRow{
				DistPoolId:       pool.ID,
				PromotionId:      pool.PromotionId,
				PrizeDenomId:     pool.PrizeDenomId,
				Amount:           reqPrize.Amount,
				Odds:             reqPrize.Odds,
				WinImageUrl:      reqPrize.WinImageUrl,
				MaxDailyWinLimit: reqPrize.MaxDailyWinLimit,
				MaxTotalWinLimit: reqPrize.MaxTotalWinLimit,
				IsActive:         true, //default
			}
			err = models.CreatePrizeWithTx(tx, &prize)

			// result
			if err != nil {
				tx.Rollback()
				fmt.Printf("%+v\n",err.Error())
				if strings.Contains(err.Error(),"1062") {
					services.NotAcceptable(c, "data already exists", err)
					return
				} else {
					services.NotAcceptable(c, "fail CreatePrizeWithTx " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
					return
				}
			}
		}
	}

	err = tx.Commit().Error
	if err != nil {
		services.NotAcceptable(c, "commit failed", err)
		return
	}
	services.Success(c, nil, promotion)	
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
		services.BadRequest(c, "Bad Body Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqUpdatePromotion
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 현재 프로모션 정보 조회
	p := types.ResGetPromotion{
		PromotionId: uint64(reqId),
	}
	err = models.QueryPromotion(&p)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail QueryPromotion"+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// TotalSupply 가 변경된 경우 v 계산
	// TotalSupply 는 사용량(cur_TotalSupply - cur_remain) 보다 작을 수 없다.
	// 이미 사용한 양이 공급량(TotalSupply) 보다 클 수 없기 때문이다.
	var remainingQty uint64
	if p.VoucherTotalSupply != req.VoucherTotalSupply {
		usedQty := p.VoucherTotalSupply - p.VoucherRemainingQty
		if req.VoucherTotalSupply < usedQty {
			err := errors.New("total supply can not be smaller than used #")
			fmt.Printf("%+v\n", err.Error())
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
		remainingQty = req.VoucherTotalSupply - usedQty
	}

	// create transaction
	tx, err := models.CreateTxInstance()
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "tx error : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// handler data
	promotion := schema.PromotionRow{
		PromotionId:           reqId,
		Title:                 req.Title,
		Desc:                  req.Desc,
		Url:                   req.Url,
		IsActive:              req.IsActive,
		IsWhitelisted:         req.IsWhitelisted,
		VoucherName:           req.VoucherName,
		VoucherExchangeRatio0: req.VoucherExchangeRatio0,
		VoucherExchangeRatio1: req.VoucherExchangeRatio1,
		VoucherTotalSupply:    req.VoucherTotalSupply,
		VoucherRemainingQty:   remainingQty,
		PromotionStartAt:      req.PromotionStartAt,
		PromotionEndAt:        req.PromotionEndAt,
		ClaimStartAt:          req.ClaimStartAt,
		ClaimEndAt:            req.ClaimEndAt,
		UpdatedAt:             time.Now(),
	}
	err = models.UpdatePromotionWithTx(tx, &promotion)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
			return
		} else {
			services.NotAcceptable(c, "fail UpdatePromotionWithTx "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}

	// update or insert or delete dist pools
	for _, dPool := range req.DistributionPools {
		switch dPool.UpdateType {
		case "":	// normal update
			// 현재 dPool 정보 조회
			dp := schema.PrizeDistPoolRow{
				DistPoolId: dPool.DistPoolId,
			}
			err = models.QueryDistPool(&dp)
			if err != nil {
				fmt.Printf("%+v\n", err.Error())
				tx.Rollback()
				services.NotAcceptable(c, "fail QueryDistPool"+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				return
			}

			// remainingQty 조회 및 계산
			var rQty uint64
			if dPool.TotalSupply != dp.TotalSupply {
				usedQty := dp.TotalSupply - dp.RemainingQty
				if dPool.TotalSupply < usedQty {
					err := errors.New("dist pool total_supply can not be smaller than used #")
					fmt.Printf("%+v\n", err.Error())
					tx.Rollback()
					services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
					return
				}
				rQty = dPool.TotalSupply - usedQty
				fmt.Printf("요청totalSupply: %+v, 사용한Qty: %+v, 변경되는 remainQty: %+v", dPool.TotalSupply, usedQty, rQty)
			}

			// handler data
			pool := schema.PrizeDistPoolRow{
				DistPoolId:   dPool.DistPoolId,
				TotalSupply:  dPool.TotalSupply,
				RemainingQty: rQty,
				IsActive:     dPool.IsActive,
				UpdatedAt:    time.Now(),
			}
			err = models.UpdateDistPoolWithTx(tx, &pool)
			// result
			if err != nil {
				tx.Rollback()
				fmt.Printf("%+v\n", err.Error())
				if strings.Contains(err.Error(), "1062") {
					services.NotAcceptable(c, "something duplicated. already exists. fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				} else {
					services.NotAcceptable(c, "fail UpdateDistPoolWithTx"+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				}
				return
			}

			// update prizes
			for _, pr := range dPool.Prizes {
				switch pr.UpdateType {
				case "":	// normal update
					// handler data
					prize := schema.PrizeRow{
						PrizeId:          pr.PrizeId,
						Odds:             pr.Odds,
						MaxDailyWinLimit: pr.MaxDailyWinLimit,
						MaxTotalWinLimit: pr.MaxTotalWinLimit,
						IsActive:         pr.IsActive,
						UpdatedAt:        time.Now(),
					}
					err = models.UpdatePrize(&prize)

					// result
					if err != nil {
						tx.Rollback()
						fmt.Printf("%+v\n",err.Error())
						if strings.Contains(err.Error(),"1062") {
							services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
						} else {
							services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
						}
						return
					}
					
				case "insert":
					// prize 생성
					prize := schema.PrizeRow{
						DistPoolId:       dPool.DistPoolId,
						PromotionId:      reqId,
						PrizeDenomId:     dPool.PrizeDenomId,
						Amount:           pr.Amount,
						Odds:             pr.Odds,
						WinImageUrl:      pr.WinImageUrl,
						MaxDailyWinLimit: pr.MaxDailyWinLimit,
						MaxTotalWinLimit: pr.MaxTotalWinLimit,
						IsActive:         true, //default
					}
					err = models.CreatePrizeWithTx(tx, &prize)

					// result
					if err != nil {
						tx.Rollback()
						fmt.Printf("%+v\n",err.Error())
						if strings.Contains(err.Error(),"1062") {
							services.NotAcceptable(c, "data already exists", err)
							return
						} else {
							services.NotAcceptable(c, "fail CreatePrizeWithTx " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
							return
						}
					}
				case "delete":
					if p.Status == "not-started" {
						// handler data
						prize := schema.PrizeRow{
							PrizeId: pr.PrizeId,
						}
						err = models.DeletePrize(&prize)
						// result
						if err != nil {
							tx.Rollback()
							services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
							return
						} 
					}
				}
			}
		case "insert":
			creatingPool := schema.PrizeDistPoolInsertRow{
				PromotionId:  reqId,
				PrizeDenomId: dPool.PrizeDenomId,
				TotalSupply:  dPool.TotalSupply,
				RemainingQty: dPool.TotalSupply,	// default
				IsActive:     true,	// default
			}
			err = models.CreateDistPoolWithTx(tx, &creatingPool)
		
			// result
			if err != nil {
				tx.Rollback()
				fmt.Printf("%+v\n",err.Error())
				if strings.Contains(err.Error(),"1062") {
					services.NotAcceptable(c, "data already exists", err)
					return
				} else {
					services.NotAcceptable(c, "fail CreateDistPoolWithTx " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
					return
				}
			}
	
			// fmt.Println(creatingPool.ID)
			// 생성한 pool 조회
			pool := creatingPool
	
			// create prizes
			for _, reqPrize := range dPool.Prizes {
				// prize 생성
				prize := schema.PrizeRow{
					DistPoolId:       pool.ID,
					PromotionId:      pool.PromotionId,
					PrizeDenomId:     pool.PrizeDenomId,
					Amount:           reqPrize.Amount,
					Odds:             reqPrize.Odds,
					WinImageUrl:      reqPrize.WinImageUrl,
					MaxDailyWinLimit: reqPrize.MaxDailyWinLimit,
					MaxTotalWinLimit: reqPrize.MaxTotalWinLimit,
					IsActive:         true, //default
				}
				err = models.CreatePrizeWithTx(tx, &prize)
	
				// result
				if err != nil {
					tx.Rollback()
					fmt.Printf("%+v\n",err.Error())
					if strings.Contains(err.Error(),"1062") {
						services.NotAcceptable(c, "data already exists", err)
						return
					} else {
						services.NotAcceptable(c, "fail CreatePrizeWithTx " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
						return
					}
				}
			}
		case "delete":
			if p.Status != "not-started" {
				tx.Rollback()
				services.NotAcceptable(c, "cat not delete dPool because promotion already started " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), errors.New("started"))
				return
			}

			// handler data
			pool := schema.PrizeDistPoolRow{
				DistPoolId: dPool.DistPoolId,
			}
			err = models.DeleteDistPool(&pool)
			// result
			if err != nil {
				tx.Rollback()
				services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				return
			}

			for _, pr := range dPool.Prizes {
				// handler data
				prize := schema.PrizeRow{
					PrizeId: pr.PrizeId,
				}
				err = models.DeletePrize(&prize)
				// result
				if err != nil {
					tx.Rollback()
					services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
					return
				} 
			}
		}
	}
		
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		services.NotAcceptable(c, "commit failed", err)
		return
	}
	services.Success(c, nil, nil)
}

