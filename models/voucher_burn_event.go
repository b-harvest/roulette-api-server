package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

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

func QueryVoucherBurnEventsByAddr(events *[](*types.ResGetVoucherBurnEvents), addr string) (err error) {
	q := `
		SELECT 
			E.*, P.voucher_name as voucher_name
		FROM voucher_burn_event E
		LEFT JOIN
			(select promotion_id, voucher_name from promotion
			) P ON E.promotion_id = P.promotion_id
		WHERE E.addr = ?
		ORDER BY E.id DESC
		`
	if err = config.DB.Raw(q, addr).Scan(events).
		Error; err != nil {
		return
	}
	return
}