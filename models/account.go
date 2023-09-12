package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
)

func QueryOrCreateAccount(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).FirstOrCreate(acc).Error
	return
}

func QueryAccBalance(accBalance *types.ResGetAccBalance) (err error) {
	sql := `
		SELECT UVB.promotion_id as promotion_id,
			ACC.addr as addr,
			ACC.ticket_amount as ticket_amount,
			UVB.current_amount as voucher_amount,
			UVB.total_reiceved_amount as total_reiceved_voucher_amount
		FROM account as ACC
			JOIN user_voucher_balance as UVB
				ON ACC.addr = UVB.addr
		WHERE ACC.addr = ?;
	`
	err = config.DB.Raw(sql, accBalance.Addr).Scan(accBalance).Error
	fmt.Println(accBalance)
	return
}

func QueryAccounts(accs *[]schema.AccountRow) (err error) {
	err = config.DB.Table("account").Find(accs).Error
	return
}

func CreateAccount(acc *types.ReqTbCreateAccount) (err error) {
	err = config.DB.Table("account").Create(acc).Error
	return
}

func QueryAccount(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).First(acc).Error
	return
}

func UpdateAccount(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).Update(acc).Error
	return
}

func DeleteAccount(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).Delete(acc).Error
	return
}
