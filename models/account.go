package models

import (
	"errors"
	"regexp"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func isEthHexAddress(ethHexAddress string) bool {
	hexAddressRegex := regexp.MustCompile("^(0x)?[0-9a-fA-F]{40}$")
	if !hexAddressRegex.MatchString(ethHexAddress) {
		return false
	}

	if strings.ToLower(ethHexAddress) == ethHexAddress || strings.ToUpper(ethHexAddress) == ethHexAddress {
		return true
	}

	isChecksumAddress := func(address string) bool {
		for i := 2; i < len(address); i++ {
			if address[i] >= 'A' && address[i] <= 'F' {
				return true
			}
		}
		return false
	}

	return isChecksumAddress(ethHexAddress)
}

func CreateAccountWithTx(tx *gorm.DB, acc *schema.AccountRow) error {
	if tx == nil {
		tx = config.DB
	}

	// If ethereum address
	// then check correct ethereum address format
	if acc.Type == "ETH" {
		if !isEthHexAddress(acc.Addr) {
			err := errors.New("Invalid ethereum hex address")
			return err
		}
	}

	err := config.DB.Table("account").Create(acc).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}
func QueryOrCreateAccount(acc *schema.AccountRow) (err error) {
	// If ethereum address
	// then check correct ethereum address format
	if acc.Type == "ETH" {
		if !isEthHexAddress(acc.Addr) {
			err = errors.New("Invalid ethereum hex address")
			return
		}
	}

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
	GROUP BY GOP.prize_denom_id;
	`
	err = config.DB.Raw(sql, addr).Scan(winTotals).Error
	return
}

func QueryAccountById(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("id = ?", acc.Id).Find(acc).Error
	return
}

func QueryAccountByAddr(acc *schema.AccountRow) (err error) {
	err = config.DB.Table("account").Where("addr = ?", acc.Addr).Find(acc).Error
	return
}
func QueryAccountByAddrWithTx(tx *gorm.DB, acc *schema.AccountRow) (bool, error) {
	if tx == nil {
		tx = config.DB
	}

	err := config.DB.Table("account").Where("addr = ?", acc.Addr).Find(acc).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		} else {
			tx.Rollback()
			return false, err
		}
	}
	return true, nil
}

func UpdateAccountById(tx *gorm.DB, acc *schema.AccountRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := tx.Table("account").Where("id = ?", acc.Id).Update(acc).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func UpdateAccountTicketById(tx *gorm.DB, acc *schema.AccountRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := tx.Table("account").Select("ticket_amount").Where("id = ?", acc.Id).Update("ticket_amount", acc.TicketAmount).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func UpdateAccountGoldTicketById(tx *gorm.DB, acc *schema.AccountRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := tx.Table("account").Select("gold_ticket_amount").Where("id = ?", acc.Id).Update("gold_ticket_amount", acc.GoldTicketAmount).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func UpdateAccountTicketByAddr(tx *gorm.DB, acc *schema.AccountRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := tx.Table("account").Select("ticket_amount").Where("addr = ?", acc.Addr).Update("ticket_amount", acc.TicketAmount).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
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

func QueryAndLockAccountByAddr(tx *gorm.DB, account *schema.AccountRow) (bool, error) {
	if tx == nil {
		tx = config.DB
	}

	sql := "SELECT * FROM account WHERE addr = ? FOR UPDATE NOWAIT"

	err := tx.Raw(sql, account.Addr).Scan(account).Error
	if err != nil {
		tx.Rollback()

		// When specific account is locked
		// then rollback and just return (true, nil)
		if err.Error() == "Error 1205: Lock wait timeout exceeded; try restarting transaction" {
			tx.Rollback()
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// Accounts 에 대해 각각 QueryAccountDetail 조회하기 전에 Account 기본 정보 세팅
func QueryAccountsDetailPrepare(accs *[]types.ResGetAccount) (err error) {
	err = config.DB.Table("account").Find(accs).Error
	return
}

// Account 상세 정보
func QueryAccountDetail(acc *types.ResGetAccount) (isNotExist bool, err error) {
	// ! Caution !
	// isNotExist = true => Not exist matched record
	// isExist = false => Exist matched record or occured error
	isNotExist = false
	defer func() {
		if err == gorm.ErrRecordNotFound {
			err = nil
			isNotExist = true
		}
	}()

	if err = config.DB.Table("account").Where("addr = ?", acc.Addr).First(acc).Error; err != nil {
		return
	}
	err = config.DB.Table("user_voucher_balance").Where("addr = ?", acc.Addr).Find(&acc.Vouchers).Error
	if err == nil {
		for i, vb := range acc.Vouchers {
			config.DB.Table("promotion").Where("promotion_id = ?", vb.PromotionId).First(&acc.Vouchers[i].Promotion)
		}
	}

	// summary
	// total_current_voucher_num, total_received_voucher_num
	sql := `
	SELECT 
		sum(current_amount) as total_current_voucher_num,
		sum(total_received_amount) as total_received_voucher_num
	FROM user_voucher_balance
	WHERE addr = ?
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	sql = `
	select count(*) as total_connect_num FROM event_wallet_conn
	WHERE addr = ?
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	sql = `
	select 
		count(*) as total_order_num
	FROM game_order
	WHERE addr = ?
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	sql = `
	select 
		count(*) as total_win_num
	FROM game_order
	WHERE addr = ? AND status > 2
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	sql = `
	select 
		count(*) as total_claimble_num
	FROM game_order GO
    LEFT JOIN promotion P ON GO.promotion_id=P.promotion_id
	WHERE 
		GO.addr = ? AND GO.status = 3 AND
        P.claim_start_at < NOW() AND P.claim_end_at > NOW()
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	// usd_of_win_order = (당첨 prize 의 amount * prize_denom_id 의 usd_value)
	// total_usd = usd_of_win_order 총합
	sql = `
	SELECT
		sum(P.amount*D.usd_price) as total_win_usd
	FROM game_order A
	LEFT JOIN prize P ON A.prize_id=P.prize_id
	LEFT JOIN prize_denom D ON P.prize_denom_id=D.prize_denom_id
	WHERE A.addr = ? AND A.is_win = true AND A.status in (3,4,5)
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

	sql = `
	SELECT
		IFNULL(sum(P.amount*D.usd_price), 0) as total_claimble_usd
	FROM game_order A
	LEFT JOIN prize P ON A.prize_id=P.prize_id
	LEFT JOIN prize_denom D ON P.prize_denom_id=D.prize_denom_id
	LEFT JOIN promotion PROM ON A.promotion_id=PROM.promotion_id
	WHERE 
	A.addr = ? AND 
			A.is_win = true AND 
			A.status = 3 AND
			PROM.claim_start_at < NOW() AND
			PROM.claim_end_at > NOW()
	`
	config.DB.Raw(sql, acc.Addr).Scan(&acc.Summary)

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
