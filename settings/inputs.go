package settings

const PROTOCOL_UDP = "udp"
const PROTOCOL_TCP = "tcp"

// Entity Input
type Input struct {
	Id       int64
	Protocol string
	Addr     string
	Enabled  int8
}

// Repository Inputs
type Inputs struct{}

// Repository Inputs: Add new input
func (inputs Inputs) Add(input *Input) error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}

	result, err := db.Exec("INSERT INTO inputs (protocol, addr, enabled, date_edit) values ($1, $2, $3, CURRENT_TIMESTAMP)",
		input.Protocol, input.Addr, input.Enabled)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	input.Id = id

	return nil
}

func (inputs Inputs) GetOne(id int64) (Input, error) {
	db, dbErr := getDb()
	if dbErr != nil {
		return Input{}, dbErr
	}

	row := db.QueryRow("SELECT id, protocol, addr, enabled FROM inputs WHERE id = $1", id)
	input := Input{}

	err := row.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)

	if err != nil {
		return input, err
	}
	return input, nil
}

func (inputs Inputs) GetAll() ([]Input, error) {

	items := []Input{}

	db, dbErr := getDb()
	if dbErr != nil {
		return items, dbErr
	}

	rows, err := db.Query("SELECT id, protocol, addr, enabled FROM inputs")
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		input := Input{}
		err := rows.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)
		if err != nil {
			return items, err
		}
		items = append(items, input)
	}
	return items, nil
}

func (inputs Inputs) Update(input *Input) error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}
	_, err := db.Exec("UPDATE inputs SET protocol = $1, addr = $2, enabled = $3, date_edit = CURRENT_TIMESTAMP WHERE id = $4", input.Protocol, input.Addr, input.Enabled, input.Id)
	return err
}

func (inputs Inputs) Delete(id int64) error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}
	_, err := db.Exec("DELETE FROM inputs WHERE id = $1", id)
	return err
}

func (inputs Inputs) Install() error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}

	sql := `DROP TABLE inputs`

	_, err := db.Exec(sql)

	if err != nil {
		return err
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
