package controllers

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"roulette-api-server/config"
	"roulette-api-server/middlewares"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func StartGame(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var order schema.OrderRow
	// Update game order
	if err = json.Unmarshal(jsonData, &order); err != nil {
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)

			tx.Rollback()

			err = errors.New("panic")
			debug.PrintStack()
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}()
	err = tx.Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 1. Check

	// Table : account
	account := schema.AccountRow{
		Addr: order.Addr,
	}
	// Check whether account is doing game
	isLocked, err := models.QueryAndLockAccountByAddr(tx, &account)
	if isLocked {
		// If is there any in progress game
		// then just ignore request
		services.Success(c, "there are some games in progress", nil)

		tx.Rollback()
		return
	}
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	var ticketAmount *uint64
	if order.PromotionId == 1 && order.GameId == 1 {
		// Berabola game
		ticketAmount = &account.TicketAmount
	} else if order.PromotionId == 2 && order.GameId == 2 {
		// Gold Berabola game
		ticketAmount = &account.GoldTicketAmount
	} else {
		err = errors.New("invalid game type")
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}

	// TODO: calculate required ticket amount for game start
	// Check whether account has sufficient ticket amount
	ticketQtyForGame := uint64(1)
	if *ticketAmount < ticketQtyForGame {
		err = errors.New("insufficient ticket amount")
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}

	// Update game order
	order.AccountId = account.Id
	order.UsedTicketQty = ticketQtyForGame

	// Table : game_order
	counter, err := models.QueryInProgressGameCnt(tx, &order)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// Check existence of in progress game
	if counter.Cnt != 0 {
		// If is there any in progress game
		// then just ignore request
		services.Success(c, "there are some games in progress", nil)

		tx.Rollback()
		return
	}

	// Table : promotion
	promotion := schema.PromotionRowWithoutID{
		PromotionId: order.PromotionId,
	}
	err = models.QueryPromotionById(tx, &promotion)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// Check whether active promotion
	if !promotion.IsActive {
		err = errors.New("promotion isn't active")
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}
	// Check whether in promotion periods
	now := time.Now()
	if now.Before(promotion.PromotionStartAt) || now.After(promotion.PromotionEndAt) {
		err = errors.New("not in promotion periods")
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}

	// Table : game_type
	gameType := schema.Game{
		GameId: order.GameId,
	}
	err = models.QueryGameType(tx, &gameType)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// Check whether active game
	if !gameType.IsActive {
		err = errors.New("game isn't active")
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}

	// 2. Game Logic

	// Update game order
	order.StartedAt = time.Now()

	// Query joined data(prize info) with prize, prize_denom, distribution_pool
	var tmpPrizeInfos []types.PrizeInfo
	err = models.QueryPrizeInfosByPromotionId(tx, &tmpPrizeInfos, promotion.PromotionId)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Generate random number 0 ~ 99999
	// n, err := rand.Int(rand.Reader, big.NewInt(100000))
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	randNum := n.Uint64()

	var resPrizeInfo *types.PrizeInfo

	// Check whether is active prize
	var prevAccumOdds uint64 = 0
	for _, prizeInfo := range tmpPrizeInfos {
		// Query the number of today win game_orders of prize
		var todayWinCounter types.Counter
		err = models.QueryTodayWins(tx, &todayWinCounter, prizeInfo.PrizeId)
		if err != nil {
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		// Only include prizes that meet the conditions
		condition1 := !prizeInfo.PIsActive || !prizeInfo.PDIsActive || !prizeInfo.DPIsActive
		condition2 := prizeInfo.MaxTotalWinLimit < (prizeInfo.WinCnt + 1)
		condition3 := prizeInfo.MaxDailyWinLimit < (todayWinCounter.Cnt + 1)
		condition4 := prizeInfo.RemainingQty < prizeInfo.Amount
		if condition1 || condition2 || condition3 || condition4 {
			continue
		}

		// prevAccumOdds <= randNum < currentAccumOdds
		odds := uint64(prizeInfo.Odds * float64(1000))
		currentAccumOdds := prevAccumOdds + odds

		// For test
		// msg := fmt.Sprintf("id %d, odds %d : ", prizeInfo.PrizeId, odds)
		// fmt.Println(msg, prevAccumOdds, "<=", randNum, "<", currentAccumOdds)

		if prevAccumOdds <= randNum && randNum < currentAccumOdds {
			resPrizeInfo = &prizeInfo
			break
		}
		prevAccumOdds = currentAccumOdds
	}

	// Update game order
	if resPrizeInfo == nil {
		// Lost
		order.PrizeId = 0
		order.IsWin = false
	} else {
		// Win
		order.PrizeId = resPrizeInfo.PrizeId
		order.IsWin = true
	}
	order.Status = 1

	// 3. Update

	// Table : prize, distribution_pool
	// If win then increase win_cnt(prize), subtract remaining_qty(distribution_pool)
	if resPrizeInfo != nil {
		prize := schema.PrizeRow{
			PrizeId: resPrizeInfo.PrizeId,
			WinCnt:  resPrizeInfo.WinCnt + 1,
		}
		err = models.UpdatePrizeByPrizeId(tx, &prize)
		if err != nil {
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}

		pool := schema.PrizeDistPoolRow{
			DistPoolId:   resPrizeInfo.DistPoolId,
			RemainingQty: resPrizeInfo.RemainingQty - resPrizeInfo.Amount,
		}
		err = models.UpdateDistPoolByPoolId(tx, &pool)
		if err != nil {
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}

	// Table : account
	// Subtract ticket amount for doing game
	*ticketAmount -= ticketQtyForGame
	account.UpdatedAt = time.Time{}
	err = models.UpdateAccountTicketById(tx, &account)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 4. Create
	// Table : game_order
	jsonData, err = json.Marshal(order)
	if err = json.Unmarshal(jsonData, &order); err != nil {
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	orderWithId := schema.OrderRowWithID{}
	err = json.Unmarshal(jsonData, &orderWithId)
	if err = json.Unmarshal(jsonData, &order); err != nil {
		services.BadRequest(c, "bad request : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	err = models.CreateOrderWithTx(tx, &orderWithId)
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, types.ResStartGame{
		OrderId:       orderWithId.ID,
		AccountId:     order.AccountId,
		Addr:          order.Addr,
		PromotionId:   order.PromotionId,
		GameId:        order.GameId,
		Status:        order.Status,
		UsedTicketQty: order.UsedTicketQty,
		StartedAt:     order.StartedAt,
	})
}

// 게임 종료 (룰렛)
func StopGame(c *gin.Context) {
	/*
		orderId 불러오기
		order 조회하여 진행 중인지 확인
		is_win 에 따른 상태 업데이트
	*/
	// Fetch body: orderId
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	var order schema.OrderRow
	if err = json.Unmarshal(jsonData, &order); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)

			tx.Rollback()

			err = errors.New("panic")
			debug.PrintStack()
			services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}()
	err = tx.Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 주문 조회
	isLocked, err := models.QueryAndLockOrderById(tx, &order)
	if isLocked {
		// If is there any in stop event
		// then just ignore request
		services.Success(c, "there are some game stop event", nil)

		tx.Rollback()
		return
	}
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// Status 가 1 인지 확인 (게임 중인지 확인)
	if order.Status != 1 {
		err = errors.New("Game not in progress")
		fmt.Printf("%+v\n", err.Error())
		services.BadRequest(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)

		tx.Rollback()
		return
	}

	// 상태 업데이트
	if order.IsWin {
		// 당첨으로인한종료
		order.Status = 3
	} else {
		// 꽝으로인한종료
		order.Status = 2
	}
	order.UpdatedAt = time.Now()
	err = models.UpdateOrder(tx, &order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 주문 상세 정보 조회
	latestOrder := types.ResGetLatestOrderByAddr{OrderId: order.OrderId}
	if err = models.QueryOrderDetailById(tx, &latestOrder); err != nil {
		services.NotAcceptable(c, "fail QueryOrderDetailById "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// 당첨 시 토큰 전송
	if latestOrder.IsWin {
		err = middlewares.SendToken(latestOrder.Addr, int(latestOrder.Prize.Amount))
		if err != nil {
			services.NotAcceptable(c, "fail SendToken "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
			return
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	services.Success(c, nil, latestOrder)
}

// TODO
/*
	- claim 리스트
*/

func GetRandom(c *gin.Context) {
	// Generate random number 1 ~ 100000
	n, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	randNum := n.Uint64() + 1
	services.Success(c, nil, randNum)
}

func StartGameSample(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.GameOrder
	err := models.QueryCurGameByAddr(&game, iAddr)
	if err != nil {
		if err.Error() != "record not found" {
			services.NotAcceptable(c, "Something went wrong! Can not query user game", err)
			return
		}
	} else {
		// 이미 진행 중이라면
		services.NotAcceptable(c, "Ongoing game already exists", err)
		return
	}

	// 티켓이 존재하는지 확인
	var account schema.Account
	err = models.QueryBalanceByAddr(&account, iAddr)
	if err != nil {
		if err.Error() == "record not found" { // if addr not exists
			services.NotAcceptable(c, "user does not exists. game can not be started", err)
		} else {
			services.NotAcceptable(c, "Something went wrong! Can not query user balance", err)
		}
		return
	}
	if account.Ticket == 0 { // 티켓이 1개만 필요하다고 가정
		services.NotAcceptable(c, "you dont have any ticket", err)
		return
	}

	// 티켓 빼기 (현재는 1개만)
	if err = models.BurnUserTicket(&account, iAddr, int64(1)); err != nil {
		services.NotAcceptable(c, "ticket burn error. can not start game", err)
		return
	}

	// 게임 계산
	// 게임의 정책이 변경될 수 있으므로 시작 시점을 기준으로 결과를 계산한다.
	// math.Random
	// Generate random number 1 ~ 100000
	n, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		services.NotAcceptable(c, "fail : "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	resultNum := n.Uint64() + 1
	var isWin bool
	var prizeID int
	fmt.Println(resultNum)
	// 룰렛판(roulette_set) 에 따라 결과를 분기
	// 풀렛판은 동시에 2개 이상의 prize 를 얻을 수 없다. 즉 결과는 하나다.
	// 풀렛판은 각 prize + startNum + endNum 으로 구성. 해당 num 안에 포함될 경우
	// inWin && prizeID 부여. 어떤 num 에도 포함되지 않으면 꽝
	// 1상품의 확률이 3%
	// 2상품의 확률이 10%
	// 3상품의 확률이 50%
	// 나머지는 꽝이다
	// 1상품: 0 ~ 30
	// 2상품: 31 ~ 130
	// 3상품: 131 ~ 630
	// 꽝: 631 ~ 999

	if resultNum > 500 { // 꽝
		isWin = false
		prizeID = 0
	} else { //당첨
		isWin = true
		prizeID = 1
	}

	// 게임 저장
	game.IsWin = isWin
	game.PrizeID = int64(prizeID)
	game.PaidTicketNum = int64(1)
	err = models.StartNewGame(&game, iAddr)
	if err != nil {
		services.NotAcceptable(c, "Can not eoroll new game", err)
	} else {
		// 결과는 바로 나오지만 진행 중인 경우 결과를 보여주지 않는다.
		var gameInProgress schema.GameInProgress
		_ = models.QueryCurGameByAddr(&game, iAddr)
		gameInProgress.Address = game.Address
		gameInProgress.GameOrderId = game.GameOrderId
		gameInProgress.PaidTicketNum = game.PaidTicketNum
		gameInProgress.Status = game.Status
		gameInProgress.Type = game.Type
		services.Success(c, nil, gameInProgress)
	}
}

func GetOngoingGame(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.GameOrder
	err := models.QueryCurGameByAddr(&game, iAddr)
	if err != nil {
		if err.Error() != "record not found" {
			services.NotAcceptable(c, "Something went wrong! Can not query user game", err)
			return
		} else { // 진행 중인 게임이 없다면
			services.Success(c, nil, game)
			//services.NotAcceptable(c, "Ongoing game no exists", err)
			return
		}
	}

	services.Success(c, nil, game)
}

//-----------------------------------------------------------------------------------------

// 게임 종류 조회
func GetGames(c *gin.Context) {
	games := make([]schema.Game, 0, 100)
	err := models.QueryGameTypes(&games)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "GetGame fail", err)
		return
	}

	services.Success(c, nil, games)
}

// 게임 생성
func CreateGame(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request", err)
		return
	}
	var req types.ReqTbCreateGame
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	game := schema.Game{
		Title:    req.Title,
		Desc:     req.Desc,
		IsActive: req.IsActive,
		Url:      req.Url,
	}
	err = models.CreateGame(&game)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "CreateGame Fail", err)
		}
	} else {
		services.Success(c, nil, game)
	}
}

func GetGame(c *gin.Context) {
	// 파라미터 조회
	strGameId := c.Param("game_id")
	gameId, err := strconv.ParseInt(strGameId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request gameId path parameter", err)
		return
	}
	game := schema.Game{
		GameId: gameId,
	}
	err = models.QueryGameType(nil, &game)

	if err != nil {
		if err.Error() == "record not found" { // if addr not exists
			services.NotAcceptable(c, "record not found", err)
		} else {
			services.NotAcceptable(c, "failed GetGame", err)
		}
	} else {
		services.Success(c, nil, game)
	}
}

// 게임 정보 수정
func UpdateGame(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strGameId := c.Param("game_id")
	gameId, err := strconv.ParseInt(strGameId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request gameId path parameter", err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request", err)
		return
	}
	var req types.ReqTbUpdateGame
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	game := schema.Game{
		GameId:    gameId,
		Title:     req.Title,
		Desc:      req.Desc,
		IsActive:  req.IsActive,
		Url:       req.Url,
		UpdatedAt: time.Now(),
	}
	err = models.UpdateGame(&game)

	// result
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		if strings.Contains(err.Error(), "1062") {
			services.NotAcceptable(c, "Title already exists", err)
		} else {
			services.NotAcceptable(c, "UpdateGame Fail", err)
		}
	} else {
		services.Success(c, nil, game)
	}
}

// 게임 정보 삭제
func DeleteGame(c *gin.Context) {
	// 파라미터 조회
	strGameId := c.Param("game_id")
	gameId, err := strconv.ParseInt(strGameId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request gameId path parameter", err)
		return
	}

	// handler data
	game := schema.Game{
		GameId: gameId,
	}
	err = models.DeleteGame(&game)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed DeleteGame", err)
	} else {
		services.Success(c, nil, game)
	}
}
