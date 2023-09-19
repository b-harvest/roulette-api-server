package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/types"
)

func QueryAccountStat(startDate string, endDate string) (stat types.ResAccountStat, err error) {

	condition1 := ""
	condition2 := ""

	if startDate != "" && endDate != "" {
		condition1 = fmt.Sprintf(" WHERE created_at BETWEEN '%s' AND '%s'", startDate, endDate)
		condition2 = fmt.Sprintf(" WHERE last_login_at BETWEEN '%s' AND '%s'", startDate, endDate)
		stat.PeriodStart = startDate
		stat.PeriodEnd = endDate
	}

	sql := `
		SELECT
			COUNT(id) AS total_account_num,
			SUM(CASE WHEN is_blacklisted = true THEN 1 ELSE 0 END) AS total_blacklist_num,
			R.register_account_num,
			L.login_account_num
		FROM
			account
		LEFT JOIN (
			SELECT COUNT(id) AS register_account_num
			FROM account
			%s
		) AS R ON 1=1
		LEFT JOIN (
			SELECT COUNT(last_login_at) AS login_account_num
			FROM account
			%s
		) AS L ON 1=1
		;
	`

	sql = fmt.Sprintf(sql, condition1, condition2)
	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}
	return
}

func QueryPromotionStat() (stat types.ResPromotionStat, err error) {

	sql := `
	SELECT
		SUM(CASE WHEN now() BETWEEN promotion_start_at AND promotion_end_at THEN 1 ELSE 0 END) AS in_progress_count,
		SUM(CASE WHEN now() > promotion_end_at THEN 1 ELSE 0 END) AS finished_count,
		SUM(CASE WHEN now() < promotion_start_at THEN 1 ELSE 0 END) AS not_started_count
	FROM
		promotion
		;
	`

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}
	return
}

func QueryFlipLinkStat(promotionId, startDate, endDate string) (stat types.ResFlipLinkStat, err error) {

	condition := ""
	if startDate != "" && endDate != "" {
		condition = fmt.Sprintf(" AND created_at BETWEEN '%s' AND '%s'", startDate, endDate)
	}

	sql := `
	SELECT
		e.promotion_id,
		P.title,
		COUNT(e.id) as total_count
	FROM
		event_flip_link e
	LEFT JOIN (SELECT promotion_id, title FROM promotion ) P ON e.promotion_id = P.promotion_id
	WHERE e.promotion_id = 
	` + promotionId + `
	` + condition

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	SELECT
		DATE_FORMAT(created_at, '%Y-%m-%d') AS date,
		COUNT(id) AS count
	FROM
		event_flip_link
	WHERE promotion_id =
	` + promotionId + `
	` + condition + `
	GROUP BY
		DATE_FORMAT(created_at, '%Y-%m-%d');	
	`

	var dailyStats []types.ResFlipLinkDailyStat
	err = config.DB.Raw(sql).Scan(&dailyStats).Error
	if err != nil {
		return
	}

	stat.Daily = dailyStats
	return
}

func QueryWalletConnectStat(promotionId, startDate, endDate string) (stat types.ResWalletConnectStat, err error) {

	condition := ""
	if startDate != "" && endDate != "" {
		condition = fmt.Sprintf(" AND created_at BETWEEN '%s' AND '%s'", startDate, endDate)
	}

	sql := `
	SELECT
		e.promotion_id,
		P.title,
		COUNT(e.id) as total_count
	FROM
		event_wallet_conn e
	LEFT JOIN (SELECT promotion_id, title FROM promotion ) P ON e.promotion_id = P.promotion_id
	WHERE e.promotion_id = 
	` + promotionId + `
	` + condition

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	SELECT
		DATE_FORMAT(created_at, '%Y-%m-%d') AS date,
		COUNT(id) AS count
	FROM
	event_wallet_conn
	WHERE promotion_id =
	` + promotionId + `
	` + condition + `
	GROUP BY
		DATE_FORMAT(created_at, '%Y-%m-%d');	
	`

	var dailyStats []types.ResWalletConnectDailyStat
	err = config.DB.Raw(sql).Scan(&dailyStats).Error
	if err != nil {
		return
	}

	stat.Daily = dailyStats

	return
}

func QueryVoucherStat(promotionId string) (stat types.ResVoucherStat, err error) {

	sql := `
		select 
		promotion_id, title, 
		voucher_total_supply as total_supply, voucher_remaining_qty as remaining_qty
		from promotion p 
		where promotion_id = 
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	select 
		sum(amount) as sent_vouchers,
		count(distinct recipient_addr) as recipient_count
	from 
		voucher_send_event
	where 
		promotion_id = 
		` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	select 
		sum(burned_voucher_amount) as burnt_vouchers,
		sum(minted_ticket_amount) as minted_tickets
	from 
		voucher_burn_event
	where promotion_id = 
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	return
}

func QueryTicketStat(promotionId string) (stat types.ResTicketStat, err error) {

	sql := `
		SELECT
			P.promotion_id,
			P.title,
			SUM(V.minted_ticket_amount) AS total_minted
		FROM
			voucher_burn_event V
		LEFT JOIN
			promotion P ON V.promotion_id = P.promotion_id
		WHERE
			V.promotion_id = 
		` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
		select 
			sum(used_ticket_qty) as total_used,
			count(distinct addr) as ticket_user_count 
		from game_order 
		where promotion_id = 
		` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
		SELECT
			go.game_id,
			gt.title,
			sum(used_ticket_qty) as used     
		FROM
			game_order go
		JOIN
			game_type gt ON go.game_id = gt.game_id
		WHERE
			go.promotion_id = 
		` + promotionId + `
		GROUP BY
			go.game_id, gt.title	
		`

	var usageStats []types.ResTicketUsageStat
	err = config.DB.Raw(sql).Scan(&usageStats).Error
	if err != nil {
		return
	}

	stat.TicketUsage = usageStats

	return
}

func QueryPrizeStat(promotionId string) (stat types.ResPrizeStat, err error) {

	sql := `
		select
			promotion_id, title
		from promotion
		where promotion_id = 
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	select
		SUM(p.amount * p.max_total_win_limit * pd.usd_price) as prize_max_total_win_limit_usd_value
	from 
		prize p
	left join 
		prize_denom pd on p.prize_denom_id = pd.prize_denom_id 
	where 
		promotion_id = 		
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	select 
		SUM(dp.total_supply * pd.usd_price) as pool_total_supply_usd_value
	from 
		distribution_pool dp 
	left join 
		prize_denom pd on dp.prize_denom_id = pd.prize_denom_id 	
	where 
		promotion_id =  		
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	sql = `
	SELECT
		SUM(CASE WHEN g.status = 3 THEN pd.usd_price * p.amount ELSE 0 END) AS not_claimed_usd_value,
		SUM(CASE WHEN g.status = 4 THEN pd.usd_price * p.amount ELSE 0 END) AS claiming_usd_value,
		SUM(CASE WHEN g.status = 5 THEN pd.usd_price * p.amount ELSE 0 END) AS paid_usd_value
	FROM
		game_order g
	LEFT JOIN
		prize p ON g.prize_id = p.prize_id
	LEFT JOIN
		prize_denom pd ON pd.prize_denom_id = p.prize_denom_id
	WHERE
		g.is_win = 1
		AND g.promotion_id =  		
	` + promotionId

	err = config.DB.Raw(sql).Scan(&stat).Error
	if err != nil {
		return
	}

	return
}
