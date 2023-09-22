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
	//-- game_order -> status: 1(진행중) 2(꽝으로인한종료) 3(클레임전) 4(클레임중) 5(클레임성공) 6(클레임실패) 7(취소)
	route := gin.Default() // gin Engine 초기화

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
	route.GET("/tb/balance/users/:addr",             tcontrollers.GetBalanceByAddr)    // 유저 밸런스 조회
	route.GET("/tb/voucher/swap/:addr/:voucher_num", tcontrollers.SwapVoucherToTicket) // 바우처 -> 티켓 스왑
	route.GET("/tb/voucher/send/:addr/:voucher_num", tcontrollers.SendVoucher)         // 바우처 send
	route.GET("/tb/game/random",                     tcontrollers.GetRandom)           // 난수 테스트
	route.GET("/tb/game/start/:addr",                tcontrollers.StartGame)           // 게임 시작
	route.GET("/tb/game/stop/:addr",                 tcontrollers.StopGame)            // 게임 종료
	route.GET("/tb/game/ongoing/:addr",              tcontrollers.GetOngoingGame)      // 현재 진행 중인 게임 조회


	//------------------------------------------------------------------------------
	// 룰렛 API 1차 개발
	//------------------------------------------------------------------------------
	route.GET  ("/promotions",               controllers.GetPromotions)     // [USER] 완료: 프로모션 정보 조회
	route.GET  ("/promotions/:promotion_id", controllers.GetPromotion)      // [USER] 완료: 프로모션 조회
	route.POST ("/promotions",               controllers.CreatePromotion)   // 프로모션 생성 (promotion + dPools + prizes)
	route.PATCH("/promotions/:promotion_id", controllers.UpdatePromotion)   // 프로모션 수정 (promotion + dPools + prizes)

	// game order
	route.POST("/game-mgmt/start", controllers.StartGame)   // [USER]TODO: review
	route.POST("/game-mgmt/stop",  controllers.StopGame)    // [USER]게임 종료
	route.POST("/game-mgmt/claim", controllers.Claim)       // [USER]특정 order 클레임

	// account
	route.GET  ("/accounts",                        controllers.GetAccounts)		      // 계정들 조회
	route.GET  ("/accounts/detail",                 controllers.GetAccountsDetail)		// 계정들 상세 조회
	route.GET  ("/accounts/:addr",                  controllers.GetAccount)	          // [USER]상세 정보
	route.PUT  ("/accounts/:addr",                  controllers.PutAccount)	          // [USER]계정 생성
	route.GET  ("/accounts/:addr/orders",           controllers.GetGameOrdersByAddr)  // 완료
	route.GET  ("/accounts/:addr/orders/latest",    controllers.GetLatestOrder)       // [USER]유저 최근 order 정보
	route.GET  ("/accounts/:addr/transfers",        controllers.GetTransferEvents)    // [USER] TODO: 유저 별
	// route.GET  ("/accounts/:addr/balances",         controllers.GetBalancesByAddr)    // [USER] -> GetAccount 에서 커버 가능
	// route.GET  ("/accounts/:addr/winning-records",  controllers.GetWinTotalByAcc)     // [USER]유저 prize 별 총 당첨 amt
	

	// metrics
	// wallet-connects
	route.GET ("/metrics/wallet-connects",       controllers.GetEventWalletConn)       // wallet접속 내역
	route.GET ("/metrics/wallet-connects/count", controllers.GetEventWalletConnCount)  // wallet접속 내역 총 cnt
	route.POST("/metrics/wallet-connects",       controllers.PostEventWalletConn)      // wallet접속 내역 생성
	// TODO: 통계 정보: 유저 별, daily 등

	// flip-links
	route.GET ("/metrics/flip-links",       controllers.GetEventFlipLink)            // link클릭 내역
	route.GET ("/metrics/flip-links/count", controllers.GetEventFlipLinkCount)       // link클릭 내역 총 cnt
	route.POST("/metrics/flip-links",       controllers.PostEventFlipLinks)          // link접속 내역 생성
	// TODO: 통계 정보: 유저 별, daily 등

	// voucher-mgmt
	route.GET ("/voucher-mgmt/events/send",        controllers.GetVoucherSendEvents)      // 바우처 send 내역
	route.POST("/voucher-mgmt/events/send",        controllers.CreateVoucherSendEvents)   // 바우처 보내기
	route.GET ("/voucher-mgmt/available-vouchers", controllers.GetAvailableVouchers)      // 프로모션 별 voucher 정보
	route.POST("/voucher-mgmt/burn",               controllers.PostVoucherBurn)

	// game-mgmt
	route.GET   ("/game-mgmt/orders/winning-results",           controllers.GetGameWinningResults) // 당첨된 모든 orders
	route.PATCH ("/game-mgmt/orders/winning-results/:order_id", controllers.UpdateGameOrderStatus) // order 상태 변경
	route.GET   ("/game-mgmt/games",          controllers.GetGames)   // 게임 모두 조회
	route.POST  ("/game-mgmt/games",          controllers.CreateGame) // 게임 생성
	route.GET   ("/game-mgmt/games/:game_id", controllers.GetGame)    // 특정 게임 조회
	route.PATCH ("/game-mgmt/games/:game_id", controllers.UpdateGame) // 게임 수정
	route.DELETE("/game-mgmt/games/:game_id", controllers.DeleteGame) // 게임 삭제

	// prize-mgmt
	route.GET   ("/prize-mgmt/denoms",      controllers.GetPrizeDenoms)    // 데놈 모두 조회
	route.POST  ("/prize-mgmt/denoms",      controllers.CreatePrizeDenom)  // 데놈 생성
	route.GET   ("/prize-mgmt/denoms/:id",  controllers.GetPrizeDenom)     // 특정 데놈 조회 
	route.PATCH ("/prize-mgmt/denoms/:id",  controllers.UpdatePrizeDenom)  // 데놈 수정
	route.DELETE("/prize-mgmt/denoms/:id",  controllers.DeletePrizeDenom)  // 데놈 삭제

	//------------------------------------------------------------------------------
	// 룰렛 only 특정 테이블 CRUD APIs
	//------------------------------------------------------------------------------

	// promotion
	route.GET("/tb/promotions", tcontrollers.GetPromotions)
	route.POST("/tb/promotions", tcontrollers.CreatePromotion)
	route.GET("/tb/promotions/:promotion_id", tcontrollers.GetPromotion)
	route.PATCH("/tb/promotions/:promotion_id/info", tcontrollers.UpdatePromotion)
	route.DELETE("/tb/promotions/:promotion_id", tcontrollers.DeletePromotion)

	// prize denom
	// TODO: denom 별 상세한 정보(속하는 dist, prize 등 통계 정보) 추가
	route.GET("/tb/prize-mgmt/denoms", tcontrollers.GetPrizeDenoms)
	route.POST("/tb/prize-mgmt/denoms", tcontrollers.CreatePrizeDenom)
	route.GET("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.GetPrizeDenom)
	route.PATCH("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.UpdatePrizeDenom)
	route.DELETE("/tb/prize-mgmt/denoms/:prize_denom_id", tcontrollers.DeletePrizeDenom)

	// prize distribution pool
	route.GET("/tb/prize-mgmt/pools", tcontrollers.GetDistPools)              // for test
	route.POST("/tb/prize-mgmt/pools", tcontrollers.CreateDistPool)           // for test
	route.GET("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.GetDistPool) // for test
	route.PATCH("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.UpdateDistPool)
	route.DELETE("/tb/prize-mgmt/pools/:dist_pool_id", tcontrollers.DeleteDistPool)

	// prize
	route.GET("/tb/prize-mgmt/prizes", tcontrollers.GetPrizes)          // for test
	route.POST("/tb/prize-mgmt/prizes", tcontrollers.CreatePrize)       // for test
	route.GET("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.GetPrize) // for test
	route.PATCH("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.UpdatePrize)
	route.DELETE("/tb/prize-mgmt/prizes/:prize_id", tcontrollers.DeletePrize)

	// account
	route.GET("/tb/accounts", tcontrollers.GetAccounts)      // for test
	route.POST("/tb/accounts", tcontrollers.CreateAccount)   // for test
	route.GET("/tb/accounts/:addr", tcontrollers.GetAccount) // for test
	route.PATCH("/tb/accounts/:addr", tcontrollers.UpdateAccount)
	route.DELETE("/tb/accounts/:addr", tcontrollers.DeleteAccount)

	// game-mgmt
	route.GET("/tb/game-mgmt/games", tcontrollers.GetGames)
	route.POST("/tb/game-mgmt/games", tcontrollers.CreateGame)
	route.GET("/tb/game-mgmt/games/:game_id", tcontrollers.GetGame)
	route.PATCH("/tb/game-mgmt/games/:game_id", tcontrollers.UpdateGame)
	route.DELETE("/tb/game-mgmt/games/:game_id", tcontrollers.DeleteGame)

	// game-order
	route.GET("/tb/game-mgmt/orders", tcontrollers.GetGameOrders)
	route.POST("/tb/game-mgmt/orders", tcontrollers.CreateGameOrder) // for test
	route.GET("/tb/game-mgmt/orders/:order_id", tcontrollers.GetGameOrder)
	route.PATCH("/tb/game-mgmt/orders/:order_id", tcontrollers.UpdateGameOrder)
	route.DELETE("/tb/game-mgmt/orders/:order_id", tcontrollers.DeleteGameOrder)

	// voucher
	// TODO: route.GET   ("/tb/vouchers", controllers.GetAccounts) // 바우처 리스트
	route.GET("/tb/voucher-mgmt/balances", tcontrollers.GetVoucherBalances)          // 바우처 밸런스 리스트
	route.POST("/tb/voucher-mgmt/balances", tcontrollers.CreateVoucherBalance)       // for test
	route.GET("/tb/voucher-mgmt/balances/:id", tcontrollers.GetVoucherBalance)       // for test
	route.PATCH("/tb/voucher-mgmt/balances/:id", tcontrollers.UpdateVoucherBalance)  // for test
	route.DELETE("/tb/voucher-mgmt/balances/:id", tcontrollers.DeleteVoucherBalance) // for test

	// voucher send history
	route.GET("/tb/voucher-mgmt/events/send", tcontrollers.GetVoucherSendEvents)
	route.POST("/tb/voucher-mgmt/events/send", tcontrollers.CreateVoucherSendEvent)       // for test
	route.GET("/tb/voucher-mgmt/events/send/:id", tcontrollers.GetVoucherSendEvent)       // for test
	route.PATCH("/tb/voucher-mgmt/events/send/:id", tcontrollers.UpdateVoucherSendEvent)  // for test
	route.DELETE("/tb/voucher-mgmt/events/send/:id", tcontrollers.DeleteVoucherSendEvent) // for test

	// voucher burn history
	route.GET("/tb/voucher-mgmt/events/burn", tcontrollers.GetVoucherBurnEvents)
	route.GET("/tb/voucher-mgmt/events/burn/:id", tcontrollers.GetVoucherBurnEvent)        // for test
	route.PATCH("/tb/voucher-mgmt/events/burn/:id", tcontrollers.UpdateVoucherBurnEvent)   // for test
	route.DELETE("/tb/voucher-mgmt/evenets/burn/:id", tcontrollers.DeleteVoucherBurnEvent) // for test

	// stats : Global
	route.GET("/stats/accounts", controllers.GetAccountStat)
	route.GET("/stats/promotions", controllers.GetPromotionStat)

	// stats : Promotion specific
	route.GET("/stats/flip-links/:promotion_id", controllers.GetFlipLinkStat)
	route.GET("/stats/wallet-connects/:promotion_id", controllers.GetWalletConnectStat)
	route.GET("/stats/vouchers/:promotion_id", controllers.GetVoucherStat)
	route.GET("/stats/tickets/:promotion_id", controllers.GetTicketStat)
	route.GET("/stats/prizes/:promotion_id", controllers.GetPrizeStat)

	return route
}
