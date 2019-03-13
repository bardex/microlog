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

func Migrate() {
	db, e := getDb()
	if e != nil {
		panic(e)
	}
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id INTEGER NOT NULL UNIQUE)`)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT id FROM migrations ORDER BY id")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	exists := make(map[int]bool)

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		exists[id] = true
	}

	for i, f := range migrations {
		if _, exist := exists[i]; !exist {
			err := f(db)
			if err != nil {
				panic(err)
			}
			_, e := db.Exec(`INSERT INTO migrations (id) VALUES ($1)`, i)
			if e != nil {
				panic(e)
			}
		}
	}

}
