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

type (
	GamesInMem struct {
		// promotion id => prize info
		Games map[int64](*PrizeInfoInMem)
	}

	PrizeInfoInMem struct {
		PrizeIds []int64
		ProbabilityList []uint64
	}
)
var gamesInMem GamesInMem
var maxNum int64

func init() {
	// Initialization
	gamesInMem = GamesInMem{
		Games: map[int64](*PrizeInfoInMem){},
	}
	maxNum = 100000
}

func getPrizeInfo(promotionId int64) (*PrizeInfoInMem, error) {
	// hit
	// Check gameId is exists in GamesInMem
	// if gamesInMem.Games[promotionId] != nil {
	// 	return gamesInMem.Games[promotionId], nil
	// }

	// miss
	// create prize Info
	var prizes []schema.PrizeRow
	err := models.QueryPrizesByPromotionId(&prizes, promotionId)
	if err != nil {
		return nil, err
	}

	gamesInMem.Games[promotionId] = &PrizeInfoInMem{
		[]int64{},
		[]uint64{},
	}

	var accum uint64 = 1
	for _, prize := range prizes {
		gamesInMem.Games[promotionId].PrizeIds = append(gamesInMem.Games[promotionId].PrizeIds, prize.PrizeId)

		accum += uint64(prize.Odds * float64(1000))
		gamesInMem.Games[promotionId].ProbabilityList = append(gamesInMem.Games[promotionId].ProbabilityList, accum)
	}
	
	return gamesInMem.Games[promotionId], nil
}

// return prizeId and error
// if prizeId == 0
// then not win
func inGame(promotionId int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(maxNum))
	if err != nil {
		return -1, err
	}
	// 1 ~ maxNum
	randNum := n.Uint64() + 1

	prizeInfo, err := getPrizeInfo(promotionId)
	if err != nil {
		return -1, err
	}
	var prevNum uint64 = 1
	for i, probability := range prizeInfo.ProbabilityList {
		// preNum <= randNum < n
		if prevNum <= randNum && randNum < probability {
			return prizeInfo.PrizeIds[i], nil
		}
		prevNum = probability
	}
	return 0, nil
}

func genOrders(order *schema.OrderRow) error {
	order.StartedAt = time.Now()

	var err error
	order.PrizeId, err = inGame(order.PromotionId)
	if err != nil {
		return err
	}
	if order.PrizeId == 0 {
		// not win
		order.IsWin = false
	} else {
		order.IsWin = true
	}

	// Game in progress
	order.Status = 1

	return nil
}

func StartGame(c *gin.Context) {
	// 진행 중인 게임이 있는지 확인

	// Fetch body
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var order schema.OrderRow
	if err = json.Unmarshal(jsonData, &order); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r!= nil {
			fmt.Println(r)

			tx.Rollback()

			err = errors.New("panic")
			debug.PrintStack()
			services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
		}
	}()
	err = tx.Error
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// 1. whether exists account
	account := schema.AccountRow{
		Addr: order.Addr,
	}
	err = models.QueryAccountByAddr(&account)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	order.AccountId = account.Id

	// 2. subtract account ticket

	// if account has insufficient ticket amount
	// then error
	// TODO: calculate required ticket amount for game start
	//if account.TicketAmount < order.UsedTicketQty {
	ticketQtyForGame := uint64(1)
	order.UsedTicketQty = ticketQtyForGame
	if account.TicketAmount < ticketQtyForGame {
		err = errors.New("Insufficient ticket amount")
		fmt.Printf("%+v\n", err.Error())
		services.BadRequest(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}

	// update account
	account.TicketAmount -= ticketQtyForGame
	err = models.UpdateAccountById(tx, &account)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// 3. do game
	err = genOrders(&order)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// create orders (insert as batch)
	err = models.CreateOrderWithTx(tx, &order)
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	//결과는 숨긴다.
	// order.IsWin
	// order.PrizeId
	services.Success(c, nil, &order)
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
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var order schema.OrderRow
	if err = json.Unmarshal(jsonData, &order); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// 주문 조회
	err = models.QueryOrderById(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	// Status 가 1 인지 확인 (게임 중인지 확인)
	if order.Status != 1 {
		err = errors.New("Game not in progress")
		fmt.Printf("%+v\n", err.Error())
		services.BadRequest(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
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
	err = models.UpdateOrder(&order)
	if err != nil {
		fmt.Printf("%+v\n", err.Error())
		services.NotAcceptable(c, "fail "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	
	// 주문 상세 정보 조회
	latestOrder := types.ResGetLatestOrderByAddr{OrderId: order.OrderId}
	if err = models.QueryOrderDetailById(&latestOrder); err != nil {
		services.NotAcceptable(c, "fail QueryOrderDetailById "+c.Request.Method+" "+c.Request.RequestURI+" : "+err.Error(), err)
		return
	}
	
	services.Success(c, nil, latestOrder)
}


// TODO
/*
	- claim 리스트
*/

func GetRandom(c *gin.Context) {
	n, err := rand.Int(rand.Reader, big.NewInt(maxNum))
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	// 1 ~ maxNum
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
	n, err := rand.Int(rand.Reader, big.NewInt(maxNum))
	if err = models.BurnUserTicket(&account, iAddr, int64(1)); err != nil {
		services.NotAcceptable(c, "ticket burn error. can not start game", err)
		return
	}
	// 1 ~ maxNum
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
	err = models.QueryGameType(&game)

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
