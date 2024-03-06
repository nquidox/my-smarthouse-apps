package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	conf "httpApiServer/config"
	"httpApiServer/db"
	"log"
	"net/http"
	"strconv"
)

var CONFIG conf.Config
var DB *sqlx.DB

func getValues() ([]byte, error) {

	type dbRecord struct {
		Name      string  `db:"name"`
		TempValue float64 `db:"temp_value"`
		Seconds   float64 `db:"seconds"`
		DateTime  string  `db:"date_time"`
	}

	var data []dbRecord

	err := DB.Select(&data, `
		SELECT sa.name, MAX(bs.temp_value) AS temp_value, MAX(bs.seconds) AS seconds,
		       MAX(bs.date_time) AS date_time
		FROM bathhouse_sensors bs
		INNER JOIN sensors_aliases sa ON bs.hex_id=sa.hex_id
		GROUP BY sa.name;
		`)

	if err != nil {
		log.Fatal("Cannot read from database: ", err)
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

func getHourValues() ([]byte, error) {
	type dbRecord struct {
		TempValue float64 `db:"temp_value"`
		Seconds   float64 `db:"seconds"`
		DateTime  string  `db:"date_time"`
	}

	type hourData struct {
		Name   string
		Params []dbRecord
	}

	var data []dbRecord
	var output []hourData
	var sensors []string

	err := DB.Select(&sensors, `SELECT name FROM sensors_aliases`)

	if err != nil {
		log.Fatal("Cannot read from database: ", err)
	}
	fmt.Println(sensors)

	for _, sensor := range sensors {
		err = DB.Select(&data, `
		SELECT bs.temp_value, bs.seconds, bs.date_time
		FROM bathhouse_sensors bs
		LEFT JOIN sensors_aliases sa on sa.hex_id=bs.hex_id
		WHERE sa.name=$1
		LIMIT 60
		`, sensor)

		if err != nil {
			log.Fatal("Cannot read from database: ", err)
		}
		output = append(output, hourData{sensor, data})
	}

	return json.Marshal(output)
}

func getHourData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		response, _ := getHourValues()
		fmt.Fprintf(w, string(response))
	}
}

func main() {
	CONFIG, _ = conf.LoadConfig(".")
	DB = db.Connection(CONFIG.Database.Driver, CONFIG.Database.DbPath)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Have a nice day.")
	})

	mux.HandleFunc("/bathhouse", getTempJSON)
	mux.HandleFunc("/hour", getHourData)

	addr := CONFIG.Server.Hostname + ":" + strconv.Itoa(int(CONFIG.Server.Port))
	log.Fatal(http.ListenAndServe(addr, mux))
}
