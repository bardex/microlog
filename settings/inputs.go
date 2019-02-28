package settings

import "microlog/listeners"

const PROTOCOL_UDP = "udp"
const PROTOCOL_TCP = "tcp"

// Entity Input
type Input struct {
	Id       int64
	Protocol string
	Addr     string
	Enabled  int8
	listener listeners.Listener
}

func (input *Input) GetListener() listeners.Listener {
	if input.listener == nil {
		switch input.Protocol {
		case PROTOCOL_UDP:
			input.listener = listeners.CreateUdp(input.Addr)
		}
	}
	return input.listener
}

// Repository inputRepository
type inputRepository struct {
	memory map[int64]*Input
}

// Repository instance
var Inputs inputRepository

func init() {
	Inputs.memory = make(map[int64]*Input)
}

// Repository inputRepository: Add new listeners
func (inputs inputRepository) Add(input *Input) error {
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

	inputs.memory[input.Id] = input

	return nil
}

func (inputs inputRepository) GetOne(id int64) (*Input, error) {
	// in memory first
	if val, ok := inputs.memory[id]; ok {
		return val, nil
	}

	db, dbErr := getDb()
	if dbErr != nil {
		return nil, dbErr
	}

	row := db.QueryRow("SELECT id, protocol, addr, enabled FROM inputs WHERE id = $1", id)
	input := &Input{}

	err := row.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)

	if err != nil {
		return nil, err
	}

	// save in memory
	inputs.memory[input.Id] = input

	return input, nil
}

func (inputs inputRepository) GetAll() ([]*Input, error) {

	// in memory first
	if len(inputs.memory) > 0 {
		items := make([]*Input, 0, len(inputs.memory))
		for _, item := range inputs.memory {
			items = append(items, item)
		}
		return items, nil
	}

	db, dbErr := getDb()
	if dbErr != nil {
		return nil, dbErr
	}

	rows, err := db.Query("SELECT id, protocol, addr, enabled FROM inputs")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]*Input, 0)

	for rows.Next() {
		input := &Input{}
		err := rows.Scan(&input.Id, &input.Protocol, &input.Addr, &input.Enabled)
		if err != nil {
			return nil, err
		}
		items = append(items, input)

		inputs.memory[input.Id] = input
	}
	return items, nil
}

func (inputs inputRepository) Update(input *Input) error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}
	_, err := db.Exec("UPDATE inputs SET protocol = $1, addr = $2, enabled = $3, date_edit = CURRENT_TIMESTAMP WHERE id = $4", input.Protocol, input.Addr, input.Enabled, input.Id)
	return err
}

func (inputs inputRepository) Delete(id int64) error {
	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}
	_, err := db.Exec("DELETE FROM inputs WHERE id = $1", id)

	delete(inputs.memory, id)

	return err
}

func (inputs inputRepository) Install() error {
	inputs.memory = make(map[int64]*Input)

	db, dbErr := getDb()
	if dbErr != nil {
		return dbErr
	}

	sql := `DROP TABLE inputs`

	db.Exec(sql)

	sql = `CREATE TABLE inputs(
	  id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	  protocol TEXT,
	  addr TEXT,
	  enabled INTEGER,
	  date_edit DATETIME
	)`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
