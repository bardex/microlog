package listeners

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db         *sql.DB
	inMemory   []*Listener
}

func (r *Repository) Init() error {
	q := `CREATE TABLE IF NOT EXISTS listeners (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
			protocol TEXT,
			extractor TEXT,
			addr TEXT,
			active INTEGER,
			date_edit DATETIME)`

	_, err := r.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() error {
	return nil
}

func (r *Repository) Add(listener *Listener) error {
	args := sqlite3.NamedArgs{"$a": 1, "$b": "demo"}
	c.Exec("INSERT INTO x VALUES($a, $b, $c)", args) // $c will be NULL


	listener.Id = id

	r.inMemory = append(r.inMemory, listener)

	return nil
}

func (r *Repository) Update(listener *Listener) error {
	return nil
}

func (r *Repository) Delete(listener *Listener) error {
	return nil
}

func (r *Repository) FindAll() ([]*Listener, error) {
	return r.inMemory, nil
}

func (r *Repository) Find(id int64) (*Listener, error) {
	all, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, l := range all {
		if l.Id == id {
			return l, nil
		}
	}

	return nil, nil
}

func NewRepository(db *sql.DB) (*Repository, error) {
	repo := &Repository{
		db:         db,
	}
	repo.Init()
	return repo, nil
}
