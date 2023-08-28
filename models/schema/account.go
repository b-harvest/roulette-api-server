package schema

// import "time"

type Account struct {
	Address 			string `gorm:"not null" json:"address"`
	HexAddr       string `gorm:"not null" json:"hexAddr" db:"hex_addr"`
	Voucher 			int64 `gorm:"not null" json:"voucher" db:"voucher"`
	Ticket 			  int64 `gorm:"not null"  json:"ticket" db:"ticket"`
	// FirstName 	string `gorm:"not null"`
	// LastName 	string `gorm:"not null"`
	// Email   	string `gorm:"unique;not null"`
	// Password   	string `gorm:"not null"`
	// Phone   	string `gorm:"not null"`
	// Address 	*string
	// CreatedAt   time.Time `sql:"DEFAULT:current_timestamp"`
	// UpdatedAt   time.Time
}

func (b *Account) TableName() string {
	return "account"
}
