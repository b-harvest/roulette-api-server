package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"

	_ "github.com/go-sql-driver/mysql"
)

func QueryCurGameByAddr(game *schema.Game, addr string) (err error) {
	// 진행 중인 게임이 1개만 존재한다고 가정. First 한개만 가져 옴
	if err = config.DB.Table("game_order").Where("address = ? and status = 1", addr).First(game).Error; err != nil {
		return err
	}
	return nil
}

func StartNewGame(game *schema.Game, addr string) (err error) {
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

func StopGame(game *schema.Game, addr string) (err error) {
	game.Status = 2
	err = config.DB.Table("game_order").Where("address = ? and status = 1", addr).Update("status", game.Status).Error
	if err != nil {
		return err
	}
	return nil
}