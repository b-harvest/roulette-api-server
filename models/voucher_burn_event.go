package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryVoucherBurnEvents(events *[]schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Find(events).Error
	return
}

func CreateVoucherBurnEvent(event *schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Create(event).Error
	return
}

func QueryVoucherBurnEvent(event *schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Where("id = ?", event.Id).First(event).Error
	return
}

func UpdateVoucherBurnEvent(event *schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Where("id = ?", event.Id).Update(event).Error
	return
}

func DeleteVoucherBurnEvent(event *schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Where("id = ?", event.Id).Delete(event).Error
	return
}
