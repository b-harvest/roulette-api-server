package schema

import (
	"time"
)

type VoucherSendEventRow struct {
	Id            int64     `json:"id" db:"id"`
	AccountId     int64     `json:"accountId" db:"account_id"`
	RecipientAddr int64     `json:"recipientAddr" db:"recipient_addr"`
	PromotionId   int64     `json:"promotionId" db:"promotion_id"`
	Amount        uint64    `json:"amount" db:"amount"`
	SentAt        time.Time `json:"sentAt" db:"sent_at"`
}

func (b *VoucherSendEventRow) TableName() string {
	return "voucher_send_event"
}
