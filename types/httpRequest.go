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
	RemainingQty           uint64    `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}

type ReqUpdateDistPool struct {
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	RemainingQty           uint64    `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool      `json:"isActive" db:"is_active"`
}