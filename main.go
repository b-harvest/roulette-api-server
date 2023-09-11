package main

import (
	"fmt"
	config "roulette-api-server/config"
	schema "roulette-api-server/models/schema"
	routes "roulette-api-server/routes"

	"github.com/jinzhu/gorm"
)
var err error
func main() {
	config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer config.DB.Close()

	config.DB.AutoMigrate(
		&schema.OAuthAccessToken{},
		&schema.OAuthClient{},
		&schema.OAuthRefreshToken{},
		&schema.User{},
	)

	r := routes.SetupRouter()
	_ = r.Run()
}
