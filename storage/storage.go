package storage

import "sync"

var once sync.Once
var storage Storage

type Message map[string]interface{}

type Messages []Message

type Storage interface {
	Init() error
	Close() error
	Write(row Message) error
	Find(filters Query, page int32, limit int32) (Messages, error)
}

func GetStorage() (Storage, error) {
	once.Do(func() {
		storage = &StorageMemory{}
	})
	return storage, nil
}
