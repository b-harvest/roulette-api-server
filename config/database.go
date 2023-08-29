package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	cfg, err := Load(DefaultConfigPath)
	if err != nil {
		panic(err)
	}
	dbConfig := DBConfig{
		Host:     cfg.DBConf.Host,
		Port:     cfg.DBConf.Port,
		User:    	cfg.DBConf.User,
		Password: cfg.DBConf.Password,
		DBName:   cfg.DBConf.DBName,
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}