package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	. "grap-data/config"
	"log"
)

var DB *sqlx.DB

func init() {
	log.Println("config: DB")
	var err error
	driver := ViperConfig.Database.Driver
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true", ViperConfig.Database.User, ViperConfig.Database.Password,
		ViperConfig.Database.Host, ViperConfig.Database.Dbname)
	DB, err = sqlx.Open(driver, dsn)

	if err != nil {
		log.Fatal(err)
	}
	DB.SetMaxOpenConns(200)
	DB.SetMaxIdleConns(10)
	return
}
