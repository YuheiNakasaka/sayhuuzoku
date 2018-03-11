package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// TODO:
// db構造体つくる
// Newでdbへの接続つくる
// dbにinsertとcreateのメソッドはやす

// MyDB : db struct
type MyDB struct {
	connection *sql.DB
}

// New : create db and keep connection
func (mydb *MyDB) New() error {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mydb.connection = db
	defer db.Close()

	q := "CREATE TABLE IF NOT EXISTS wakati_shopname ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", word VARCHAR(255) NOT NULL"
	q += ", position INTEGER NOT NULL"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"

	_, err = db.Exec(q)
	if err != nil {
		panic(err)
	}

	return err
}
