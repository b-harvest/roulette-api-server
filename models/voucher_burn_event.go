package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func QueryVoucherBurnEvents(events *[]schema.VoucherBurnEventRow) (err error) {
	err = config.DB.Table("voucher_burn_event").Find(events).Error
	return
}

func CreateVoucherBurnEvent(tx *gorm.DB, event *schema.VoucherBurnEventRow) (err error) {
	if tx == nil {
		tx = config.DB
	}

	err = config.DB.Table("voucher_burn_event").Create(event).Error
	if err != nil {
		tx.Rollback()
		return
	}
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
