package types

import "time"

type GameOrderStatusRow struct {
	OrderId         int64      `json:"orderId" db:"order_id"`
	Status          int        `json:"status" db:"status"`
	ClaimFinishedAt *time.Time `json:"claimFinishedAt" db:"claim_finished_at" gorm:"default:null"`
}

// game_order -> status: 1(진행중) 2(꽝으로인한종료) 3(클레임전) 4(클레임중) 5(클레임성공) 6(클레임실패) 7(취소)
type GameOrderStatus struct {
	N_A        int
	InProgress int
	Lost       int
	NotClaimed int
	Claiming   int
	Paid       int
	Failed     int
	Cancelled  int
}

var GameOrderStatusInt = map[int]string{
	0: "N/A",
	1: "In Progress",
	2: "Lost",
	3: "Not Claimed",
	4: "Claiming",
	5: "Paid",
	6: "Failed",
	7: "Cancelled",
}

var GameOrderStatusString = map[string]int{
	"N/A":         0,
	"In Progress": 1,
	"Lost":        2,
	"Not Claimed": 3,
	"Claiming":    4,
	"Paid":        5,
	"Failed":      6,
	"Cancelled":   7,
}
