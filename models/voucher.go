package models

import (
	"errors"
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func Swap(account *schema.Account, addr string, cntVoucher int64) (err error) {
	// query account
	if err = config.DB.Table("account").Where("address = ?", addr).First(account).Error; err != nil {
		fmt.Println(err)
		return err
	}
	if account.Voucher < cntVoucher {
		return errors.New(fmt.Sprintf("not enough vouchers. requested:%v, you have:%v", cntVoucher, account.Voucher))
	}

	// update
	account.Voucher = account.Voucher - cntVoucher
	account.Ticket = account.Ticket + cntVoucher
	fmt.Printf("new voucher: %v\n", account.Voucher)
	fmt.Printf("new ticket: %v\n", account.Ticket)
	// if err = config.DB.Table("account").Where("address = ?", addr).Update(account).Error; err != nil {
	// 	return err
	// }
	config.DB.Table("account").Where("address = ?", addr).Update("voucher", account.Voucher)
	config.DB.Table("account").Where("address = ?", addr).Update("ticket", account.Ticket)
	return nil
}


func SendVoucher(account *schema.Account, addr string, cntVoucher int64) (err error) {
	// query account
	if err = config.DB.Table("account").Where("address = ?", addr).First(account).Error; err != nil {
		fmt.Println(err)
		return err
	}

	// update
	account.Voucher = account.Voucher + cntVoucher
	config.DB.Table("account").Where("address = ?", addr).Update("voucher", account.Voucher)
	return nil
}

func BurnUserTicket(account *schema.Account, addr string, cntTicket int64) (err error) {
	// query account
	if err = config.DB.Table("account").Where("address = ?", addr).First(account).Error; err != nil {
		fmt.Println(err)
		return err
	}

	// update
	account.Ticket = account.Ticket - cntTicket
	config.DB.Table("account").Where("address = ?", addr).Update("ticket", account.Ticket)
	return nil
}
