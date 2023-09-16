package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
)

func QueryOrderById(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).Find(order).Error
	return
}

func UpdateOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).Update(order).Error
	return
}

func QueryOrders(orders *[]schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Find(orders).Error
	return
}

func CreateOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Create(order).Error
	return
}

func QueryOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).First(order).Error
	return
}

func DeleteOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).Delete(order).Error
	return
}

func QueryGameWinningResults(results *[](*types.ResGetGameWinningResults)) (err error) {
	q := "SELECT " +
		"    G.addr, P.title, G.used_ticket_qty,  " +
		"	 D.name as 'prize_name', R.amount as 'prize_amount', G.status, " +
		"	 G.claimed_at, G.claim_finished_at,  " +
		"	 D.type as 'prize_type', R.prize_denom_id, G.prize_id " +
		"  FROM GAME_ORDER G " +
		"   LEFT JOIN (SELECT title, promotion_id FROM promotion) P ON G.promotion_id = P.promotion_id " +
		"   LEFT JOIN prize R ON G.prize_id = R.prize_id " +
		"   LEFT JOIN (SELECT name, type, prize_denom_id FROM prize_denom) D ON R.prize_denom_id = D.prize_denom_id " +
		"  WHERE G.is_win = 1 " +
		"  ORDER BY G.claimed_at DESC, P.title ASC, G.addr ASC "

	if err = config.DB.Raw(q).Scan(results).
		Error; err != nil {
		return
	}
	return
}
