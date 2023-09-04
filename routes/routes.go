package routes

import (
	"roulette-api-server/controllers"
	"roulette-api-server/middlewares"
	tcontrollers "roulette-api-server/tcontrollers"

	// controllers "roulette-api-server/controllers"

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
	route.POST("auth/signin", middlewares.IsClientAuthenticated, tcontrollers.AuthSignin)
	route.DELETE("auth/signout", middlewares.IsUserAuthenticated, tcontrollers.AuthSignout)

	route.GET("users", middlewares.IsUserAuthenticated, tcontrollers.UserFetchAll)
	route.GET("users/:id", middlewares.IsUserAuthenticated, tcontrollers.UserFetchSingle)
	route.POST("users", middlewares.IsClientAuthenticated, tcontrollers.UserCreate)
	route.PUT("users/:id", middlewares.IsUserAuthenticated, tcontrollers.UserUpdate)
	route.DELETE("users/:id", middlewares.IsUserAuthenticated, tcontrollers.UserDelete)
	// end of boilplate
	
	//------------------------------------------------------------------------------
	// 룰렛 API samples
	//------------------------------------------------------------------------------
	route.GET("/tb/balance/users/:addr", tcontrollers.GetBalanceByAddr)										// 유저 밸런스 조회
	//TODO: change method GET -> POST
	route.GET("/tb/voucher/swap/:addr/:voucher_num", tcontrollers.SwapVoucherToTicket)		// 바우처 -> 티켓 스왑
	//TODO: change method GET -> POST
	route.GET("/tb/voucher/send/:addr/:voucher_num", tcontrollers.SendVoucher)						// 바우처 send
	route.GET("/tb/game/random", tcontrollers.GetRandom)																	// 난수 테스트
	//TODO: change method GET -> POST
	route.GET("/tb/game/start/:addr", tcontrollers.StartGame)															// 게임 시작
	//TODO: change method GET -> POST
	route.GET("/tb/game/stop/:addr", tcontrollers.StopGame)																// 게임 종료
	route.GET("/tb/game/ongoing/:addr", tcontrollers.GetOngoingGame)											// 현재 진행 중인 게임 조회

	//------------------------------------------------------------------------------
	// 룰렛 API 1차 개발
	//------------------------------------------------------------------------------
	route.GET   ("/promotions", controllers.GetPromotions)	// 프로모션 정보 조회
	

	//------------------------------------------------------------------------------
	// 룰렛 only 특정 테이블 CRUD APIs
	//------------------------------------------------------------------------------

	// promotion
	// TODO: should query promotion, dist_pool, prize at once											
	route.GET   ("/tb/promotions", tcontrollers.GetPromotions)
	// TODO: should create promotion, dist_pool, prize at once											
	route.POST  ("/tb/promotions", tcontrollers.CreatePromotion)									
	// TODO: should query promotion, dist_pool, prize at once											
	route.GET   ("/tb/promotions/:promotion_id", tcontrollers.GetPromotion)						
	route.PATCH ("/tb/promotions/:promotion_id/info", tcontrollers.UpdatePromotion)						
	// TODO: seperate promotion only & all promotion data including dist_pool, prize
	route.DELETE("/tb/promotions/:promotion_id", tcontrollers.DeletePromotion)	

	// prize denom
	// TODO: denom 별 상세한 정보(속하는 dist, prize 등 통계 정보) 추가
	route.GET   ("/tb/prize-mgmt/denoms", tcontrollers.GetPrizeDenoms)											
	route.POST  ("/tb/prize-mgmt/denoms", tcontrollers.CreatePrizeDenom)									
	route.GET   ("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.GetPrizeDenom)						
	route.PATCH ("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.UpdatePrizeDenom)						
	route.DELETE("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.DeletePrizeDenom)	

	// prize distribution pool
	// TODO: pool 별 속하는 prize 정보 함께 제공
	route.GET   ("/tb/prize-mgmt/pools", tcontrollers.GetDistPools) // for test
	route.POST  ("/tb/prize-mgmt/pools", tcontrollers.CreateDistPool) // for test				
	route.GET   ("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.GetDistPool) // for test					
	route.PATCH ("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.UpdateDistPool)	
	route.DELETE("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.DeleteDistPool)

	// prize
	route.GET   ("/tb/prize-mgmt/prizes", tcontrollers.GetPrizes) // for test
	route.POST  ("/tb/prize-mgmt/prizes", tcontrollers.CreatePrize) // for test				
	route.GET   ("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.GetPrize) // for test					
	route.PATCH ("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.UpdatePrize)	
	route.DELETE("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.DeletePrize)	

	// account
	route.GET   ("/tb/accounts", tcontrollers.GetAccounts) // for test
	route.POST  ("/tb/accounts", tcontrollers.CreateAccount) // for test
	route.GET   ("/tb/accounts/:addr", tcontrollers.GetAccount) // for test
	route.PATCH ("/tb/accounts/:addr", tcontrollers.UpdateAccount)
	route.DELETE("/tb/accounts/:addr", tcontrollers.DeleteAccount)

	// game-mgmt
	route.GET   ("/tb/game-mgmt/games", tcontrollers.GetGames)											
	route.POST  ("/tb/game-mgmt/games", tcontrollers.CreateGame)									
	route.GET   ("/tb/game-mgmt/games/:game_id", tcontrollers.GetGame)						
	route.PATCH ("/tb/game-mgmt/games/:game_id", tcontrollers.UpdateGame)						
	route.DELETE("/tb/game-mgmt/games/:game_id", tcontrollers.DeleteGame)		

	// game-order
	route.GET   ("/tb/game-mgmt/orders", tcontrollers.GetGameOrders)
	route.POST  ("/tb/game-mgmt/orders", tcontrollers.CreateGameOrder) // for test
	route.GET   ("/tb/game-mgmt/orders/:order_id", tcontrollers.GetGameOrder)
	route.PATCH ("/tb/game-mgmt/orders/:order_id", tcontrollers.UpdateGameOrder)
	route.DELETE("/tb/game-mgmt/orders/:order_id", tcontrollers.DeleteGameOrder)

	// voucher
	// TODO: route.GET   ("/tb/vouchers", controllers.GetAccounts) // 바우처 리스트
	route.GET   ("/tb/voucher-mgmt/balances", tcontrollers.GetVoucherBalances) // 바우처 밸런스 리스트
	route.POST  ("/tb/voucher-mgmt/balances", tcontrollers.CreateVoucherBalance) // for test
	route.GET   ("/tb/voucher-mgmt/balances/:id", tcontrollers.GetVoucherBalance) // for test
	route.PATCH ("/tb/voucher-mgmt/balances/:id", tcontrollers.UpdateVoucherBalance) // for test
	route.DELETE("/tb/voucher-mgmt/balances/:id", tcontrollers.DeleteVoucherBalance) // for test

	// voucher send history
	route.GET   ("/tb/voucher-mgmt/events/send", tcontrollers.GetVoucherSendEvents)
	route.POST  ("/tb/voucher-mgmt/events/send", tcontrollers.CreateVoucherSendEvent) // for test
	route.GET   ("/tb/voucher-mgmt/events/send/:id", tcontrollers.GetVoucherSendEvent) // for test
	route.PATCH ("/tb/voucher-mgmt/events/send/:id", tcontrollers.UpdateVoucherSendEvent) // for test
	route.DELETE("/tb/voucher-mgmt/events/send/:id", tcontrollers.DeleteVoucherSendEvent) // for test

	// voucher send history
	route.GET   ("/tb/voucher-mgmt/events/burn", tcontrollers.GetVoucherBurnEvents)
	route.POST  ("/tb/voucher-mgmt/events/burn", tcontrollers.CreateVoucherBurnEvent) // for test
	route.GET   ("/tb/voucher-mgmt/events/burn/:id", tcontrollers.GetVoucherBurnEvent) // for test
	route.PATCH ("/tb/voucher-mgmt/events/burn/:id", tcontrollers.UpdateVoucherBurnEvent) // for test
	route.DELETE("/tb/voucher-mgmt/evenets/burn/:id", tcontrollers.DeleteVoucherBurnEvent) // for test

	return route
}