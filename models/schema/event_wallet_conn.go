package schema

import "time"

type EventwalletConnRow struct {
	Id          uint64    `json:"id" db:"id"`
	Addr        string    `json:"addr" db:"addr"`
	AddrType    string     `json:"type" db:"type"`
	PromotionId int64     `json:"promotionId" db:"promotion_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func (r *EventwalletConnRow) TableName() string {
	return "event_wallet_conn"
}
