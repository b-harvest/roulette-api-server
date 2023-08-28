package routes

import (
	"roulette-api-server/controllers"
	"roulette-api-server/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))

	//------------------------------------------------------------------------------
	// boilplate
	//------------------------------------------------------------------------------
	route.POST("auth/signin", middlewares.IsClientAuthenticated, controllers.AuthSignin)
	route.DELETE("auth/signout", middlewares.IsUserAuthenticated, controllers.AuthSignout)

	route.GET("users", middlewares.IsUserAuthenticated, controllers.UserFetchAll)
	route.GET("users/:id", middlewares.IsUserAuthenticated, controllers.UserFetchSingle)
	route.POST("users", middlewares.IsClientAuthenticated, controllers.UserCreate)
	route.PUT("users/:id", middlewares.IsUserAuthenticated, controllers.UserUpdate)
	route.DELETE("users/:id", middlewares.IsUserAuthenticated, controllers.UserDelete)
	// end of boilplate
	
	// 유저 밸런스 조회
	route.GET("/balance/users/:addr", controllers.GetBalanceByAddr)
	// 바우처 -> 티켓 스왑
	route.GET("/voucher/swap/:addr/:voucher_num", controllers.SwapVoucherToTicket)
	// 바우처 send
	route.GET("/voucher/send/:addr/:voucher_num", controllers.SendVoucher)
	// 게임 시작
	route.GET("/game/random", controllers.GetRandom)	// 난수 테스트
	route.GET("/game/start/:addr", controllers.StartGame)
	// 게임 종료
	route.GET("/game/stop/:addr", controllers.StopGame)
	// 현재 진행 중인 게임 조회
	route.GET("/game/ongoing/:addr", controllers.GetOngoingGame)

	return route
}