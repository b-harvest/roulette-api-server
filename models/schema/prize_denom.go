package schema

import (
	"time"
)

type PrizeDenomRow struct {
	PrizeDenomId int64     `json:"prizeDenomId" db:"prize_denom_id"`
	Name         string    `json:"name" db:"name"`
	Type         string    `json:"type" db:"type"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

func (r *PrizeDenomRow) TableName() string {
	return "prize_denom"
}
