package schema

import (
	"time"
)

// promotion 조회 시 sub 정보
type PrizeRow struct {
	PrizeId           int64     `json:"prizeId" db:"prize_id"`
	DistPoolId        int64     `json:"distPoolId" db:"dist_pool_id"`
	PromotionId       int64     `json:"promotionId" db:"promotion_id"`
	PrizeDenomId      int64     `json:"prizeDenomId" db:"dist_pool_id"`
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

func (r *PrizeRow) TableName() string {
	return "promotion"
}



