package db

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Record struct {
	HexId     string  `db:"hex_id" json:"data,omitempty"`
	TempValue float64 `db:"temp_value"`
	Seconds   float64 `db:"seconds"`
	DateTime  string  `db:"date_time"`
}

func Connection(driver, path string) *sqlx.DB {
	db, err := sqlx.Open(driver, path)

	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	return db
}

func GetAllSensors(DB *sqlx.DB) (sensors []string) {
	err := DB.Select(&sensors, `SELECT name FROM sensors_aliases`)

	if err != nil {
		log.Fatal("Cannot read from database: ", err)
	}
	return
}

func GetSensorData(DB *sqlx.DB, sensor string, limit int) (data []Record) {
	err := DB.Select(&data, `
		SELECT bs.temp_value, bs.seconds, bs.date_time
		FROM bathhouse_sensors bs
		LEFT JOIN sensors_aliases sa on sa.hex_id=bs.hex_id
		WHERE sa.name=$1
		ORDER BY bs.id DESC
		LIMIT $2
		`, sensor, limit)

	if err != nil {
		log.Fatal("Cannot read from database: ", err)
	}
	return
}
