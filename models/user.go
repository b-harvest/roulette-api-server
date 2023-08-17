package models

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/validations"

	_ "github.com/go-sql-driver/mysql"
)

func UserFetchAll(user *[]schema.User) (err error) {
	if err = config.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

func UserFetchSingle(user *schema.User, userId string) (err error) {
	if err = config.DB.Where("id = ?", userId).First(user).Error; err != nil {
		return err
	}
	return nil
}

func UserCreate(request *validations.UserCreate) (err error) {
	if err = config.DB.Table("users").Create(request).Error; err != nil {
		return err
	}
	return nil
}

func UserUpdate(request *validations.UserUpdate, userId string) (err error) {
	if err = config.DB.Table("users").Where("id = ?", userId).Update(request).Error; err != nil {
		return err
	}
	return nil
}

func UserDelete(user *schema.User, userId string) (err error) {
	if err = config.DB.Where("id = ?", userId).Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func UserFetchWithEmail(user *schema.User, email string) (err error) {
	if err = config.DB.Where("email = ?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}