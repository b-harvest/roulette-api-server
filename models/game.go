package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func QueryCurGameByAddr(game *schema.GameOrder, addr string) (err error) {
	// 진행 중인 게임이 1개만 존재한다고 가정. First 한개만 가져 옴
	if err = config.DB.Table("game_order").Where("address = ? and status = 1", addr).First(game).Error; err != nil {
		return err
	}
	return nil
}

func StartNewGame(game *schema.GameOrder, addr string) (err error) {
	game.Address = addr
	game.Status = 1
	if err = config.DB.Table("game_order").Create(game).Error; err != nil {
		return err
	}

	// 생성된 게임 조회
	if err = config.DB.Table("game_order").Where("address = ? and status = 1", addr).First(game).Error; err != nil {
		return err
	}
	return nil
}

func StopGame(game *schema.GameOrder, addr string) (err error) {
	game.Status = 2
	err = config.DB.Table("game_order").Where("address = ? and status = 1", addr).Update("status", game.Status).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryGameTypes(games *[]schema.Game) (err error) {
	err = config.DB.Table("game_type").Find(games).Error
	return
}

func CreateGame(game *schema.Game) (err error) {
	// err = config.DB.Table("game_type").FirstOrCreate(game).Error
	err = config.DB.Table("game_type").Create(game).Error
	return
}

func QueryGameType(tx *gorm.DB, game *schema.Game) (err error) {
	if tx == nil {
		tx = config.DB
	}

	err = tx.Table("game_type").Where("game_id = ?", game.GameId).First(game).Error
	if err != nil {
		tx.Rollback()
	}
	return
}

func UpdateGame(game *schema.Game) (err error) {
	err = config.DB.Table("game_type").Where("game_id = ?", game.GameId).Update(game).Error
	if err != nil {
		return err
	}
	err = config.DB.Table("game_type").Where("game_id = ?", game.GameId).UpdateColumn("is_active", game.IsActive).Error
	if err != nil {
		return err
	}
	return
}

func DeleteGame(game *schema.Game) (err error) {
	//err = config.DB.Table("game_type").Where("game_id = ?", game.GameId).Update(game).Error
	err = config.DB.Table("game_type").Where("game_id = ?", game.GameId).Delete(game).Error
	return
}
