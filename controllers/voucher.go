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

// user_voucher_balance 조회
func GetVoucherBalances(c *gin.Context) {
	bals := make([]schema.VoucherBalanceRow, 0, 100)
	err := models.QueryVoucherBalances(&bals)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, bals)
}

// user_voucher_balance 생성
func CreateVoucherBalance(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbCreateVoucherBalance
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// data handling
	bal := schema.VoucherBalanceRow{
		AccountId:           req.AccountId,
		Addr:                req.Addr,
		PromotionId:         req.PromotionId,
		CurrentAmount:       req.CurrentAmount,
		TotalReceivedAmount: req.TotalReceivedAmount,
	}
	err = models.CreateVoucherBalance(&bal)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, bal)
	}
}

// 특정 user_voucher_balance 조회
func GetVoucherBalance(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	bal := schema.VoucherBalanceRow{
		Id: reqId,
	}
	err = models.QueryVoucherBalance(&bal)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, bal)
	}
}

// user_voucher_balance 정보 수정
func UpdateVoucherBalance(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Body request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbUpdateVoucherBalance
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	bal := schema.VoucherBalanceRow{
		Id:                  reqId,
		CurrentAmount:       req.CurrentAmount,
		TotalReceivedAmount: req.TotalReceivedAmount,
		UpdatedAt:           time.Now(),
	}
	err = models.UpdateVoucherBalance(&bal)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, bal)
	}
}

// user_voucher_balance 삭제
func DeleteVoucherBalance(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// handler data
	bal := schema.VoucherBalanceRow{
		Id: reqId,
	}
	err = models.DeleteVoucherBalance(&bal)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, bal)
	}
}

// voucher_send_history 전체 조회
func GetVoucherSendEvents(c *gin.Context) {

	events := make([]*types.ResGetVoucherSendEvents, 0, 100)
	err := models.QueryVoucherSendEvents(&events)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, events)
}

// voucher send, voucher burn 내역 조회
func GetTransferEvents(c *gin.Context) {
	addr := c.Param("addr")
	res  := types.ResTransfersHistoryByAddr{
		Addr: addr,
		VoucherSendEvents: make([]*types.ResGetVoucherSendEvents, 0, 100),
		VoucherBurnEvents: make([]*types.ResGetVoucherBurnEvents, 0, 100),
	}

	err := models.QueryVoucherSendEventsByAddr(&res.VoucherSendEvents, addr)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	err = models.QueryVoucherBurnEventsByAddr(&res.VoucherBurnEvents, addr)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, res)
}

// voucher_send_history 생성
func CreateTbVoucherSendEvent(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbCreateVoucherSendEvent
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// data handling
	event := schema.VoucherSendEventRow{
		AccountId:     req.AccountId,
		RecipientAddr: req.RecipientAddr,
		PromotionId:   req.PromotionId,
		Amount:        req.Amount,
		SentAt:        time.Now(),
	}
	err = models.CreateVoucherSendEvent(&event)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, event)
	}
}

// 특정 voucher_send_history 조회
func GetVoucherSendEvent(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	event := schema.VoucherSendEventRow{
		Id: reqId,
	}
	err = models.QueryVoucherSendEvent(&event)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, event)
	}
}

// voucher_send_history 정보 수정
func UpdateVoucherSendEvent(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Body request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbUpdateVoucherSendEvent
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	event := schema.VoucherSendEventRow{
		Id:     reqId,
		Amount: req.Amount,
	}
	err = models.UpdateVoucherSendEvent(&event)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, event)
	}
}

// voucher_send_history 삭제
func DeleteVoucherSendEvent(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// handler data
	event := schema.VoucherSendEventRow{
		Id: reqId,
	}
	err = models.DeleteVoucherSendEvent(&event)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, event)
	}
}

// voucher_burn_event 전체 조회
func GetVoucherBurnEvents(c *gin.Context) {
	events := make([]schema.VoucherBurnEventRow, 0, 100)
	err := models.QueryVoucherBurnEvents(&events)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, events)
}

// voucher_burn_event 생성
func CreateVoucherBurnEvent(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbCreateVoucherBurnEvent
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// data handling
	event := schema.VoucherBurnEventRow{
		AccountId:           req.AccountId,
		Addr:                req.Addr,
		PromotionId:         req.PromotionId,
		BurnedVoucherAmount: req.BurnedVoucherAmount,
		MintedTicketAmount:  req.MintedTicketAmount,
		BurnedAt:            time.Now(),
	}
	err = models.CreateVoucherBurnEvent(&event)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, event)
	}
}

// 특정 voucher_burn_event 조회
func GetVoucherBurnEvent(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	event := schema.VoucherBurnEventRow{
		Id: reqId,
	}
	err = models.QueryVoucherBurnEvent(&event)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, event)
	}
}

