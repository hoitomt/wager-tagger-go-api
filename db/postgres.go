package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Run `heroku pg:credentials DATABASE` from the ruby sportsbook_api for the connection string
func GetDb() gorm.DB {
	db, err := gorm.Open("postgres",
		"dbname=sportsbook_api host=localhost port=5432 user=hoitomt password=badger sslmode=disable")
	// db, err := gorm.Open("postgres",
	//  "dbname=d8ihclom31aprt host=ec2-54-197-246-197.compute-1.amazonaws.com port=5432 user=mffrovlynusepg password=uL6PkJjdATMqsSnvhwFzCIlh3X sslmode=require")
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	return db
}
