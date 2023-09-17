package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func PutAccount(c *gin.Context) {
	// type
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req schema.AccountRow
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	req.Addr         = c.Param("addr")
	req.LastLoginAt  = time.Now()
	req.TicketAmount = 0

	err = models.QueryOrCreateAccount(&req)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, &req)
}

func GetBalancesByAddr(c *gin.Context) {
	addr := c.Param("addr")
	bals := make([]types.ResGetBalanceByAcc, 0, 100)

	err := models.QueryBalancesByAddr(&bals, addr)
	if err != nil {
		if err.Error() != "record not found" {
			fmt.Printf("%+v\n", err.Error())
			services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}

	services.Success(c, nil, &bals)
}

func GetGameOrdersByAddr(c *gin.Context) {
	var orders []schema.OrderRow

	err := models.QueryOrdersByAcc(&orders, c.Param("address"))
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, &orders)
}

func GetWinTotalByAcc(c *gin.Context) {
	var resp []types.ResGetWinTotalByAcc

	err := models.QueryWinTotalByAcc(&resp, c.Param("addr"))
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, &resp)
}

// SELECT order_id, addr, is_win, promotion_id, claimed_at, claim_finished_at
// FROM game_order
// WHERE order_id = 9;
// if is_win == 1, claimed_at == null, claim_finished_at == null
// 	then (promotion_id 아래 sql 실행)
// 	else (errors.New("Can't claim (Not win or Already caimed)"))

// SELECT promotion_id, claim_start_at, claim_end_at
// FROM promotion
// WHERE promotion_id = 2;

// if now.After(b) && now.Before(a) || b == now || now == a {
// 	fmt.Println("true")
// }

// if claim_start_at <= now() <= claim_end_at
// 	then (아래 sql 실행)
// 	else (errors.New("Can't claim due to not claimable period"))

// UPDATE game_order
// SET claimed_at = now()
// WHERE order_id = ?

func PatchClaim(c *gin.Context) {
	orderId, err := strconv.ParseInt(c.Param("order-id"), 10, 64)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	order := schema.OrderRow{
		OrderId: orderId,
	}
	err = models.QueryOrderById(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Claimable only in below condition
	// status: before claim(3), win(true) and claimed(null), claim(null) finished not yet
	if order.Status != 3 && !order.IsWin || !order.ClaimedAt.IsZero() || !order.ClaimFinishedAt.IsZero() {
		err = errors.New("Can't claim due to not win or already claimed.")
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	promotion := schema.PromotionRow{
		PromotionId: order.PromotionId,
	}
	err = models.QueryPromotionById(&promotion)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Claimable only in below condition
	// claim-start <= now <= claim-end
	now := time.Now()
	if now.Before(promotion.ClaimStartAt) || now.After(promotion.ClaimEndAt) {
		err = errors.New("Can't claim due to not claimable period")
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	order.Status = 4
	order.ClaimedAt = now

	err = models.UpdateOrder(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, nil)
}

// accounts 조회
func GetAccounts(c *gin.Context) {
	accs := make([]schema.AccountRow, 0, 100)
	err := models.QueryAccounts(&accs)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, accs)
}

// account 생성
func CreateAccount(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbCreateAccount
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(string(jsonData))
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// data handling
	// acc := schema.AccountRow{
	// 	Addr:       req.Addr,
	// 	TicketAmount:      req.TicketAmount,
	// 	AdminMemo:     req.AdminMemo,
	// 	Type:             req.Type,
	// 	IsBlacklisted:      req.IsBlacklisted,
	// 	// LastLoginAt: time.Time{},
	// 	LastLoginAt: time.Now(),
	// }
	err = models.CreateAccount(&req)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, req)
	}
}

// 특정 Account 조회
func GetAccount(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("addr")
	acc := schema.AccountRow{
		Addr: strId,
	}
	err := models.QueryAccount(&acc)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, acc)
	}
}

// Account 정보 수정
func UpdateAccount(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("address")
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Body request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var req types.ReqTbUpdateAccount
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	acc := schema.AccountRow{
		Addr:          strId,
		TicketAmount:  req.TicketAmount,
		AdminMemo:     req.AdminMemo,
		Type:          req.Type,
		IsBlacklisted: req.IsBlacklisted,
		UpdatedAt:     time.Now(),
	}
	err = models.UpdateAccount(&acc)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		}
	} else {
		services.Success(c, nil, acc)
	}
}

// Account 삭제
func DeleteAccount(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("address")

	// handler data
	acc := schema.AccountRow{
		Addr: strId,
	}
	err := models.DeleteAccount(&acc)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, acc)
	}
}
