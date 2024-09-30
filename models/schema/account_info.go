package schema

import "time"

type AccountInfoRow struct {
	ID               int64     `json:"id" db:"id"`
	Addr             string    `json:"addr" db:"addr"`
	DelegationAmount float64   `json:"delegationAmount" db:"delegation_amount"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
}

func (b *AccountInfoRow) TableName() string {
	return "account_info"
}