// voucher_burn_event 정보 수정
func UpdateVoucherBurnEvent(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Body request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbUpdateVoucherBurnEvent
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	event := schema.VoucherBurnEventRow{
		Id:                  reqId,
		BurnedVoucherAmount: req.BurnedVoucherAmount,
		MintedTicketAmount:  req.MintedTicketAmount,
	}
	err = models.UpdateVoucherBurnEvent(&event)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, event)
	}
}

// voucher_burn_event 삭제
func DeleteVoucherBurnEvent(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// handler data
	event := schema.VoucherBurnEventRow{
		Id: reqId,
	}
	err = models.DeleteVoucherBurnEvent(&event)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, event)
	}
}

// 현재 가용한 Voucher 목록 조회 ( Status : In progress, Not started )
func GetAvailableVouchers(c *gin.Context) {

	vouchers := make([]*types.ResGetAvailableVouchers, 0, 100)
	err := models.QueryAvailableVouchers(&vouchers)

	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, vouchers)
}

// voucher_send_history 생성
func CreateVoucherSendEvents(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqCreateVoucherSendEvent
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 프로모션 조회
	// convert req.PromotionId to uint64
	strPromotionId := fmt.Sprintf("%d", req.PromotionId)
	promotionId, err := strconv.ParseUint(strPromotionId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request Id path parameter", err)
		return
	}
	promotion := types.ResGetPromotion{
		PromotionId: promotionId,
	}
	err = models.QueryPromotion(&promotion)
	if err != nil {
		fmt.Printf("QueryPromotion 564 : %+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// 잔고 체크
	req.TotalAmount = uint64(len(req.RecipientAddrs) * int(req.Amount))
	if promotion.VoucherRemainingQty < req.TotalAmount {
		services.NotAcceptable(c, "remainingQty is smaller than totalSendAmount "+c.Request.Method+" "+c.Request.RequestURI, nil)
		return
	}
	// 프로모션 Remaining Qty 업데이트
	err = models.UpdatePromotion(&schema.PromotionRow{
		PromotionId:         req.PromotionId,
		VoucherRemainingQty: promotion.VoucherRemainingQty - req.Amount,
		UpdatedAt:           time.Now(),
	})
	if err != nil {
		fmt.Printf("UpdatePromotion 577 : %+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}


	for _, addr := range req.RecipientAddrs {

		fmt.Printf("addr : %s\n", addr)

		// 사용자 Account Id 조회
		
		addrType := services.GetAddressType(addr)
		acc := schema.AccountRow{
			Addr: addr,
			Type: addrType,
		}

		err := models.QueryOrCreateAccount(&acc)
		if err != nil {
			fmt.Printf("QueryAccount 485 : %+v\n", err.Error())
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
		fmt.Printf("acc : %+v\n", acc)
		// 바우처 전송 이벤트 생성
		err = models.CreateVoucherSendEvent(&schema.VoucherSendEventRow{
			AccountId:     acc.Id,
			RecipientAddr: addr,
			PromotionId:   req.PromotionId,
			Amount:        req.Amount,
			SentAt:        time.Now(),
		})

		if err != nil {
			fmt.Printf("CreateVoucherSendEvent 500 : %+v\n", err.Error())
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		// 사용자 바우처 잔고 조회
		voucherBalance := schema.VoucherBalanceRow{
			Addr:        addr,
			PromotionId: req.PromotionId,
		}

		err = models.QueryVoucherBalanceByAddrPromotionId(&voucherBalance)
		if err != nil && err.Error() != "record not found" {
			fmt.Printf("QueryVoucherBalanceByAddrPromotionId 515 : %+v\n", err.Error())
			voucherBalance.Id = 0
		}

		if voucherBalance.Id == 0 {
			// 사용자 바우처 잔고 레코드 생성
			err = models.CreateVoucherBalance(&schema.VoucherBalanceRow{
				AccountId:           acc.Id,
				Addr:                addr,
				PromotionId:         req.PromotionId,
				CurrentAmount:       req.Amount,
				TotalReceivedAmount: req.Amount,
				UpdatedAt:           time.Now(),
			})
			if err != nil {
				fmt.Printf("CreateVoucherBalance 529 : %+v\n", err.Error())
				services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				return
			}
		} else {
			// 사용자 바우처 잔고 업데이트
			err = models.UpdateVoucherBalance(&schema.VoucherBalanceRow{
				Id:                  voucherBalance.Id,
				CurrentAmount:       voucherBalance.CurrentAmount + req.Amount,
				TotalReceivedAmount: voucherBalance.TotalReceivedAmount + req.Amount,
				UpdatedAt:           time.Now(),
			})
			if err != nil {
				fmt.Printf("UpdateVoucherBalance 542 : %+v\n", err.Error())
				services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
				return
			}
		}
	}

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, nil)
	}
}
