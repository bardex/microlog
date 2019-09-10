package storage

import (
	"testing"
)

func TestStorageMemoryPaginator(t *testing.T) {
	s := StorageMemory{}
	for i := 0; i < 10; i++ {
		s.Write(Row{"id": i + 1})
	}

	{
		// first page
		rows, err := s.Find(Query{}, 1, 2)

		if err != nil {
			t.Fatal(err)
		}

		if len(rows) != 2 {
			t.Fatalf("Expected len:2, actual:%d", len(rows))
		}

		if rows[0]["id"] != 1 {
			t.Fatalf("Expected id:1, actual:%d", rows[0]["id"])
		}

		if rows[1]["id"] != 2 {
			t.Fatalf("Expected id:2, actual:%d", rows[1]["id"])
		}
	}

	{
		// last page
		rows, err := s.Find(Query{}, 5, 2)

		if err != nil {
			t.Fatal(err)
		}

		if len(rows) != 2 {
			t.Fatalf("Expected len:2, actual:%d", len(rows))
		}

		if rows[0]["id"] != 9 {
			t.Fatalf("Expected id:9, actual:%d", rows[0]["id"])
		}

		if rows[1]["id"] != 10 {
			t.Fatalf("Expected id:10, actual:%d", rows[1]["id"])
		}
	}

	{
		// oversize page
		rows, err := s.Find(Query{}, 6, 2)

		if err != nil {
			t.Fatal(err)
		}

		if len(rows) != 0 {
			t.Fatalf("Expected len:0, actual:%d", len(rows))
		}
	}
}

func TestStorageMemoryFind(t *testing.T) {
	s := StorageMemory{}
	s.Write(Row{"id": 1, "status": 0})
	s.Write(Row{"id": 2, "status": 1})
	s.Write(Row{"id": 3, "status": 0})
	s.Write(Row{"id": 4, "status": 1})
	s.Write(Row{"id": 5, "status": 0}) //valid
	s.Write(Row{"id": 6, "status": 1})
	s.Write(Row{"id": 7, "status": 0}) //valid
	s.Write(Row{"id": 8, "status": 1})
	s.Write(Row{"id": 9, "status": 0}) //valid
	s.Write(Row{"id": 10, "status": 1})
	s.Write(Row{"id": 11, "status": 0, "time": "2019"})

	qb := QueryBuilder{}
	q := qb.And(qb.Equal("status", "0"), qb.GreaterOrEqual("id", "5"), qb.NotExists("time"))
	rows, err := s.Find(q, 1, 5)

	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 3 {
		t.Fatalf("Expected len:3, actual:%d", len(rows))
	}

	if rows[0]["id"] != 5 {
		t.Fatalf("Expected id:5, actual:%d", rows[0]["id"])
	}
	if rows[1]["id"] != 7 {
		t.Fatalf("Expected id:7, actual:%d", rows[1]["id"])
	}
	if rows[2]["id"] != 9 {
		t.Fatalf("Expected id:9, actual:%d", rows[2]["id"])
	}
}
