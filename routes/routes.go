package routes

import (
	"fmt"
	"roulette-api-server/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//-- game_order -> status: 1(진행중) 2(꽝으로인한종료) 3(클레임전) 4(클레임중) 5(클레임성공) 6(클레임실패) 7(취소)
	route := gin.Default() // gin Engine 초기화

	// CORS 설정
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowCredentials = false
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))
	fmt.Printf("%+v\n", config)
	cors.Default()

	route.GET("/promotions", controllers.GetPromotions)                   // [USER] 완료: 프로모션 정보 조회
	route.GET("/promotions/:promotion_id", controllers.GetPromotion)      // [USER] 완료: 프로모션 조회
	route.POST("/promotions", controllers.CreatePromotion)                // 프로모션 생성 (promotion + dPools + prizes)

	// game order
	route.POST("/game-mgmt/start", controllers.StartGame)    // [USER]TODO: review
	route.POST("/game-mgmt/stop", controllers.StopGame)      // [USER]게임 종료
	route.PATCH("/game-mgmt/claim/:addr", controllers.ClaimForBB)

	// account
	route.GET("/accounts", controllers.GetAccounts)                        // 계정들 조회
	route.GET("/accounts/detail", controllers.GetAccountsDetail)           // 계정들 상세 조회
	route.GET("/accounts/:addr", controllers.GetAccount)                   // [USER]상세 정보
	route.PUT("/accounts/:addr", controllers.PutAccount)                   // [USER]계정 생성
	route.GET("/accounts/:addr/orders", controllers.GetGameOrdersByAddr)   // 완료
	route.GET("/accounts/:addr/orders/latest", controllers.GetLatestOrder) // [USER]유저 최근 order 정보
	route.GET("/accounts/:addr/transfers", controllers.GetTransferEvents)  // [USER] TODO: 유저 별

	// game-mgmt
	route.GET("/game-mgmt/orders/winning-results", controllers.GetGameWinningResults)             // 당첨된 모든 orders
	route.GET("/game-mgmt/games", controllers.GetGames)                                           // 게임 모두 조회
	route.POST("/game-mgmt/games", controllers.CreateGame)                                        // 게임 생성
	route.GET("/game-mgmt/games/:game_id", controllers.GetGame)                                   // 특정 게임 조회

	// healthcheck
	route.GET("/healthcheck", controllers.GetHealthcheck)

	return route
}
