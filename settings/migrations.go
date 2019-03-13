package settings

import "database/sql"

// APPEND ONLY END !
var migrations = []func(db *sql.DB) error{

	// create INPUTS table
	func(db *sql.DB) error {
		var err error
		_, err = db.Exec(`DROP TABLE IF EXISTS inputs`)
		if err != nil {
			return err
		}
		q := `CREATE TABLE inputs(
	 		 	id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	  			protocol TEXT,
	  			extractor TEXT,
	  			addr TEXT,
	  			enabled INTEGER,
	  			date_edit DATETIME)`
		_, err = db.Exec(q)
		if err != nil {
			return err
		}
		return nil
	},
}
