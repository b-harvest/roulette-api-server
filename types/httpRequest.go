package types

import "time"

type ReqTbCreateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
	//CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	//UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ReqTbUpdateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
}

type ReqTbCreatePromotion struct {
	Title 			           string    `json:"title" db:"title"`
	Desc 			             string    `json:"desc" db:"desc"`
	IsActive               bool      `json:"isActive" db:"is_active"`
	IsWhitelisted          bool      `json:"isWhitelisted" db:"is_whitelisted"`
	VoucherName 			     string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0  int       `json:"voucherExchangeRatio0" db:"voucher_exchange_ratio_0"`
	VoucherExchangeRatio1  int       `json:"voucherExchangeRatio1" db:"voucher_exchange_ratio_1"`
	VoucherTotalSupply     uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	PromotionStartAt       time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt         time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt           time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt             time.Time `json:"claimEndAt" db:"claim_end_at"`
}

type ReqTbUpdatePromotion struct {
	Title 			           string    `json:"title" db:"title"`
	Desc 			             string    `json:"desc" db:"desc"`
	IsActive               bool      `json:"isActive" db:"is_active"`
	IsWhitelisted          bool      `json:"isWhitelisted" db:"is_whitelisted"`
	VoucherName 			     string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0  int       `json:"voucherExchangeRatio0" db:"voucher_exchange_ratio_0"`
	VoucherExchangeRatio1  int       `json:"voucherExchangeRatio1" db:"voucher_exchange_ratio_1"`
	VoucherTotalSupply     uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	VoucherRemainingQty    uint64    `json:"voucherRemainingQty" db:"voucher_remaining_qty"`
	PromotionStartAt       time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt         time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt           time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt             time.Time `json:"claimEndAt" db:"claim_end_at"`
}

type ReqTbCreatePrizeDenom struct {
	Name 			             string    `json:"name" db:"name"`
	Type 			             string    `json:"type" db:"type"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqTbUpdatePrizeDenom struct {
	Name 			             string    `json:"name" db:"name"`
	Type 			             string    `json:"type" db:"type"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqTbCreateDistPool struct {
	PromotionId            int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId           int64     `json:"prizeDenomId" db:"prize_denom_id"`
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqTbUpdateDistPool struct {
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	RemainingQty           uint64    `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqTbCreatePrize struct {
	DistPoolId        int64     `json:"distPoolId" db:"dist_pool_id"`
	PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId      int64     `json:"prizeDenomId" db:"dist_pool_id"`
	Amount            uint64    `json:"amount" db:"amount"`
	Odds              float64   `json:"odds" db:"odds"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	IsActive          bool      `json:"isActive" db:"is_active"`
}

type ReqTbUpdatePrize struct {
	Odds              float64   `json:"odds" db:"odds"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	IsActive          bool      `json:"isActive" db:"is_active"`
}

type ReqTbCreateAccount struct {
	Addr 			     string    `json:"addr" db:"addr"`
	TicketAmount   uint64    `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo      string    `json:"adminMemo" db:"admin_memo"`
	Type 			     string    `json:"type" db:"type"`
	IsBlacklisted  bool      `json:"isBlacklisted" db:"is_blacklisted"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ReqTbUpdateAccount struct {
	TicketAmount   uint64    `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo      string    `json:"adminMemo" db:"admin_memo"`
	Type 			     string    `json:"type" db:"type"`
	IsBlacklisted  bool      `json:"isBlacklisted" db:"is_blacklisted"`
}

type ReqTbCreateOrder struct {
	AccountId         int64     `json:"accountId" db:"account_id"`
	Addr 			        string    `json:"addr" db:"addr"`
	PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	GameId            int64     `json:"gameId" db:"game_id"`
	Status            int       `json:"status" db:"status"`
	UsedTicketQty     uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
}

type ReqTbUpdateOrder struct {
	IsWin             bool      `json:"isWin" db:"is_win"`
	Status            int       `json:"status" db:"status"`
	UsedTicketQty     uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
	ClaimedAt         time.Time `json:"claimedAt" db:"claimed_at"`
	ClaimFinishedAt   time.Time `json:"claimFinishedAt" db:"claim_finished_at"`
}

type ReqTbCreateVoucherBalance struct {
	AccountId            int64     `json:"accountId" db:"account_id"`
	Addr 			           string    `json:"addr" db:"addr"`
	PromotionId          int64     `json:"promotionId" db:"promotion_id"`
	CurrentAmount        uint64     `json:"currentAmount" db:"current_amount"`
	TotalReceivedAmount  uint64     `json:"totalReceivedAmount" db:"total_received_amount"`
}

type ReqTbUpdateVoucherBalance struct {
	CurrentAmount        uint64     `json:"currentAmount" db:"current_amount"`
	TotalReceivedAmount  uint64     `json:"totalReceivedAmount" db:"total_received_amount"`
}

type ReqTbCreateVoucherSendEvent struct {
	AccountId          int64     `json:"accountId" db:"account_id"`
	RecipientAddr      int64     `json:"recipientAddr" db:"recipient_addr"`
	PromotionId        int64     `json:"promotionId" db:"promotion_id"`
	Amount             uint64    `json:"amount" db:"amount"`
	SentAt             time.Time `json:"sentAt" db:"sent_at"`
}

type ReqTbUpdateVoucherSendEvent struct {
	Amount             uint64    `json:"amount" db:"amount"`
}

type ReqTbCreateVoucherBurnEvent struct {
	AccountId             int64     `json:"accountId" db:"account_id"`
	Addr                  int64     `json:"addr" db:"addr"`
	PromotionId           int64     `json:"promotionId" db:"promotion_id"`
	BurnedVoucherAmount   uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount    uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
	BurnedAt              time.Time `json:"burnedAt" db:"burned_at"`
}

type ReqTbUpdateVoucherBurnEvent struct {
	BurnedVoucherAmount   uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount    uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
}

type ReqCreatePromotion struct {
	Title 			           string    `json:"title" db:"title"`
	Desc 			             string    `json:"desc" db:"desc"`
	Url 			             string    `json:"url" db:"url"`
	IsActive               bool      `json:"isActive" db:"is_active"`
	IsWhitelisted          bool      `json:"isWhitelisted" db:"is_whitelisted"`
	VoucherName 			     string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0  int       `json:"voucherExchangeRatio0" db:"voucher_exchange_ratio_0"`
	VoucherExchangeRatio1  int       `json:"voucherExchangeRatio1" db:"voucher_exchange_ratio_1"`
	VoucherTotalSupply     uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	PromotionStartAt       time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt         time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt           time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt             time.Time `json:"claimEndAt" db:"claim_end_at"`
	DistributionPools      []ReqCreateDistPool `json:"distributionPools"`
}

type ReqCreateDistPool struct {
	PrizeDenomId           int64     `json:"prizeDenomId" db:"prize_denom_id"`
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	Prizes                 []ReqCreatePrize `json:"prizes"`
}

type ReqCreatePrize struct {
	Amount            uint64    `json:"amount" db:"amount"`
	Odds              float64   `json:"odds" db:"odds"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	// IsActive          bool      `json:"isActive" db:"is_active"`
}