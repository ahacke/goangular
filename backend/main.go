package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type entryStruct struct {
	Username   string
	Departname string
	Created    int64
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./app/dist/angular-app")))

	http.HandleFunc("/api/hello", helloWorld)
	http.HandleFunc("/api/entry", postEntry)
	http.HandleFunc("/api/entries", getEntries)
	http.ListenAndServe(":8080", nil)

}

func getEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	db, err := sql.Open("sqlite3", "./foo.db")
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
	db, err := sql.Open("sqlite3", "./foo.db")
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

func sqliteTest() {

	// sqlite

	// table
	/*
			CREATE TABLE `userinfo` (
		        `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
		        `username` VARCHAR(64) NULL,
		        `departname` VARCHAR(64) NULL,
		        `created` DATE NULL
			);
	*/

	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	// update
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	rows.Close() //good habit to close

	// delete
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
