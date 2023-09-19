package types

type ResAccountStat struct {
	TotalAccountNum    uint64 `json:"totalAccountNum" db:"total_account_num" gorm:"column:total_account_num;type:uint64"`
	TotalBlackListNum  uint64 `json:"totalBlackListNum" db:"totla_blacklist_num" gorm:"column:totla_blacklist_num;type:uint64"`
	RegisterAccountNum uint64 `json:"registerAccountNum" db:"register_account_num" gorm:"column:register_account_num;type:uint64"`
	LoginAccountNum    uint64 `json:"loginAccountNum" db:"login_account_num" gorm:"column:login_account_num;type:uint64"`
	PeriodStart        string `json:"periodStart" db:"period_start" gorm:"column:period_start;type:varchar(255)"`
	PeriodEnd          string `json:"periodEnd" db:"period_end" gorm:"column:period_end;type:varchar(255)"`
}

type ResPromotionStat struct {
	InProgressCount uint64 `json:"inProgrressCount" db:"in_progress_count" gorm:"column:in_progress_count;type:uint64"`
	FinishedCount   uint64 `json:"finishedCount" db:"finished_count" gorm:"column:finished_count;type:uint64"`
	NotStartedCount uint64 `json:"notStartedCount" db:"not_started_count" gorm:"column:not_started_count;type:uint64"`
}

type ResFlipLinkStat struct {
	PromotionId uint64                 `json:"promotionId" db:"promotion_id" gorm:"column:promotion_id;type:uint64"`
	Title       string                 `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	TotalCount  uint64                 `json:"totalCount" db:"total_count" gorm:"column:total_count;type:uint64"`
	Daily       []ResFlipLinkDailyStat `json:"daily" db:"daily" gorm:"column:daily;type:json)"`
}

type ResFlipLinkDailyStat struct {
	Date  string `json:"date" db:"date" gorm:"column:date;type:varchar(255)"`
	Count uint64 `json:"count" db:"count" gorm:"column:count;type:uint64"`
}

type ResWalletConnectStat struct {
	PromotionId uint64                      `json:"promotionId" db:"promotion_id" gorm:"column:promotion_id;type:uint64"`
	Title       string                      `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	TotalCount  uint64                      `json:"totalCount" db:"total_count" gorm:"column:total_count;type:uint64"`
	Daily       []ResWalletConnectDailyStat `json:"daily" db:"daily" gorm:"column:daily;type:json)"`
}

type ResWalletConnectDailyStat struct {
	Date  string `json:"date" db:"date" gorm:"column:date;type:varchar(255)"`
	Count uint64 `json:"count" db:"count" gorm:"column:count;type:uint64"`
}

type ResVoucherStat struct {
	PromotionId     uint64 `json:"promotionId" db:"promotion_id" gorm:"column:promotion_id;type:uint64"`
	Title           string `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	TotalSupply     uint64 `json:"totalSupply" db:"total_supply" gorm:"column:total_supply;type:uint64"`
	RemainingQty    uint64 `json:"remainingQty" db:"remaining_qty" gorm:"column:remaining_qty;type:uint64"`
	SentVouchers    uint64 `json:"sentVouchers" db:"sent_vouchers" gorm:"column:sent_vouchers;type:uint64"`
	ReceipientCount uint64 `json:"recipientCount" db:"recipient_count" gorm:"column:recipient_count;type:uint64"`
	BurntVouchers   uint64 `json:"burntVouchers" db:"burnt_vouchers" gorm:"column:burnt_vouchers;type:uint64"`
	MintedTickets   uint64 `json:"mintedTickets" db:"minted_tickets" gorm:"column:minted_tickets;type:uint64"`
}

type ResTicketStat struct {
	PromotionId     uint64               `json:"promotionId" db:"promotion_id" gorm:"column:promotion_id;type:uint64"`
	Title           string               `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	TotalMinted     uint64               `json:"totalMinted" db:"total_minted" gorm:"column:total_minted;type:uint64"`
	TotalUsed       uint64               `json:"totalUsed" db:"total_used" gorm:"column:total_used;type:uint64"`
	TicketUserCount uint64               `json:"ticketUserCount" db:"ticket_user_count" gorm:"column:ticket_user_count;type:uint64"`
	TicketUsage     []ResTicketUsageStat `json:"ticketUsage" db:"ticket_usage" gorm:"column:ticket_usage;type:json)"`
}

type ResTicketUsageStat struct {
	GameId uint64 `json:"gameId" db:"game_id" gorm:"column:game_id;type:uint64"`
	Title  string `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	Used   uint64 `json:"used" db:"used" gorm:"column:used;type:uint64"`
}

type ResPrizeStat struct {
	PromotionId                   uint64  `json:"promotionId" db:"promotion_id" gorm:"column:promotion_id;type:uint64"`
	Title                         string  `json:"title" db:"title" gorm:"column:title;type:varchar(255)"`
	PoolTotalSupplyUsdValue       float64 `json:"poolTotalSupplyUsdValue" db:"pool_total_supply_usd_value" gorm:"column:pool_total_supply_usd_value;type:float64"`
	PrizeMaxTotalWinLimitUsdValue float64 `json:"prizeMaxTotalWinLimitUsdValue" db:"prize_max_total_win_limit_usd_value" gorm:"column:prize_max_total_win_limit_usd_value;type:float64"`
	PaidUsdValue                  float64 `json:"paidUsdValue" db:"paid_usd_value" gorm:"column:paid_usd_value;type:float64"`
	NotClaimedUsdValue            float64 `json:"notClaimedUsdValue" db:"not_claimed_usd_value" gorm:"column:not_claimed_usd_value;type:float64"`
	ClaimingUsdValue              float64 `json:"claimingUsdValue" db:"claiming_usd_value" gorm:"column:claiming_usd_value;type:float64"`
}
