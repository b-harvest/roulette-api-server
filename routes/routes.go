package routes

import (
	"roulette-api-server/controllers"
	"roulette-api-server/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()			// gin Engine 초기화

	// CORS 설정
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))

	//------------------------------------------------------------------------------
	// test
	//------------------------------------------------------------------------------
	route.POST("auth/signin", middlewares.IsClientAuthenticated, controllers.AuthSignin)
	route.DELETE("auth/signout", middlewares.IsUserAuthenticated, controllers.AuthSignout)

	route.GET("users", middlewares.IsUserAuthenticated, controllers.UserFetchAll)
	route.GET("users/:id", middlewares.IsUserAuthenticated, controllers.UserFetchSingle)
	route.POST("users", middlewares.IsClientAuthenticated, controllers.UserCreate)
	route.PUT("users/:id", middlewares.IsUserAuthenticated, controllers.UserUpdate)
	route.DELETE("users/:id", middlewares.IsUserAuthenticated, controllers.UserDelete)
	// end of boilplate
	
	//------------------------------------------------------------------------------
	// 룰렛 API samples
	//------------------------------------------------------------------------------
	route.GET("/balance/users/:addr", controllers.GetBalanceByAddr)										// 유저 밸런스 조회
	//TODO: change method GET -> POST
	route.GET("/voucher/swap/:addr/:voucher_num", controllers.SwapVoucherToTicket)		// 바우처 -> 티켓 스왑
	//TODO: change method GET -> POST
	route.GET("/voucher/send/:addr/:voucher_num", controllers.SendVoucher)						// 바우처 send
	route.GET("/game/random", controllers.GetRandom)																	// 난수 테스트
	//TODO: change method GET -> POST
	route.GET("/game/start/:addr", controllers.StartGame)															// 게임 시작
	//TODO: change method GET -> POST
	route.GET("/game/stop/:addr", controllers.StopGame)																// 게임 종료
	route.GET("/game/ongoing/:addr", controllers.GetOngoingGame)											// 현재 진행 중인 게임 조회

	//------------------------------------------------------------------------------
	// 룰렛 API 1차 개발
	//------------------------------------------------------------------------------
	
	// game-mgmt
	route.GET("/game-mgmt/games", controllers.GetGames)											
	route.POST("/game-mgmt/games", controllers.CreateGame)									
	route.PATCH("/game-mgmt/games/:game_id", controllers.UpdateGame)						
	route.DELETE("/game-mgmt/games/:game_id", controllers.DeleteGame)						
	


	return route
}