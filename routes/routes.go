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
		
	// promotion
	// TODO: should query promotion, dist_pool, prize at once											
	route.GET   ("/promotions", controllers.GetPromotions)
	// TODO: should create promotion, dist_pool, prize at once											
	route.POST  ("/promotions", controllers.CreatePromotion)									
	// TODO: should query promotion, dist_pool, prize at once											
	route.GET   ("/promotions/:promotion_id", controllers.GetPromotion)						
	route.PATCH ("/promotions/:promotion_id/info", controllers.UpdatePromotion)						
	// TODO: seperate promotion only & all promotion data including dist_pool, prize
	route.DELETE("/promotions/:promotion_id", controllers.DeletePromotion)	

	// prize denom
	// TODO: denom 별 상세한 정보(속하는 dist, prize 등 통계 정보) 추가
	route.GET   ("/prize-mgmt/denoms", controllers.GetPrizeDenoms)											
	route.POST  ("/prize-mgmt/denoms", controllers.CreatePrizeDenom)									
	route.GET   ("/prize-mgmt/denoms/:prize_denom_id", controllers.GetPrizeDenom)						
	route.PATCH ("/prize-mgmt/denoms/:prize_denom_id", controllers.UpdatePrizeDenom)						
	route.DELETE("/prize-mgmt/denoms/:prize_denom_id", controllers.DeletePrizeDenom)	

	// prize distribution pool
	// TODO: pool 별 속하는 prize 정보 함께 제공
	route.GET   ("/prize-mgmt/pools", controllers.GetDistPools) // for test
	route.POST  ("/prize-mgmt/pools", controllers.CreateDistPool) // for test				
	route.GET   ("/prize-mgmt/pools/:dist_pool_id", controllers.GetDistPool) // for test					
	route.PATCH ("/prize-mgmt/pools/:dist_pool_id", controllers.UpdateDistPool)	
	route.DELETE("/prize-mgmt/pools/:dist_pool_id", controllers.DeleteDistPool)

	// prize
	route.GET   ("/prize-mgmt/prizes", controllers.GetPrizes) // for test
	route.POST  ("/prize-mgmt/prizes", controllers.CreatePrize) // for test				
	route.GET   ("/prize-mgmt/prizes/:prize_id", controllers.GetPrize) // for test					
	route.PATCH ("/prize-mgmt/prizes/:prize_id", controllers.UpdatePrize)	
	route.DELETE("/prize-mgmt/prizes/:prize_id", controllers.DeletePrize)	

	// account
	route.GET   ("/accounts", controllers.GetAccounts) // for test
	route.POST  ("/accounts", controllers.CreateAccount) // for test
	route.GET   ("/accounts/:addr", controllers.GetAccount) // for test
	route.PATCH ("/accounts/:addr", controllers.UpdateAccount)
	route.DELETE("/accounts/:addr", controllers.DeleteAccount)

	// game-mgmt
	route.GET   ("/game-mgmt/games", controllers.GetGames)											
	route.POST  ("/game-mgmt/games", controllers.CreateGame)									
	route.GET   ("/game-mgmt/games/:game_id", controllers.GetGame)						
	route.PATCH ("/game-mgmt/games/:game_id", controllers.UpdateGame)						
	route.DELETE("/game-mgmt/games/:game_id", controllers.DeleteGame)		

	// game-order
	route.GET   ("/game-mgmt/orders", controllers.GetGameOrders)
	route.POST  ("/game-mgmt/orders", controllers.CreateGameOrder) // for test
	route.GET   ("/game-mgmt/orders/:order_id", controllers.GetGameOrder)
	route.PATCH ("/game-mgmt/orders/:order_id", controllers.UpdateGameOrder)
	route.DELETE("/game-mgmt/orders/:order_id", controllers.DeleteGameOrder)

	// voucher
	// TODO: route.GET   ("/vouchers", controllers.GetAccounts) // 바우처 리스트
	route.GET   ("/voucher-mgmt/balances", controllers.GetVoucherBalances) // 바우처 밸런스 리스트
	route.POST  ("/voucher-mgmt/balances", controllers.CreateVoucherBalance) // for test
	route.GET   ("/voucher-mgmt/balances/:id", controllers.GetVoucherBalance) // for test
	route.PATCH ("/voucher-mgmt/balances/:id", controllers.UpdateVoucherBalance) // for test
	route.DELETE("/voucher-mgmt/balances/:id", controllers.DeleteVoucherBalance) // for test

	// voucher send history
	route.GET   ("/voucher-mgmt/events/send", controllers.GetVoucherSendEvents)
	route.POST  ("/voucher-mgmt/events/send", controllers.CreateVoucherSendEvent) // for test
	route.GET   ("/voucher-mgmt/events/send/:id", controllers.GetVoucherSendEvent) // for test
	route.PATCH ("/voucher-mgmt/events/send/:id", controllers.UpdateVoucherSendEvent) // for test
	route.DELETE("/voucher-mgmt/events/send/:id", controllers.DeleteVoucherSendEvent) // for test

	// voucher send history
	route.GET   ("/voucher-mgmt/events/burn", controllers.GetVoucherBurnEvents)
	route.POST  ("/voucher-mgmt/events/burn", controllers.CreateVoucherBurnEvent) // for test
	route.GET   ("/voucher-mgmt/events/burn/:id", controllers.GetVoucherBurnEvent) // for test
	route.PATCH ("/voucher-mgmt/events/burn/:id", controllers.UpdateVoucherBurnEvent) // for test
	route.DELETE("/voucher-mgmt/evenets/burn/:id", controllers.DeleteVoucherBurnEvent) // for test


	//------------------------------------------------------------------------------
	// 룰렛 API 1차 개발 optonal APIs
	//------------------------------------------------------------------------------


	return route
}