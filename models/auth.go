package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
)

func AuthAccessTokenCreate(data *schema.OAuthAccessToken) (err error) {
	if err = config.DB.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func AuthRefreshTokenCreate(data *schema.OAuthRefreshToken) (err error) {
	if err = config.DB.Create(data).Error; err != nil {
		return err
	}
	return nil
}
