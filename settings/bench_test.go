package settings

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"testing"
)

//
// RUN this:  go test -bench . -benchmem  bench_test.go
//

func BenchmarkEveryRowTransaction(b *testing.B) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		b.Fatal(err)
	}

	sql := `CREATE TABLE logs(
	  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	  message TEXT,
	  date_creation DATETIME
	)`

	_, err = db.Exec(sql)
	if err != nil {
		b.Fatal(err)
	}

	msg := strings.Repeat(" test ", 1000)

	for i := 0; i < b.N; i++ {
		tx, _ := db.Begin()
		_, err := db.Exec("INSERT INTO logs (message, date_creation) values ($1, CURRENT_TIMESTAMP)", msg)
		if err != nil {
			b.Fatal(err)
		}
		tx.Commit()
	}

	db.Close()
	os.Remove("test.db")
}

func BenchmarkOneTransaction(b *testing.B) {
	db, err := sql.Open("sqlite3", "test2.db")
	if err != nil {
		b.Fatal(err)
	}

	sql := `CREATE TABLE logs(
	  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	  message TEXT,
	  date_creation DATETIME
	)`

	_, err = db.Exec(sql)
	if err != nil {
		b.Fatal(err)
	}

	msg := strings.Repeat(" test ", 1000)

	tx, _ := db.Begin()

	for i := 0; i < b.N; i++ {
		_, err := db.Exec("INSERT INTO logs (message, date_creation) values ($1, CURRENT_TIMESTAMP)", msg)
		if err != nil {
			b.Fatal(err)
		}
	}

	tx.Commit()

	db.Close()
	os.Remove("test2.db")
}
