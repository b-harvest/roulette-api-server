package schema

import "time"

type EventFlipLinkRow struct {
	Id          uint64    `json:"id" db:"id"`
	Addr        string    `json:"addr" db:"addr"`
	AddrType    string    `json:"addrType" db:"addr_type"`
	PromotionId int64     `json:"promotionId" db:"promotion_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func (r *EventFlipLinkRow) TableName() string {
	return "event_flip_link"
}
