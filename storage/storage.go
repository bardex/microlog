package storage

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

type Storage interface {
	Write(row map[string]interface{}) error
	Find(filters SearchFilter, page int32, limit int32) ([]Message, error)
}
