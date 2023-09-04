package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	- claim 리스트
*/

func GetRandom(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(1000))
	services.Success(c, nil, rand.Intn(1000))
}

func StartGame(c *gin.Context) {
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

	if resultNum > 500 {	// 꽝
		isWin = false
		prizeID = 0
	} else {	//당첨
		isWin = true
		prizeID = 1
	}

	// 게임 저장
	game.IsWin  = isWin
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

func StopGame(c *gin.Context) {
	iAddr := c.Param("addr")

	// 이미 게임이 진행 중인지 확인
	var game schema.GameOrder
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
	var game schema.GameOrder
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

//-----------------------------------------------------------------------------------------

// 게임 종류 조회
func GetGames(c *gin.Context) {
	games := make([]schema.Game, 0, 100)
	err := models.QueryGameTypes(&games)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
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
	var req types.ReqCreateGame
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	game := schema.Game{
		Title: req.Title,
		Desc: req.Desc,
		IsActive: req.IsActive,
		Url: req.Url,
	}
	err = models.CreateGame(&game)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
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
		if err.Error() == "record not found" {	// if addr not exists
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
	var req types.ReqUpdateGame
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	game := schema.Game{
		GameId: gameId,
		Title: req.Title,
		Desc: req.Desc,
		IsActive: req.IsActive,
		Url: req.Url,
		UpdatedAt: time.Now(),
	}
	err = models.UpdateGame(&game)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
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