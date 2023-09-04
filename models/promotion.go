package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func QueryTbPromotions(promotions *[]schema.PromotionRow) (err error) {
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

//---------------------------------------------------------

func QueryPromotions(promotions *[](*types.ResGetPromotions)) (err error) {
	//q := "SELECT P.*, D. FROM promotion P "
	//err = config.DB.Table("promotion").Exec(q).Find(promotions).Error
	if err = config.DB.Table("promotion").Find(promotions).Error; err != nil {
		return
	}
	for _, v := range *promotions {
		fmt.Printf("%+v\n", v)
		if time.Now().After(v.PromotionEndAt) {
			v.Status = "finished"
		} else if time.Now().After(v.PromotionStartAt) {
			v.Status = "in progress"
		} else {
			v.Status = "not started"
		}
	}
	return
}

