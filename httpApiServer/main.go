package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

var config Config

func getValues() ([]byte, error) {
	db, err := sqlx.Open(config.Database.Driver, config.Database.DbPath)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	type dbRecord struct {
		Id        string  `db:"id"`
		HexId     string  `db:"hex_id"`
		TempValue float64 `db:"temp_value"`
		Seconds   float64 `db:"seconds"`
		DateTime  string  `db:"date_time"`
	}

	var data []dbRecord

	err2 := db.Select(&data, `SELECT * FROM bathhouse_sensors LIMIT 3`)

	if err2 != nil {
		log.Fatal("Cannot read from database: ", err2)
	}

	return json.Marshal(data)
}

func getTempJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		response, _ := getValues()
		fmt.Fprintf(w, string(response))
	}
}

func main() {
	config, _ = loadConfig(".")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Have a nice day.")
	})
	mux.HandleFunc("/bathhouse", getTempJSON)

	addr := config.Server.Hostname + ":" + strconv.Itoa(int(config.Server.Port))
	log.Fatal(http.ListenAndServe(addr, mux))
}
