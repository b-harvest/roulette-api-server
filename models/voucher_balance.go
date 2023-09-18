package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

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

func QueryVoucherBalanceByAddrPromotionId(bal *schema.VoucherBalanceRow) (err error) {
	err = config.DB.Table("user_voucher_balance").Where("addr = ? and promotion_id = ?", bal.Addr, bal.PromotionId).First(bal).Error
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

func QueryAvailableVouchers(promotions *[](*types.ResGetAvailableVouchers)) (err error) {
	q :=
		"SELECT promotion_id, title, voucher_name, voucher_total_supply, voucher_remaining_qty " +
			"FROM promotion  " +
			// "WHERE is_active = 1 " +
			// "AND is_whitelisted = 1  " +
			// "AND promotion_end_at > NOW()"
			"WHERE promotion_end_at > NOW()"

	// if err = config.DB.Table("promotion").Order("promotion_start_at DESC").Find(promotions).
	if err = config.DB.Raw(q).Scan(promotions).
		Error; err != nil {
		return
	}
	return
}
