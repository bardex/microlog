package listeners

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"microlog/storage"
	"testing"
)

func TestRepository(t *testing.T) {
	dbfile := "listener_repo_test.db"
	db, err := sql.Open("sqlite3", dbfile)

	if err != nil {
		t.Fatal(err)
	}

	//defer os.Remove(dbfile)

	storage := &storage.StorageMemory{}
	repo, err := NewRepository(db, storage)

	if err != nil {
		t.Fatal(err)
	}

	defer repo.Close()

	listener := NewListenerByParams(ProtocolTcp, ":8090", ExtractorJson)

	repo.Add(listener)

	if listener.Id == 0 {
		t.Fatal("listener.Id not set")
	}

}
