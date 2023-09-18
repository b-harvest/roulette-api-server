package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func QueryOrCreateAccount(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).FirstOrCreate(acc).Error
	return
}

func QueryBalancesByAddr(bals *[]types.ResGetBalanceByAcc, addr string) (err error) {
	sql := `
	SELECT UVB.promotion_id as promotion_id,
		ACC.addr as addr,
		ACC.ticket_amount as ticket_amount,
		UVB.current_amount as voucher_amount,
		UVB.total_received_amount as total_received_voucher_amount
	FROM account as ACC
		JOIN user_voucher_balance as UVB
			ON ACC.addr = UVB.addr
	WHERE ACC.addr = ?;
	`
	err = config.DB.Raw(sql, addr).Scan(bals).Error
	return
}

func QueryOrdersByAcc(orders *[]schema.OrderRow, addr string) (err error) {
	err = config.DB.Table("game_order").Where("addr = ?", addr).Find(orders).Error
	return
}

func QueryWinTotalByAcc(winTotals *[]types.ResGetWinTotalByAcc, addr string) (err error) {
	sql := `
	SELECT PD.name as prize_name, SUM(GOP.amount) as total_amount
	FROM (
		SELECT GO.order_id, P.prize_id, P.prize_denom_id, P.amount
		FROM (
				SELECT *
				FROM game_order
				WHERE addr = ?
					AND is_win = 1
			) as GO
			JOIN prize as P
				ON GO.prize_id = P.prize_id
		) as GOP
		JOIN prize_denom as PD
			ON GOP.prize_denom_id = PD.prize_denom_id
	GROUP BY GOP.prize_id;
	`
	err = config.DB.Raw(sql, addr).Scan(winTotals).Error
	return
}

func QueryAccountByAddr(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).Find(acc).Error
	return
}

func UpdateAccountById(tx *gorm.DB, acc *schema.AccountRow) error {
	err := tx.Table("account").Where("id = ?", acc.Id).Update(acc).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
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
