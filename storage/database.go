package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var sqlitePath = "storage.db"

var db *sql.DB

func openDb() error {
	var err error
	if db == nil {
		database, err := sql.Open("sqlite3", sqlitePath)
		if err == nil {
			_, err = database.Exec("PRAGMA synchronous=NORMAL")
			if err != nil {
				return err
			}
			_, err = database.Exec("PRAGMA journal_mode=WAL")
			if err != nil {
				return err
			}
			db = database
		}
	}
	return err
}

func closeDb() error {
	err := db.Close()
	if err == nil {
		db = nil
	}
	return err
}
