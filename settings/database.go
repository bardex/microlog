package settings

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var sqlitePath = "settings.db"

var db *sql.DB

func getDb() (*sql.DB, error) {
	var err error
	if db == nil {
		database, err := sql.Open("sqlite3", sqlitePath)
		if err == nil {
			db = database
		}
	}
	return db, err
}

func closeDb() error {
	err := db.Close()
	if err == nil {
		db = nil
	}
	return err
}
