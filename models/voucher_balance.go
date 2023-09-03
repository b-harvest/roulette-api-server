package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryVoucherBalances(bals *[]schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Find(bals).Error
	return
}

func CreateVoucherBalance(bal *schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Create(bal).Error
	return
}

func QueryVoucherBalance(bal *schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Where("id = ?", bal.Id).First(bal).Error
	return
}

func UpdateVoucherBalance(bal *schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Where("id = ?", bal.Id).Update(bal).Error
	return
}

func DeleteVoucherBalance(bal *schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Where("id = ?", bal.Id).Delete(bal).Error
	return
}
