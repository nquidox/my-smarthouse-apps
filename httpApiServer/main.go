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

func getValues(limit int) ([]byte, error) {
	type hourData struct {
		Name   string
		Params []db.Record
	}

	var output []hourData

	for _, sensor := range db.GetAllSensors(DB) {

		sensorData := db.GetSensorData(DB, sensor, limit)

		output = append(output, hourData{sensor, sensorData})
	}

	return json.Marshal(output)
}

func tempValuesHandler(limit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := getValues(limit)
		fmt.Fprintf(w, string(jsonResponse))
	}
}

func main() {
	CONFIG, _ = conf.LoadConfig(".")
	DB = db.Connection(CONFIG.Database.Driver, CONFIG.Database.DbPath)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /latest", tempValuesHandler(1))
	mux.HandleFunc("GET /hour", tempValuesHandler(60))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Have a nice day.")
	})

	addr := CONFIG.Server.Hostname + ":" + strconv.Itoa(int(CONFIG.Server.Port))
	log.Fatal(http.ListenAndServe(addr, mux))
}
