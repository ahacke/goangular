package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
	_ "github.com/mattn/go-sqlite3"
)

const (
	applicationName  string = "angolite"
	databaseFileName string = applicationName + ".db"
)

type entryStruct struct {
	Username   string
	Departname string
	Created    int64
}

func main() {
	initializeDatabase()

	http.Handle("/", http.FileServer(rice.MustFindBox("app").HTTPBox()))

	http.HandleFunc("/api/hello", helloWorld)
	http.HandleFunc("/api/entry", postEntry)
	http.HandleFunc("/api/entries", getEntries)
	http.ListenAndServe(":8080", nil)
}

func initializeDatabase() {
	db, err := sql.Open("sqlite3", "./"+databaseFileName)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS userinfo (uid INTEGER PRIMARY KEY AUTOINCREMENT,username VARCHAR(64) NULL,departname VARCHAR(64) NULL,created INTEGER NULL)")
	checkErr(err)

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	db, err := sql.Open("sqlite3", "./"+databaseFileName)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM userinfo")
	defer rows.Close()
	checkErr(err)

	data := []entryStruct{}
	var entry entryStruct
	var id int

	for rows.Next() {
		err = rows.Scan(&id, &entry.Username, &entry.Departname, &entry.Created)
		checkErr(err)
		data = append(data, entry)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}

func postEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
	}
	db, err := sql.Open("sqlite3", "./"+databaseFileName)
	checkErr(err)
	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var t entryStruct
	err = decoder.Decode(&t)
	checkErr(err)
	defer r.Body.Close()
	log.Printf("%+v\n", t)

	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(t.Username, t.Departname, t.Created)
	checkErr(err)

	w.WriteHeader(http.StatusOK)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Message string
	}{
		"Hello, World",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
