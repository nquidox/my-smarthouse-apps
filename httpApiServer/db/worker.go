package db

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func Connection(driver, path string) *sqlx.DB {
	db, err := sqlx.Open(driver, path)

	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	return db
}
