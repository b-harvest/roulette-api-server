package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryPrizeDenoms(denoms *[]schema.PrizeDenomRow) (err error) {
	err = config.DB.Table("prize_denom").Find(denoms).Error
	return
}

func CreatePrizeDenom(denom *schema.PrizeDenomRow) (err error) {
	err = config.DB.Table("prize_denom").Create(denom).Error
	return
}

func QueryPrizeDenom(denom *schema.PrizeDenomRow) (err error) {
	err = config.DB.Table("prize_denom").Where("prize_denom_id = ?", denom.PrizeDenomId).First(denom).Error
	return
}

func UpdatePrizeDenom(denom *schema.PrizeDenomRow) (err error) {
	err = config.DB.Table("prize_denom").Where("prize_denom_id = ?", denom.PrizeDenomId).Update(denom).Error
	if err != nil {
		return err
	}
	err = config.DB.Table("prize_denom").Where("prize_denom_id = ?", denom.PrizeDenomId).UpdateColumn("is_active", denom.IsActive).Error
	if err != nil {
		return err
	}
	return
}

func DeletePrizeDenom(denom *schema.PrizeDenomRow) (err error) {
	err = config.DB.Table("prize_denom").Where("prize_denom_id = ?", denom.PrizeDenomId).Delete(denom).Error
	return
}

func UpdatePrizeDenomByName(denom *schema.PrizeDenomRow) (err error) {
	return config.DB.Table("prize_denom").Where("name = ?", denom.Name).Update(denom).Error
}