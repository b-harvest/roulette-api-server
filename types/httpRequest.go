package types

import "time"

type ReqCreateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
	//CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	//UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ReqUpdateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
}

type ReqCreatePromotion struct {
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

type ReqUpdatePromotion struct {
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

type ReqCreatePrizeDenom struct {
	Name 			             string    `json:"name" db:"name"`
	Type 			             string    `json:"type" db:"type"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqUpdatePrizeDenom struct {
	Name 			             string    `json:"name" db:"name"`
	Type 			             string    `json:"type" db:"type"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqCreateDistPool struct {
	PromotionId            int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId           int64     `json:"prizeDenomId" db:"prize_denom_id"`
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqUpdateDistPool struct {
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	RemainingQty           uint64    `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqCreatePrize struct {
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

type ReqUpdatePrize struct {
	Odds              float64   `json:"odds" db:"odds"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	IsActive          bool      `json:"isActive" db:"is_active"`
}

type ReqCreateAccount struct {
	Addr 			     string    `json:"addr" db:"addr"`
	TicketAmount   uint64    `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo      string    `json:"adminMemo" db:"admin_memo"`
	Type 			     string    `json:"type" db:"type"`
	IsBlacklisted  bool      `json:"isBlacklisted" db:"is_blacklisted"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ReqUpdateAccount struct {
	TicketAmount   uint64    `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo      string    `json:"adminMemo" db:"admin_memo"`
	Type 			     string    `json:"type" db:"type"`
	IsBlacklisted  bool      `json:"isBlacklisted" db:"is_blacklisted"`
}

type ReqCreateOrder struct {
	AccountId         int64     `json:"accountId" db:"account_id"`
	Addr 			        string    `json:"addr" db:"addr"`
	PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	GameId            int64     `json:"gameId" db:"game_id"`
	Status            int       `json:"status" db:"status"`
	UsedTicketQty     uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
}

type ReqUpdateOrder struct {
	IsWin             bool      `json:"isWin" db:"is_win"`
	Status            int       `json:"status" db:"status"`
	UsedTicketQty     uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
	ClaimedAt         time.Time `json:"claimedAt" db:"claimed_at"`
	ClaimFinishedAt   time.Time `json:"claimFinishedAt" db:"claim_finished_at"`
}

type ReqCreateVoucherBalance struct {
	AccountId            int64     `json:"accountId" db:"account_id"`
	Addr 			           string    `json:"addr" db:"addr"`
	PromotionId          int64     `json:"promotionId" db:"promotion_id"`
	CurrentAmount        uint64     `json:"currentAmount" db:"current_amount"`
	TotalReceviedAmount  uint64     `json:"totalReceviedAmount" db:"total_recevied_amount"`
}

type ReqUpdateVoucherBalance struct {
	CurrentAmount        uint64     `json:"currentAmount" db:"current_amount"`
	TotalReceviedAmount  uint64     `json:"totalReceviedAmount" db:"total_recevied_amount"`
}

type ReqCreateVoucherSendEvent struct {
	AccountId          int64     `json:"accountId" db:"account_id"`
	RecipientAddr      int64     `json:"recipientAddr" db:"recipient_addr"`
	PromotionId        int64     `json:"promotionId" db:"promotion_id"`
	Amount             uint64    `json:"amount" db:"amount"`
	SentAt             time.Time `json:"sentAt" db:"sent_at"`
}

type ReqUpdateVoucherSendEvent struct {
	Amount             uint64    `json:"amount" db:"amount"`
}

type ReqCreateVoucherBurnEvent struct {
	AccountId             int64     `json:"accountId" db:"account_id"`
	Addr                  int64     `json:"addr" db:"addr"`
	PromotionId           int64     `json:"promotionId" db:"promotion_id"`
	BurnedVoucherAmount   uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount    uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
	BurnedAt              time.Time `json:"burnedAt" db:"burned_at"`
}

type ReqUpdateVoucherBurnEvent struct {
	BurnedVoucherAmount   uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount    uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
}





