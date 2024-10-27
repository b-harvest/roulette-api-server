package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"roulette-api-server/config"
	"roulette-api-server/middlewares"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
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
	req.Addr = c.Param("addr")
	req.LastLoginAt = time.Now()
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
	var orders []*types.ResGetLatestOrderByAddr

	err := models.QueryOrdersByAddr(&orders, c.Param("addr"), c.Query("is-win"))
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, &orders)
}

// 유저 별 최근 주문 (game-id 필요)
func GetLatestOrder(c *gin.Context) {
	addr := c.Param("addr")
	strGameId := c.Query("game-id")
	if strGameId == "" {
		services.NotAcceptable(c, "invalid gameId "+c.Request.Method+" "+c.Request.RequestURI, nil)
		return
	}
	gameId, err := strconv.ParseInt(strGameId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "invalid gameId "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	order := types.ResGetLatestOrderByAddr{
		Addr:   addr,
		GameId: gameId,
	}

	err = models.QueryLatestOrderByAddr(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	if order.OrderId == 0 {
		services.Success(c, nil, nil)
	} else {
		services.Success(c, nil, &order)
	}
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

// This is claim function for BeraBola
func ClaimBola(c *gin.Context) {
	gameType := c.Param("type")
	addr := c.Param("addr")

	// 1. Check
	acc := schema.AccountRow{
		Addr: addr,
	}
	err := models.QueryAccountByAddr(&acc)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail query account "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	var ticketAmount *uint64
	var claimAmount int
	if gameType == "bola" {
		ticketAmount = &acc.TicketAmount
		claimAmount = 1
	} else if gameType == "gold" {
		ticketAmount = &acc.GoldTicketAmount
		claimAmount = 100
	} else {
		err = errors.New("invalid game type")
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	if *ticketAmount < 1 {
		err = errors.New("Can't claim due to not enough ticket amount")
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 2. Sending Token

	err = middlewares.SendToken(acc.Addr, claimAmount)
	if err != nil {
		services.NotAcceptable(c, "fail SendToken "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 3. Update Table

	// Table : account
	*ticketAmount = *ticketAmount - 1
	err = models.UpdateAccountTicketById(nil, &acc)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	services.Success(c, nil, acc)
}

// This is real claim function
func Claim(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var order schema.OrderRow
	if err = json.Unmarshal(jsonData, &order); err != nil {
		fmt.Println(string(jsonData))
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	err = models.QueryOrderByIdAndAddr(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail query order "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Claimable only in below condition
	// status: before claim(3), win(true) and claimed(null), claim(null) finished not yet
	if order.Status != 3 {
		err = errors.New("Can't claim due to not win or already claimed.")
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	promotion := schema.PromotionRowWithoutID{
		PromotionId: order.PromotionId,
	}
	err = models.QueryPromotionById(nil, &promotion)
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

	err = models.UpdateOrder(nil, &order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, nil)
}

func ClaimAll(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var order schema.OrderRow
	if err = json.Unmarshal(jsonData, &order); err != nil {
		fmt.Println(string(jsonData))
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	var orders []schema.OrderRow
	err = models.QueryJustOrdersByAddr(&orders, order.Addr)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail query order "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	var orderIds []string
	for _, order := range orders {
		// Claimable only in below condition
		// status: before claim(3), win(true) and claimed(null), claim(null) finished not yet
		if order.Status != 3 {
			continue
		}

		promotion := schema.PromotionRowWithoutID{
			PromotionId: order.PromotionId,
		}
		err = models.QueryPromotionById(nil, &promotion)
		if err != nil {
			continue
		}

		// Claimable only in below condition
		// claim-start <= now <= claim-end
		now := time.Now()
		if now.Before(promotion.ClaimStartAt) || now.After(promotion.ClaimEndAt) {
			continue
		}

		orderIds = append(orderIds, strconv.FormatInt(order.OrderId, 10))
	}

	err = models.UpdateOrdersByOrderIds(&orderIds)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, &types.ResAllClaim{
		Addr:            order.Addr,
		NumClaimedOrder: len(orderIds),
		Status:          4,
	})
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

// This function customized for-bb
func GetAccount(c *gin.Context) {
	addr := c.Param("addr")
	account := schema.AccountRow{
		Addr: addr,
	}
	accInfoRow := schema.AccountInfoRow{
		Addr: addr,
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)

			tx.Rollback()

			err := errors.New("panic")
			debug.PrintStack()
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}()
	err := tx.Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 1. Check logics

	// Check account is exist or not
	isAccExist, err := models.QueryAccountByAddrWithTx(tx, &account)
	if err != nil {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Check account info is exist or not
	isAccInfoExist, err := models.QueryAccountInfoWithTx(tx, &accInfoRow)
	if err != nil {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	chDelegated := make(chan *middlewares.IsSomethingReturnType, 1)
	chDelErr := make(chan error, 1)

	// TODO
	// Please check if you create or remove go function
	goFuncNum := 1

	wg := sync.WaitGroup{}
	wg.Add(goFuncNum)

	// Check delegation amount
	// if don't any deletegation, then return nil
	go func() {
		defer wg.Done()

		delegated, err := middlewares.IsDelegated(c.Param("addr"))
		if err != nil {
			chDelErr <- err
			return
		}
		chDelegated <- delegated
		chDelErr <- nil
	}()

	// 2. Create logics

	// If account or account info not exist create new one
	if !isAccExist {
		account.LastLoginAt = time.Now()
		account.TicketAmount = 0
		account.Type = "ETH"
		err = models.CreateAccountWithTx(tx, &account)
		if err != nil {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

			tx.Rollback()
			return
		}
	}
	if !isAccInfoExist {
		err = models.CreateAccountInfoWithTx(tx, &accInfoRow)
		if err != nil {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

			tx.Rollback()
			return
		}
	}

	// 3. Get goroutine results
	wg.Wait()

	err = <-chDelErr
	if err != nil {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}
	delegated := <-chDelegated

	// 4. Update logics

	delCondition := delegated != nil && delegated.Amount != accInfoRow.DelegationAmount
	// If amount increased, then increase ticket amount
	// 1000000000000000000 == 1BGT
	epsilon := math.Nextafter(1, 2) - 1
	goldBolaCondition := delCondition && ((delegated.Amount + epsilon) >= (accInfoRow.DelegationAmount + 199000000000000000000))
	bolaCondition := delCondition && ((delegated.Amount + epsilon) >= (accInfoRow.DelegationAmount + 1000000000000000000))
	if goldBolaCondition {
		account.GoldTicketAmount = account.GoldTicketAmount + 1
		accInfoRow.DelegationAmount = delegated.Amount
	} else if bolaCondition {
		account.TicketAmount = account.TicketAmount + 1
		accInfoRow.DelegationAmount = delegated.Amount
	}

	// Only update row If need to account and account_info
	if bolaCondition || goldBolaCondition {
		// Update account_info
		err = models.UpdateAccountInfoById(tx, &accInfoRow)
		if err != nil {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

			tx.Rollback()
			return
		}

		// Update account
		err = models.UpdateAccountById(tx, &account)
		if err != nil {
			services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

			tx.Rollback()
			return
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	acc := types.ResGetAccount{
		Account: account,
	}

	_, err = models.QueryAccountDetail(&acc)
	if err != nil {
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, acc)
}

// 특정 Account 조회
func GetAccountsDetail(c *gin.Context) {
	accounts := make([]types.ResGetAccount, 0, 500)
	err := models.QueryAccountsDetailPrepare(&accounts)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	for i := range accounts {
		models.QueryAccountDetail(&accounts[i])
	}

	// result
	if err != nil {
		services.NotAcceptable(c, "fail QueryAccountsDetail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
	} else {
		services.Success(c, nil, accounts)
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
