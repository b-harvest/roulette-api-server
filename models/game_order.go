package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

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

func UpdateOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).Update(order).Error
	return
}

func DeleteOrder(order *schema.OrderRow) (err error) {
	err = config.DB.Table("game_order").Where("order_id = ?", order.OrderId).Delete(order).Error
	return
}
