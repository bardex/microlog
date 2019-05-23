package listeners

import (
	"fmt"
	"testing"
)

func TestStorageMemory(t *testing.T) {

	s := StorageMemory{}

	row := make(map[string]interface{})

	row = map[string]interface{}{"id": 1, "level": 6}
	s.Write(row)

	row = map[string]interface{}{"id": 2, "level": 6}
	s.Write(row)

	row = map[string]interface{}{"id": 3, "level": -5}
	s.Write(row)

	b := SearchBuilder{}
	f := b.Or(
		b.And(
			b.Equal("level", "6"),
			b.Equal("id", "2"),
		),
		b.Equal("id", "3"),
	)
	r := s.Search(f)

	fmt.Printf("%#v", r)
}
