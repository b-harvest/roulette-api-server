package types

import (
	"time"
)

type ResGetPromotions struct {
	PromotionId           uint64    `json:"promotionId" db:"promotion_id"`
	Title                 string    `json:"title" db:"title"`
	Desc                  string    `json:"desc" db:"desc"`
	Url                   string    `json:"url" db:"url"`
	IsActive              bool      `json:"isActive" db:"is_active"`
	IsWhitelisted         bool      `json:"isWhitelisted" db:"is_whitelisted"`
	VoucherName           string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0 int       `json:"voucherExchangeRatio0"`
	VoucherExchangeRatio1 int       `json:"voucherExchangeRatio1"`
	VoucherTotalSupply    uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	VoucherRemainingQty   uint64    `json:"voucherRemainingQty" db:"voucher_remaining_qty"`
	PromotionStartAt      time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt        time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt          time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt            time.Time `json:"claimEndAt" db:"claim_end_at"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
	// ---------- additional
	ParticipantCnt uint64 `json:"voucherReceivedUserCnt" db:"participant_cnt"`
	// status: not started / in progress / finished
	Status            string           `json:"status" db:"status"`
	DistributionPools *[]PrizeDistPool `json:"distributionPools"`
}

type ResGetPromotion struct {
	PromotionId           uint64    `json:"promotionId" db:"promotion_id"`
	Title                 string    `json:"title" db:"title"`
	Desc                  string    `json:"desc" db:"desc"`
	Url                   string    `json:"url" db:"url"`
	IsActive              bool      `json:"isActive" db:"is_active"`
	IsWhitelisted         bool      `json:"isWhitelisted" db:"is_whitelisted"`
	VoucherName           string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0 int       `json:"voucherExchangeRatio0"`
	VoucherExchangeRatio1 int       `json:"voucherExchangeRatio1"`
	VoucherTotalSupply    uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	VoucherRemainingQty   uint64    `json:"voucherRemainingQty" db:"voucher_remaining_qty"`
	PromotionStartAt      time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt        time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt          time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt            time.Time `json:"claimEndAt" db:"claim_end_at"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
	// status: not started / in progress / finished
	ParticipantCnt    uint64                 `json:"voucherReceivedUserCnt" db:"participant_cnt"`
	Status            string                 `json:"status" db:"status"`
	DistributionPools *[]PrizeDistPoolDetail `json:"distributionPools"`
	Summary           PromotionSummary       `json:"promotionSummary"`
}

// TBD
type ResGetLatestOrderByAddr struct {
	OrderId         int64              `json:"orderId" db:"order_id"`
	Addr            string             `json:"addr" db:"addr"`
	PromotionId     int64              `json:"promotionId" db:"promotion_id"`
	GameId          int64              `json:"gameId" db:"game_id"`
	IsWin           bool               `json:"isWin" db:"is_win"`
	Status          int                `json:"status" db:"status"`
	UsedTicketQty   uint64             `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId         int64              `json:"prizeId" db:"prize_id"`
	StartedAt       time.Time          `json:"startedAt" db:"started_at"`
	ClaimedAt       time.Time          `json:"claimedAt" db:"claimed_at" gorm:"column:claimed_at; type:timestamp; default:null"`
	ClaimFinishedAt time.Time          `json:"claimFinishedAt" db:"claim_finished_at" gorm:"default:null"`
	CreatedAt       time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time          `json:"updatedAt" db:"updated_at"`
	Prize           ResOrderPrize      `json:"prize"`
	Promotion       ResSimplePromotion `json:"promotion"`
	RemainingTime   time.Duration      `json:"remainingTime" db:"remaining_time"`
	IsClaimable     bool               `json:"isClaimable" db:"is_claimable"`
}

type ResSimplePromotion struct {
	Title                 string    `json:"title" db:"title"`
	VoucherName           string    `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0 int       `json:"voucherExchangeRatio0"`
	VoucherExchangeRatio1 int       `json:"voucherExchangeRatio1"`
	PromotionStartAt      time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	PromotionEndAt        time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	ClaimStartAt          time.Time `json:"claimStartAt" db:"claim_start_at"`
	ClaimEndAt            time.Time `json:"claimEndAt" db:"claim_end_at"`
	Status                string    `json:"status" db:"status"`
}

