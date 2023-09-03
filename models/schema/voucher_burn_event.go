package schema

import (
	"time"
)

type VoucherBurnEventRow struct {
	Id                    int64     `json:"id" db:"id"`
	AccountId             int64     `json:"accountId" db:"account_id"`
	Addr                  int64     `json:"addr" db:"addr"`
	PromotionId           int64     `json:"promotionId" db:"promotion_id"`
	BurnedVoucherAmount   uint64    `json:"burnedVoucherAmount" db:"burned_voucher_amount"`
	MintedTicketAmount    uint64    `json:"mintedTicketAmount" db:"minted_ticket_amount"`
	BurnedAt              time.Time `json:"burnedAt" db:"burned_at"`
}

func (b *VoucherBurnEventRow) TableName() string {
	return "account"
}
