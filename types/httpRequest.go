package types

type ReqCreateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
	//CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	//UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ReqUpdateGame struct {
	//GameId         int64  `json:"gameId" db:"game_id"`
	Title 			   string `json:"title" db:"title"`
	Desc 			     string `json:"desc" db:"desc"`
	IsActive       bool   `json:"isActive" db:"is_active"`
	Url 			     string `json:"url" db:"url"`
	//CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	//UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}