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

const COUNT_INSERTS = 1000

func BenchmarkEveryRowTransaction(b *testing.B) {
	os.Remove("test.db")
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
		for j := 0; j < COUNT_INSERTS; j++ {
			tx, _ := db.Begin()
			_, err := tx.Exec("INSERT INTO logs (message, date_creation) values ($1, CURRENT_TIMESTAMP)", msg)
			if err != nil {
				b.Fatal(err)
			}
			tx.Commit()
		}
	}

	db.Close()
	os.Remove("test.db")
}

func BenchmarkOneTransaction(b *testing.B) {
	os.Remove("test2.db")
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
		for j := 0; j < COUNT_INSERTS; j++ {
			_, err := tx.Exec("INSERT INTO logs (message, date_creation) values ($1, CURRENT_TIMESTAMP)", msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	tx.Commit()

	db.Close()
	os.Remove("test2.db")
}
