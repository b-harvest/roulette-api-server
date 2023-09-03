package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryVoucherSendEvents(events *[]schema.VoucherSendEventRow) (err error) {
	err = config.DB.Table("voucher_send_event").Find(events).Error
	return
}

func CreateVoucherSendEvent(event *schema.VoucherSendEventRow) (err error) {
	err = config.DB.Table("voucher_send_event").Create(event).Error
	return
}

func QueryVoucherSendEvent(event *schema.VoucherSendEventRow) (err error) {
	err = config.DB.Table("voucher_send_event").Where("id = ?", event.Id).First(event).Error
	return
}

func UpdateVoucherSendEvent(event *schema.VoucherSendEventRow) (err error) {
	err = config.DB.Table("voucher_send_event").Where("id = ?", event.Id).Update(event).Error
	return
}

func DeleteVoucherSendEvent(event *schema.VoucherSendEventRow) (err error) {
	err = config.DB.Table("voucher_send_event").Where("id = ?", event.Id).Delete(event).Error
	return
}
