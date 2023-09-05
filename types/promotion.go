package types

/*
	- 당첨확률 dist_pool
	- 총 주문 수 game_order
	- 총 주문자 수 game_order
	- 총 당첨 수 game_order
	- 총 당첨자 수 game_order
	- 총 claim 요청자 수 game_order
	- 미처리 claim 수 game_order
	- 총 바우처 burn 갯수
*/
type PromotionSummary struct {
	//-- game_order -> status: 1(진행중) 2(꽝으로인한종료) 3(클레임전) 4(클레임중) 5(클레임성공) 6(클레임실패) 7(취소)
	TotalOdds               float64   `json:"totalOdds" db:"total_odds" gorm:"column:total_odds;type:uint64"`
	TotalOrderNum           uint64    `json:"totalOrderNum" db:"total_order_num"`
	TotalOrderUserNum       uint64    `json:"totalOrderUserNum" db:"total_order_user_num"`
	TotalWinNum             uint64    `json:"totalWinNum" db:"total_win_num"`
	TotalWinUserNum         uint64    `json:"totalWinUserNum" db:"total_win_user_num"`
	TotalClaimedNum         uint64    `json:"totalClaimedNum" db:"total_claimed_num"`
	TotalClaimedUserNum     uint64    `json:"totalClaimedUserNum" db:"total_claimed_user_num"`
	InProgressClaimNum      uint64    `json:"inProgressClaimNum" db:"in_progress_claim_num"`
	// 지갑 연결 수
	// 지갑 연결자 수
}