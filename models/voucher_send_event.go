package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
)

func QueryVoucherSendEvents(events *[](*types.ResGetVoucherSendEvents)) (err error) {
	q :=
		"SELECT E.*, P.voucher_name as voucher_name from voucher_send_event E " +
			"LEFT JOIN " +
			"  (select promotion_id, voucher_name from promotion " +
			"   ) P ON E.promotion_id = P.promotion_id"

	if err = config.DB.Raw(q).Scan(events).
		Error; err != nil {
		return
	}
	return
}

//---

func QueryTbVoucherSendEvents(events *[]schema.VoucherSendEventRow) (err error) {
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