type ResOrderPrize struct {
	PrizeId      int64              `json:"prizeId" db:"prize_id"`
	PrizeDenomId int64              `json:"prizeDenomId" db:"dist_pool_id"`
	Amount       uint64             `json:"amount" db:"amount"`
	WinImageUrl  string             `json:"winImageUrl" db:"win_image_url"`
	PrizeDenom   ResOrderPrizeDenom `json:"prizeDenom"`
}

type ResOrderPrizeDenom struct {
	PrizeDenomId int64  `json:"prizeDenomId" db:"prize_denom_id"`
	Name         string `json:"name" db:"name"`
	Type         string `json:"type" db:"type"`
}

type ResGetBalanceByAcc struct {
	PromotionID                uint64 `json:"promotionId" db:"promotion_id"`
	Addr                       string `json:"addr" db:"addr"`
	TicketAmount               uint64 `json:"ticketAmount" db:"ticket_amount"`
	VoucherAmount              uint64 `json:"voucherAmount" db:"voucher_amount"`
	TotalreceivedVoucherAmount uint64 `json:"totalreceivedVoucherAmount" db:"total_received_voucher_amount"`
}

type ResGetWinTotalByAcc struct {
	PrizeName   string `json:"prizeName" db:"prize_name"`
	TotalAmount int    `json:"totalAmount" db:"total_amount"`
}

type ResGetEventCount struct {
	Cnt uint64 `json:"cnt" db:"cnt"`
}

type ResGetVoucherSendEvents struct {
	Id            int64     `json:"id" db:"id"`
	AccountId     int64     `json:"accountId" db:"account_id"`
	RecipientAddr string    `json:"recipientAddr" db:"recipient_addr"`
	PromotionID   uint64    `json:"promotionId" db:"promotion_id"`
	VoucherName   string    `json:"voucherName" db:"voucher_name"`
	Amount        uint64    `json:"amount" db:"amount"`
	SentAt        time.Time `json:"sentAt" db:"sent_at"`
}

type ResGetAvailableVouchers struct {
	PromotionId         uint64 `json:"promotionId" db:"promotion_id"`
	Title               string `json:"title" db:"title"`
	VoucherName         string `json:"voucherName" db:"voucher_name"`
	VoucherTotalSupply  uint64 `json:"voucherTotalSupply" db:"voucher_total_supply"`
	VoucherRemainingQty uint64 `json:"voucherRemainingQty" db:"voucher_remaining_qty"`
}

type ResGetGameWinningResults struct {
	OrderId         uint64    `json:"orderId" db:"order_id"`
	AccountId       int64     `json:"accountId" db:"account_id"`
	Addr            string    `json:"addr" db:"addr"`
	Title           string    `json:"title" db:"title"`
	UsedTicketQty   uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeName       string    `json:"prizeName" db:"prize_name"`
	PrizeAmount     uint64    `json:"prizeAmount" db:"prize_amount"`
	Status          string    `json:"status" db:"status"`
	ClaimedAt       time.Time `json:"claimedAt" db:"claimed_at"`
	ClaimFinishedAt time.Time `json:"claimFinishedAt" db:"claim_finished_at"`
	PrizeType       string    `json:"prizeType" db:"prize_type"`
	PrizeDenomId    uint64    `json:"prizeDenomId" db:"prize_denom_id"`
	PrizeId         uint64    `json:"prizeId" db:"prize_id"`
}

type ResGetAccount struct {
	Id int64 `json:"id" db:"id"`
	// UserId        int64                `json:"userId" db:"user_id"`
	Addr          string               `json:"addr" db:"addr"`
	TicketAmount  uint64               `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo     string               `json:"adminMemo" db:"admin_memo"`
	Type          string               `json:"type" db:"type"`
	IsBlacklisted bool                 `json:"isBlacklisted" db:"is_blacklisted"`
	LastLoginAt   time.Time            `json:"lastLoginAt" db:"last_login_at" gorm:"default:null"`
	CreatedAt     time.Time            `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time            `json:"updatedAt" db:"updated_at"`
	Vouchers      []*ResGetVoucher     `json:"vouchers"`
	Summary       ResGetAccountSummary `json:"summary"`
}

