package schema

import (
	"time"
)

// sample
type Account struct {
	// origin
	Ticket  int64  `json:"ticket" db:"ticket"`
	Voucher int64  `json:"voucher" db:"voucher"`
	Address string `json:"address" db:"address"`
	HexAddr string `json:"hexAddr" db:"hexAddr"`
}

type AccountRow struct {
	Id            int64     `json:"id" db:"id"`
	UserId        int64     `json:"userId" db:"user_id"`
	Addr          string    `json:"addr" db:"addr"`
	TicketAmount  uint64    `json:"ticketAmount" db:"ticket_amount"`
	AdminMemo     string    `json:"adminMemo" db:"admin_memo"`
	Type          string    `json:"type" db:"type"`
	IsBlacklisted bool      `json:"isBlacklisted" db:"is_blacklisted"`
	LastLoginAt   time.Time `json:"lastLoginAt" db:"last_login_at"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

func (b *AccountRow) TableName() string {
	return "account"
}
