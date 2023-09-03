package schema

import (
	"time"
)

type VoucherBalanceRow struct {
	Id                   int64     `json:"id" db:"id"`
	AccountId            int64     `json:"accountId" db:"account_id"`
	Addr 			           string    `json:"addr" db:"addr"`
	PromotionId          int64     `json:"promotionId" db:"promotion_id"`
	CurrentAmount        uint64     `json:"currentAmount" db:"current_amount"`
	TotalReceviedAmount  uint64     `json:"totalReceviedAmount" db:"total_recevied_amount"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`
}

func (b *VoucherBalanceRow) TableName() string {
	return "account"
}
