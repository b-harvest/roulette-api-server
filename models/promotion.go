package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryPromotions(promotions *[]schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Find(promotions).Error
	return
}

func CreatePromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Create(promotion).Error
	return
}

func QueryPromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).First(promotion).Error
	return
}

func UpdatePromotion(promotion *schema.PromotionRow) (err error) {
	// TBD: TotalSupply 가 변경된 경우 remainingQty 계산을 백앤드에서?
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).Update(promotion).Error
	return
}

func DeletePromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).Delete(promotion).Error
	return
}
