package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// MyDB : db struct
type MyDB struct {
	Connection *sql.DB
}

// InitDB : initialize database
var InitDB = false

// New : create db and keep connection
func (mydb *MyDB) New() error {
	absDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbFile := absDir + "/db/data.db"
	if InitDB == true {
		os.Remove(dbFile)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mydb.Connection = db

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
