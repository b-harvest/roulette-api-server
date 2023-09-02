package schema

import "time"

// import "time"

type GameOrder struct {
	GameOrderId   int64 `gorm:"not null" json:"gameOrderId" db:"game_order_id"`
	Address 			string `gorm:"not null" json:"address"`
	PaidTicketNum int64 `gorm:"not null" json:"paidTicketNum" db:"paidTicketNum"`
	Type          int64 `gorm:"not null" json:"type" db:"type"`
	Status        int64 `gorm:"not null" json:"status" db:"status"`
	IsWin         bool `gorm:"not null" json:"isWin" db:"is_win"`
	PrizeID       int64 `json:"prizeID" db:"prize_id"`
	// CreatedAt   time.Time `sql:"DEFAULT:current_timestamp"`
	// UpdatedAt   time.Time
}

type GameInProgress struct {
	GameOrderId   int64 `gorm:"not null" json:"gameOrderId" db:"game_order_id"`
	Address 			string `gorm:"not null" json:"address"`
	PaidTicketNum int64 `gorm:"not null" json:"paidTicketNum" db:"paidTicketNum"`
	Type          int64 `gorm:"not null" json:"type" db:"type"`
	Status        int64 `gorm:"not null" json:"status" db:"status"`
	//IsWin        bool `gorm:"not null" json:"isWin" db:"is_win"`
	// CreatedAt   time.Time `sql:"DEFAULT:current_timestamp"`
	// UpdatedAt   time.Time
}

type Game struct {
	GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

func (b *Game) TableName() string {
	return "game"
}
