package inputs

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var sqlitePath = "settings.db"

var db *sql.DB

type Record struct {
	Id       int
	Protocol string
	Addr     string
	Enabled  int8
}

func Connect() error {
	database, err := sql.Open("sqlite3", sqlitePath)
	db = database
	return err
}

func Disconnect() error {
	return db.Close()
}

func Install() error {

	sql := `DROP TABLE inputs`

	_, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}

	sql = `CREATE TABLE inputs(
	  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	  protocol TEXT,
	  addr TEXT,
	  enabled INTEGER,
	  date_edit DATETIME
	)`

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func Add(input Record) (int64, error) {
	result, err := db.Exec("INSERT INTO inputs (protocol, addr, enabled, date_edit) values ($1, $2, $3, CURRENT_TIMESTAMP)",
		input.Addr, input.Protocol, input.Enabled)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetAll() ([]Record, error) {

	inputs := []Record{}

	rows, err := db.Query("SELECT id, protocol, addr, enabled FROM inputs")
	if err != nil {
		return inputs, err
	}
	defer rows.Close()

	for rows.Next() {
		input := Record{}
		err := rows.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)
		if err != nil {
			fmt.Println(err)
			continue
		}
		inputs = append(inputs, input)
	}
	return inputs, nil
}

func GetOne(id int64) (Record, error) {

	row := db.QueryRow("SELECT id, protocol, addr, enabled FROM inputs WHERE id = $1", id)
	input := Record{}
	err := row.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)

	if err != nil {
		return input, err
	}
	return input, nil
}

func Update(input Record) (int64, error) {
	result, err := db.Exec("UPDATE inputs SET protocol = $2, addr = $1, enabled = $3, date_edit = CURRENT_TIMESTAMP WHERE id = $4", input.Addr, input.Protocol, input.Enabled, input.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func Delete(id int64) (int64, error) {

	result, err := db.Exec("DELETE FROM inputs WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
