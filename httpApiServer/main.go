package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func getValues() ([]byte, error) {
	dbFile := "database.db"
	db, err := sqlx.Open("sqlite3", dbFile)
	_ = db
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
	fmt.Println(data)

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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Have a nice day.")
	})
	mux.HandleFunc("/bathhouse", getTempJSON)

	log.Fatal(http.ListenAndServe(":9009", mux))

}
