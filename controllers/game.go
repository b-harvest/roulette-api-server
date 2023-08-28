package controllers

import (
	"fmt"
	"math/rand"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRandom(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(1000))
	services.Success(c, nil, rand.Intn(1000))
}

func StartGame(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.Game
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
		if err.Error() == "record not found" {	// if addr not exists
			services.NotAcceptable(c, "user does not exists. game can not be started", err)
		} else {
			services.NotAcceptable(c, "Something went wrong! Can not query user balance", err)
		}
		return
	}
	if account.Ticket == 0 {	// 티켓이 1개만 필요하다고 가정
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
	rand.Seed(time.Now().UnixNano())
	resultNum := rand.Intn(1000)
	var isWin bool
	var giftId int
	fmt.Println(resultNum)
	// 룰렛판(roulette_set) 에 따라 결과를 분기
	// 풀렛판은 동시에 2개 이상의 gift 를 얻을 수 없다. 즉 결과는 하나다.
	// 풀렛판은 각 gift + startNum + endNum 으로 구성. 해당 num 안에 포함될 경우
	// inWin && giftId 부여. 어떤 num 에도 포함되지 않으면 꽝
	if resultNum > 500 {	// 꽝
		isWin = false
		giftId = 0
	} else {	//당첨
		isWin = true
		giftId = 1
	}

	// 게임 저장
	game.IsWin  = isWin
	game.GiftId = int64(giftId)
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

func StopGame(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.Game
	err := models.QueryCurGameByAddr(&game, iAddr)
	if err != nil {
		if err.Error() != "record not found" {
			services.NotAcceptable(c, "Something went wrong! Can not query user game", err)
			return
		} else {		// 진행 중인 게임이 없다면
			services.NotAcceptable(c, "Ongoing game no exists", err)
			return
		}
	}

	// 상태 업데이트
	game.Status = 2
	err = models.StopGame(&game, iAddr)
	if err != nil {
		services.NotAcceptable(c, "Can not stop existing game", err)
	} else {
		_ = models.QueryCurGameByAddr(&game, iAddr)
		services.Success(c, nil, game)
	}
}

func GetOngoingGame(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.Game
	err := models.QueryCurGameByAddr(&game, iAddr)
	if err != nil {
		if err.Error() != "record not found" {
			services.NotAcceptable(c, "Something went wrong! Can not query user game", err)
			return
		} else {		// 진행 중인 게임이 없다면
			services.Success(c, nil, game)
			//services.NotAcceptable(c, "Ongoing game no exists", err)
			return
		}
	}

	services.Success(c, nil, game)
}