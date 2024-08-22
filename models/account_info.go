package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	"github.com/jinzhu/gorm"
)

func QueryAccountInfoWithTx(tx *gorm.DB, acc_info *schema.AccountInfoRow) (bool, error) {
	if tx == nil {
		tx = config.DB
	}

	err := config.DB.Table("account_info").Where("addr = ?", acc_info.Addr).First(acc_info).Error
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

func CreateAccountInfoWithTx(tx *gorm.DB, acc_info *schema.AccountInfoRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := config.DB.Table("account_info").Create(acc_info).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func UpdateAccountInfoById(tx *gorm.DB, acc_info *schema.AccountInfoRow) error {
	if tx == nil {
		tx = config.DB
	}

	err := tx.Table("account_info").Where("id = ?", acc_info.ID).Update(acc_info).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}
