package listeners

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestRepository(t *testing.T) {
	dbfile := "listener_repo_test.db"
	db, err := sql.Open("sqlite3", dbfile)

	if err != nil {
		t.Fatal(err)
	}

	//defer os.Remove(dbfile)

	repo, err := NewRepository(db)

	if err != nil {
		t.Fatal(err)
	}

	defer repo.Close()

	listener := NewListenerByParams(ProtocolTcp, ":8090", ExtractorJson)


	err = repo.Add(listener)

	if err != nil {
		t.Fatal(err)
	}

	if listener.Id == 0 {
		t.Fatal("listener.Id not set")
	}

	l, err := repo.Find(listener.Id)

	if err != nil {
		t.Fatal(err)
	}

	if l == nil {
		t.Fatal("Listener not found")
	}

	if l != listener {
		t.Fatal("Listeners not equals")
	}
}
