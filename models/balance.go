package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryBalanceByAddr(account *schema.Account, addr string) (err error) {
	if err = config.DB.Table("account").Where("address = ?", addr).First(account).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func InsertNewAddr(account *schema.Account, addr string) (err error) {
	account.Address = addr
	account.HexAddr = addr
	if err = config.DB.Table("account").Create(account).Error; err != nil {
		return err
	}
	return nil
}