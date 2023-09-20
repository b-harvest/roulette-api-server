package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func QueryPrizeInfosByPromotionId(prizeInfos *[]types.PrizeInfo, promotionId int64) (err error) {
	sql := `
		SELECT P.prize_id AS prize_id, P.dist_pool_id AS dist_pool_id, P.prize_denom_id AS prize_denom_id, P.amount AS amount, P.odds AS odds, P.win_cnt AS win_cnt, P.win_image_url AS win_image_url, P.max_daily_win_limit AS max_daily_win_limit, P.max_total_win_limit AS max_total_win_limit, P.is_active AS p_is_active,
			DP.total_supply AS total_supply, DP.remaining_qty AS remaining_qty, DP.is_active AS dp_is_active,
			PD.name AS name, PD.type AS type, PD.usd_price AS usd_price, PD.is_active AS pd_is_active
		FROM prize as P
			JOIN distribution_pool as DP
				ON P.dist_pool_id = DP.dist_pool_id
			JOIN prize_denom as PD
				ON P.prize_denom_id = PD.prize_denom_id
		WHERE P.promotion_id = ?
	`
	err = config.DB.Raw(sql, promotionId).Scan(prizeInfos).Error
	return
}

func UpdatePrizeByPrizeId(tx *gorm.DB, prize *schema.PrizeRow) (err error) {
	if tx == nil {
		tx = config.DB
	}

	err = tx.Table("prize").Where("prize_id = ?", prize.PrizeId).Update(prize).Error
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func QueryPrizes(prizes *[]schema.PrizeRow) (err error) {
	err = config.DB.Table("prize").Find(prizes).Error
	return
}

func CreatePrize(prize *schema.PrizeRow) (err error) {
	err = config.DB.Table("prize").Create(prize).Error
	return
}

// Prize 생성 with Tx
func CreatePrizeWithTx(tx *gorm.DB, prize *schema.PrizeRow) (err error) {
	err = tx.Table("prize").Create(prize).Error
	return
}

func QueryPrize(prize *schema.PrizeRow) (err error) {
	err = config.DB.Table("prize").Where("prize_id = ?", prize.PrizeId).First(prize).Error
	return
}

func UpdatePrize(prize *schema.PrizeRow) (err error) {
	err = config.DB.Table("prize").Where("prize_id = ?", prize.PrizeId).Update(prize).Error
	if err != nil {
		return err
	}
	err = config.DB.Table("prize").Where("prize_id = ?", prize.PrizeId).UpdateColumn("is_active", prize.IsActive).Error
	if err != nil {
		return err
	}
	return
}

func DeletePrize(prize *schema.PrizeRow) (err error) {
	err = config.DB.Table("prize").Where("prize_id = ?", prize.PrizeId).Delete(prize).Error
	return
}
