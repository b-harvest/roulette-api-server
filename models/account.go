package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
)

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
