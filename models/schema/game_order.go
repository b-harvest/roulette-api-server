package schema

import (
	"time"
)

type OrderRow struct {
	OrderId         int64     `json:"orderId" db:"order_id"`
	AccountId       int64     `json:"accountId" db:"account_id"`
	Addr            string    `json:"addr" db:"addr"`
	PromotionId     int64     `json:"promotionId" db:"promotion_id"`
	GameId          int64     `json:"gameId" db:"game_id"`
	IsWin           bool      `json:"isWin" db:"is_win"`
	Status          int       `json:"status" db:"status"`
	UsedTicketQty   uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId         int64     `json:"prizeId" db:"prize_id"`
	StartedAt       time.Time `json:"startedAt" db:"started_at"`
	ClaimedAt       time.Time `json:"claimedAt" db:"claimed_at" gorm:"column:claimed_at; type:timestamp; default:null"`
	ClaimFinishedAt time.Time `json:"claimFinishedAt" db:"claim_finished_at" gorm:"default:null"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type OrderRowWithID struct {
	ID     			int64     
	OrderId         int64     `json:"orderId" db:"order_id"`
	AccountId       int64     `json:"accountId" db:"account_id"`
	Addr            string    `json:"addr" db:"addr"`
	PromotionId     int64     `json:"promotionId" db:"promotion_id"`
	GameId          int64     `json:"gameId" db:"game_id"`
	IsWin           bool      `json:"isWin" db:"is_win"`
	Status          int       `json:"status" db:"status"`
	UsedTicketQty   uint64    `json:"usedTicketQty" db:"used_ticket_qty"`
	PrizeId         int64     `json:"prizeId" db:"prize_id"`
	StartedAt       time.Time `json:"startedAt" db:"started_at"`
	ClaimedAt       time.Time `json:"claimedAt" db:"claimed_at" gorm:"column:claimed_at; type:timestamp; default:null"`
	ClaimFinishedAt time.Time `json:"claimFinishedAt" db:"claim_finished_at" gorm:"default:null"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

func (r *OrderRow) TableName() string {
	return "game_order"
}

func (r *OrderRowWithID) TableName() string {
	return "game_order"
}
