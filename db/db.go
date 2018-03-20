package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to find path to database: %v", err)
	}
	dbFile := ""
	for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
		p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/sayhuuzoku/db/data.db"))
		if _, err = os.Stat(p); err == nil {
			dbFile = p
			break
		}
	}
	if dbFile == "" {
		return fmt.Errorf("Failed to find path to database: %v", err)
	}

	if InitDB == true {
		os.Remove(dbFile)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}
	mydb.Connection = db

	q := "CREATE TABLE IF NOT EXISTS wakati_shopname ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", word VARCHAR(255) NOT NULL"
	q += ", position INTEGER NOT NULL"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"

	_, err = db.Exec(q)
	return err
}