type ResGetAccountSummary struct {
	TotalWinUsd             float64                `json:"totalWinUsd" db:"total_win_usd"`
	TotalClaimbleUsd        float64                `json:"totalClaimbleUsd" db:"total_claimble_usd"`
	TotalCurrentVoucherNum  uint64                 `json:"totalCurrentVoucherCnt" db:"total_current_voucher_num"`
	TotalReceivedVoucherNum uint64                 `json:"totalReceivedVoucherCnt" db:"total_received_voucher_num"`
	TotalConnectNum         uint64                 `json:"totalConnectCnt" db:"total_connect_num"`
	TotalOrderNum           uint64                 `json:"totalOrderCnt" db:"total_order_num"`
	TotalWinNum             uint64                 `json:"totalWinOrderCnt" db:"total_win_num"`
	TotalClaimbleNum        uint64                 `json:"totalClaimbleOrderCnt" db:"total_claimble_num"`
	WinPrizes               *[]ResGetWinTotalByAcc `json:"winPrizes"`
}

type ResGetVoucher struct {
	CurrentAmount       uint64                 `json:"currentAmount" db:"current_amount"`
	TotalReceivedAmount uint64                 `json:"totalReceivedAmount" db:"total_received_amount"`
	CreatedAt           time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time              `json:"updatedAt" db:"updated_at"`
	PromotionId         int64                  `json:"promotionId" db:"promotion_id"`
	Promotion           ResGetVoucherPromotion `json:"promotion"`
	// Id                  int64     `json:"id" db:"id"`
	// AccountId           int64     `json:"accountId" db:"account_id"`
	// Addr                string    `json:"addr" db:"addr"`
}

type ResGetVoucherPromotion struct {
	Title                 string `json:"title" db:"title"`
	Desc                  string `json:"desc" db:"desc"`
	VoucherName           string `json:"voucherName" db:"voucher_name"`
	VoucherExchangeRatio0 int    `json:"voucherExchangeRatio0"`
	VoucherExchangeRatio1 int    `json:"voucherExchangeRatio1"`
	// ID                    int64
	// PromotionId           int64     `json:"promotionId" db:"promotion_id"`
	// Url                   string    `json:"url" db:"url"`
	// IsActive              bool      `json:"isActive" db:"is_active" gorm:"column:is_active"`
	// IsWhitelisted         bool      `json:"isWhitelisted" db:"is_whitelisted" gorm:"column:is_whitelisted"`
	// VoucherTotalSupply    uint64    `json:"voucherTotalSupply" db:"voucher_total_supply"`
	// VoucherRemainingQty   uint64    `json:"voucherRemainingQty" db:"voucher_remaining_qty"`
	// PromotionStartAt      time.Time `json:"promotionStartAt" db:"promotion_start_at"`
	// PromotionEndAt        time.Time `json:"promotionEndAt" db:"promotion_end_at"`
	// ClaimStartAt          time.Time `json:"claimStartAt" db:"claim_start_at"`
	// ClaimEndAt            time.Time `json:"claimEndAt" db:"claim_end_at"`
	// CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	// UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}

type ResTransfersHistoryByAddr struct {
	Addr              string                     `json:"addr"`
	VoucherSendEvents []*ResGetVoucherSendEvents `json:"voucherSendEvents" db:"voucher_send_events"`
	VoucherBurnEvents []*ResGetVoucherBurnEvents `json:"voucherBurnEvents" db:"voucher_burn_events"`
}

type ResGetVoucherBurnEvents struct {
	Id                  int64     `json:"id" db:"id"`
	AccountId           int64     `json:"accountId" db:"account_id"`
	Addr                string    `json:"addr" db:"addr"`
	PromotionId         int64     `json:"promotionId" db:"promotion_id"`
	VoucherName         string    `json:"voucherName" db:"voucher_name"`
	BurnedVoucherAmount uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount  uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
	BurnedAt            time.Time `json:"burnedAt" db:"burned_at"`
}

type ResPostVoucherBurn struct {
	PromotionID   int64  `json:"promotionId"`
	Addr          string `json:"addr"`
	VoucherAmount uint64 `json:"voucher_amount"`
	TicketAmount  uint64 `json:"ticket_amount"`
}

type ResStartGame struct {
	OrderId       int64     `json:"orderId" db:"order_id"`
	AccountId     int64     `json:"accountId" db:"account_id"`
	Addr          string    `json:"addr" db:"addr"`
	PromotionId   int64     `json:"promotionId" db:"promotion_id"`
	GameId        int64     `json:"gameId" db:"game_id"`
	Status        int       `json:"status" db:"status"`
	UsedTicketQty uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	StartedAt     time.Time `json:"startedAt" db:"started_at"`
}

type ResAllClaim struct {
	Addr            string `json:"addr"`
	NumClaimedOrder int    `json:"numClaimedOrder"`
	Status          int    `json:"status"`
}

type ResHealthcheck struct {
	Status string `json:"status"`
}
