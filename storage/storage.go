package storage

import "sync"

var once sync.Once
var storage Storage

// Entity Field
type Field struct {
	Key   string
	Value string
}

// Entity Message
type Message struct {
	MessageId int64
	Time      string
	Fields    []Field
}

type Row map[string]interface{}

type Rows []Row

type Storage interface {
	Init() error
	Close() error
	Write(row Row) error
	Find(filters Query, page int32, limit int32) (Rows, error)
}

func CreateStorage() (Storage, error) {
	once.Do(func() {
		storage = &StorageMemory{}
	})
	return storage, nil
}
