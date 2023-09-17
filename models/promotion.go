package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func QueryTbPromotions(promotions *[]schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Find(promotions).Error
	return
}

func CreateTbPromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Create(promotion).Error
	return
}

func QueryTbPromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).First(promotion).Error
	return
}

func UpdatePromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).Update(promotion).Error
	return
}

// 프로모션 업데이트 with tx
func UpdatePromotionWithTx(tx *gorm.DB, promotion *schema.PromotionRow) (err error) {
	err = tx.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).Update(promotion).Error
	return
}

func DeletePromotion(promotion *schema.PromotionRow) (err error) {
	err = config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).Delete(promotion).Error
	return
}

//---------------------------------------------------------

func QueryPromotions(promotions *[](*types.ResGetPromotions), qMap types.QueryFilterMap) (r [](*types.ResGetPromotions), err error) {
	fmt.Printf("%+v\n", qMap)
	q :=
		"SELECT P.*, IFNULL(CNT.participant_cnt, 0) as participant_cnt from promotion P " +
			"LEFT JOIN " +
			"  (select promotion_id, count(*) as participant_cnt from user_voucher_balance " +
			"   group by promotion_id) CNT ON P.promotion_id = CNT.promotion_id "
	// if err = config.DB.Raw(q).Scan(promotions).
	// 	Error; err != nil {
	// 	return
	// }
	cnt := 0
	for k, v := range qMap {
		if k != "status" {
			if cnt == 0 {
				q += " WHERE P."
			} else {
				q += " AND P."
			}
			fmt.Printf("key:%+v val:%+v\n", k, v)
			q += k + " = " + v + " "
			cnt++
		}
	}
	fmt.Println(q)
	sql := config.DB.Raw(q)
	if err = sql.Scan(promotions).Error; err != nil {
		return
	}

	// status 처리
	for _, v := range *promotions {
		if time.Now().After(v.PromotionEndAt) {
			v.Status = "finished"
		} else if time.Now().After(v.PromotionStartAt) {
			v.Status = "in-progress"
		} else {
			v.Status = "not-started"
		}
	}

	// status filter 처리
	if val, OK := qMap["status"]; OK {
		r = make([]*types.ResGetPromotions, 0, 100)
		for _, v := range *promotions {
			if v.Status == val {
				r = append(r, v)
			}
		}
	} else {
		r = *promotions
	}
	
	return
}

func QueryPromotion(promotion *types.ResGetPromotion) (err error) {
	q :=
		"SELECT P.*, IFNULL(CNT.participant_cnt, 0) as participant_cnt from promotion P " +
			"LEFT JOIN " +
			"  (select promotion_id, count(*) as participant_cnt from user_voucher_balance " +
			"   group by promotion_id) CNT ON P.promotion_id = CNT.promotion_id " +
			"WHERE P.promotion_id=?"
	// if err = config.DB.Table("promotion").Order("promotion_start_at DESC").Find(promotions).
	if err = config.DB.Raw(q, promotion.PromotionId).Scan(promotion).
		Error; err != nil {
		return
	}

	if time.Now().After(promotion.PromotionEndAt) {
		promotion.Status = "finished"
	} else if time.Now().After(promotion.PromotionStartAt) {
		promotion.Status = "in-progress"
	} else {
		promotion.Status = "not-started"
	}
	return
}

// 프로모션 요약 정보
func QueryPromotionSummary(promotionId uint64) (promSummary types.PromotionSummary, err error) {
	// TotalOdds
	q :=
		"  SELECT sum(P.odds) as total_odds FROM distribution_pool DP " +
			"  LEFT JOIN prize P ON DP.dist_pool_id=P.dist_pool_id        " +
			"  WHERE DP.promotion_id=? "
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	// TotalOrderNum
	q =
		"  SELECT count(*) as total_order_num FROM game_order " +
			"  WHERE promotion_id=? "
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	// TotalOrderUserNum
	q =
		"  SELECT account_id FROM game_order         " +
			"  WHERE promotion_id =? group by account_id "
	if err = config.DB.Raw(q, promotionId).Count(&promSummary.TotalOrderUserNum).Error; err != nil {
		if strings.Contains(err.Error(), "no rows") {
			promSummary.TotalOrderUserNum = 0
		} else {
			return
		}
	}

	// TotalWinNum
	q =
		"  SELECT count(*) as total_win_num FROM game_order " +
			"  WHERE promotion_id=? AND is_win=true "
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	// TotalWinUserNum
	q =
		"  SELECT account_id FROM game_order         " +
			"  WHERE promotion_id =? AND is_win=true group by account_id "
	if err = config.DB.Raw(q, promotionId).Count(&promSummary.TotalWinUserNum).Error; err != nil {
		if strings.Contains(err.Error(), "no rows") {
			promSummary.TotalOrderUserNum = 0
		} else {
			return
		}
	}

	// TotalWinNum
	q =
		"  SELECT count(*) as total_win_num FROM game_order " +
			"  WHERE promotion_id=? AND is_win=true "
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	// TotalWinUserNum
	q =
		"  SELECT account_id FROM game_order         " +
			"  WHERE promotion_id =? AND is_win=true group by account_id "
	if err = config.DB.Raw(q, promotionId).Count(&promSummary.TotalWinUserNum).Error; err != nil {
		if strings.Contains(err.Error(), "no rows") {
			promSummary.TotalOrderUserNum = 0
		} else {
			return
		}
	}

	// TotalClaimedNum
	q =
		"  SELECT count(*) as total_claimed_num FROM game_order " +
			"  WHERE promotion_id=? AND status >= 4"
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	// TotalClaimedUserNum
	q =
		"  SELECT account_id FROM game_order         " +
			"  WHERE promotion_id =? AND status >= 4 group by account_id "
	if err = config.DB.Raw(q, promotionId).Count(&promSummary.TotalClaimedUserNum).Error; err != nil {
		if strings.Contains(err.Error(), "no rows") {
			promSummary.TotalOrderUserNum = 0
		} else {
			return
		}
	}

	// InProgressClaimNum
	q =
		"  SELECT count(*) as in_progress_claim_num FROM game_order " +
			"  WHERE promotion_id=? AND status = 4"
	if err = config.DB.Raw(q, promotionId).Scan(&promSummary).Error; err != nil {
		return
	}

	err = nil
	return
}

// 프로모션 생성
func CreatePromotion(promotion *schema.PromotionRow) (err error) {
	if err = config.DB.Table("promotion").Create(promotion).Error; err != nil {
		return
	}
	return
}

// 프로모션 생성 with TX
func CreatePromotionWithTx(tx *gorm.DB, promotion *schema.PromotionRow) (err error) {
	if err = tx.Table("promotion").Create(promotion).Error; err != nil {
		return
	}
	return
}

func QueryPromotionById(promotion *schema.PromotionRow) error {
	return config.DB.Table("promotion").Where("promotion_id = ?", promotion.PromotionId).First(promotion).Error
}
