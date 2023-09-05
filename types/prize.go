package types

import "time"

// 프로모션 조회 시 서브 정보
type PrizeDistPool struct {
	DistPoolId             int64     `json:"distPoolId" db:"dist_pool_id"`
	// PromotionId            int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId           int64     `json:"prizeDenomId" db:"prize_denom_id"`
	TotalSupply            uint64    `json:"totalSupply" db:"total_supply"`
	RemainingQty           uint64    `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool      `json:"isActive" db:"is_active"`
	// CreatedAt              time.Time `json:"createdAt" db:"created_at"`
	// UpdatedAt              time.Time `json:"updatedAt" db:"updated_at"`
	
	Name string `json:"prizeDenomName" db:"name" gorm:"column:name"`
	Type string `json:"prizeDenomType" db:"type" gorm:"column:type"`
	Prizes *[]Prize `json:"prizes"`
}

// 프로모션 조회 시 서브 정보 Detail
type PrizeDistPoolDetail struct {
	DistPoolId             int64          `json:"distPoolId" db:"dist_pool_id"`
	// PromotionId            int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId           int64          `json:"prizeDenomId" db:"prize_denom_id"`
	TotalSupply            uint64         `json:"totalSupply" db:"total_supply"`
	RemainingQty           uint64         `json:"remainingQty" db:"remaining_qty"`
	IsActive               bool           `json:"isActive" db:"is_active"`
	CreatedAt              time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt              time.Time      `json:"updatedAt" db:"updated_at"`
	Name                   string         `json:"prizeDenomName" db:"name" gorm:"column:name"`
	Type                   string         `json:"prizeDenomType" db:"type" gorm:"column:type"`
	Prizes                 *[]PrizeDetail `json:"prizes"`
}

// 프로모션 조회 시 서브 정보
type Prize struct {
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
	// DistPoolId        int64     `json:"distPoolId" db:"dist_pool_id"`
	// PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	// PrizeDenomId      int64     `json:"prizeDenomId" db:"dist_pool_id"`
	Amount            uint64    `json:"amount" db:"amount"`
	Odds              float64   `json:"odds" db:"odds"`
	WinCnt            uint64    `json:"winCnt" db:"win_cnt"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	IsActive          bool      `json:"isActive" db:"is_active"`
	// CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	// UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

// 프로모션 조회 시 서브 정보 Detail
type PrizeDetail struct {
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
	// DistPoolId        int64     `json:"distPoolId" db:"dist_pool_id"`
	// PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	// PrizeDenomId      int64     `json:"prizeDenomId" db:"dist_pool_id"`
	Amount            uint64    `json:"amount" db:"amount"`
	Odds              float64   `json:"odds" db:"odds"`
	WinCnt            uint64    `json:"winCnt" db:"win_cnt"`
	WinImageUrl 			string    `json:"winImageUrl" db:"win_image_url"`
	MaxDailyWinLimit  uint64    `json:"maxDailyWinLimit" db:"max_daily_win_limit"`
	MaxTotalWinLimit  uint64    `json:"maxTotalWinLimit" db:"max_total_win_limit"`
	IsActive          bool      `json:"isActive" db:"is_active"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}