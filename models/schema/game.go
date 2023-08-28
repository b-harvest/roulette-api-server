package schema

// import "time"

type Game struct {
	GameOrderId   int64 `gorm:"not null" json:"gameOrderId" db:"game_order_id"`
	Address 			string `gorm:"not null" json:"address"`
	PaidTicketNum int64 `gorm:"not null" json:"paidTicketNum" db:"paidTicketNum"`
	Type          int64 `gorm:"not null" json:"type" db:"type"`
	Status        int64 `gorm:"not null" json:"status" db:"status"`
	IsWin        bool `gorm:"not null" json:"isWin" db:"is_win"`
	GiftId        int64 `json:"giftId" db:"gift_id"`
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
	//GiftId        int64 `json:"giftId" db:"gift_id"`
	// CreatedAt   time.Time `sql:"DEFAULT:current_timestamp"`
	// UpdatedAt   time.Time
}

func (b *Game) TableName() string {
	return "game_order"
}
