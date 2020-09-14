/*
Copyright 2019-2020 Special Brands Holding BV
Copyright 2014 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// outyet is a web server that announces whether or not a particular Go version
// has been tagged.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

// Command-line flags.
var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func main() {
	flag.Parse()

	db := connectToDatabase()

	var router = mux.NewRouter()

	router.HandleFunc("/get-data/{title}", passDBToHandler(db, getData)).Methods("GET")
	router.HandleFunc("/post-data/{title}", passDBToHandler(db, postData)).Methods("GET")

	log.Fatal(http.ListenAndServe(*httpAddr, router))

}

func passDBToHandler(db *sql.DB, f func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, db)
	}
}

// ServeHTTP implements the HTTP user interface.
func postData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	title := vars["title"]

	// write to database here
	// localhost:5000/post-data/(title)
	// {Title: “Name”}
	sqlStatement := `INSERT INTO app (title) VALUES ($1)`

	_, err := db.Query(sqlStatement, title)
	if err != nil {
		log.Print("Query:", err)
	}

	response := map[string]string{"Title": title}
	w.Header().Set("Content-Type", "application/json")

	// TODO: return data base on userAgent json or documentation
	json.NewEncoder(w).Encode(response)
}

func getData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	title := vars["title"]

	data := struct {
		Title     string
		UUID      string
		Timestamp string
	}{}

	// read from database here
	// localhost:5000/get-data/(title)
	// {Title: “Name”,UUID4: “12345678-1234-5678-1234-567812345678”, Timestamp: “2020-08-19T02:56:24+00:00”}
	sqlStatement := `SELECT * FROM app WHERE title = $1`
	rows, err := db.Query(sqlStatement, title)
	if err != nil {
		log.Print("Query:", err)
	}
	response := []map[string]string{}
	for rows.Next() {
		rows.Scan(&data.UUID, &data.Title, &data.Timestamp)
		log.Println("Value: ", data.Title, data.UUID, data.Timestamp)
		response = append(response, map[string]string{
			"Title":     data.Title,
			"UUID":      data.UUID,
			"Timestamp": data.Timestamp,
		})
	}

	// TODO: return data base on userAgent json or documentation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func connectToDatabase() (db *sql.DB) {
	url := os.Getenv("POSTGRES_URL")
	// FIXME: Why doesn't it take the argument out of the environment?
	db, err := sql.Open("postgres", url)
	//db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres.internal:5432/postgres?sslmode=disable")
	log.Print(url)
	if err != nil {
		log.Fatal("Connect:", err)
	}
	//	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Ping:", err)
	}
	return
}
