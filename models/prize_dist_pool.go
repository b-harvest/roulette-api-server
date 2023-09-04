package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
)

func QueryTbDistPools(pools *[]schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Find(pools).Error
	return
}

func CreateDistPool(pool *schema.PrizeDistPoolRow) (err error) {
	err = config.DB.Table("distribution_pool").Create(pool).Error
	return
}

func QueryTbDistPool(pool *schema.PrizeDistPoolRow) (err error) {
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

// dist_pool 리스트 및 해당 pool 에 속하는 prize 리스트
func QueryDistPoolsByPromId(id uint64) (*[]types.PrizeDistPool, error) {
	pools := make([]types.PrizeDistPool, 0, 100)

	// dist_pool 리스트
	q := 
		"SELECT D.*, PD.name AS name, PD.type AS type FROM distribution_pool D " + 
		"LEFT JOIN prize_denom PD ON D.prize_denom_id=PD.prize_denom_id " +
		"WHERE D.promotion_id = ?"
	// err := config.DB.Table("distribution_pool").Exec(q).Where("promotion_id = ?", id).Find(&pools).Error
	err := config.DB.Raw(q, id).Scan(&pools).Error
	if err != nil {
		return nil, err
	} 

	// dist_pool 별 prize 리스트
	for i, v := range pools {
		// qPrizes :=
		tmpPrizes := make([]types.Prize, 0, 100)
		if err = config.DB.Table("prize").Where("dist_pool_id=?", v.DistPoolId).
			Find(&tmpPrizes).Error; err != nil {
				return nil, err
		}
		pools[i].Prizes = &tmpPrizes
	}

	return &pools, err
}

