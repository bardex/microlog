package listeners

import (
	"database/sql"
	"microlog/storage"
)

type Repository struct {
	db         *sql.DB
	inMemory   []*Listener
	logStorage storage.Storage
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
	id := int64(1)

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

func NewRepository(db *sql.DB, logStorage storage.Storage) (*Repository, error) {
	repo := &Repository{
		db:         db,
		logStorage: logStorage,
	}
	repo.Init()
	return repo, nil
}
