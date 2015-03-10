package db

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Run `heroku pg:credentials DATABASE` from the ruby sportsbook_api for the connection string
func GetDb() gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("PG_DB_URL"))

	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	return db
}
