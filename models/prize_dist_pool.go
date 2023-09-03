package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryDistPools(pools *[]schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Find(pools).Error
	return
}

func CreateDistPool(pool *schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Create(pool).Error
	return
}

func QueryDistPool(pool *schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Where("dist_pool_id = ?", pool.DistPoolId).First(pool).Error
	return
}

func UpdateDistPool(pool *schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Where("dist_pool_id = ?", pool.DistPoolId).Update(pool).Error
	return
}

func DeleteDistPool(pool *schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Where("dist_pool_id = ?", pool.DistPoolId).Delete(pool).Error
	return
}
