package models

import (
	"roulette-api-server/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func CreateTxInstance() (tx *gorm.DB, err error) {
	tx = config.DB.Begin()
	err = tx.Error
	return
}
