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

type Row map[string]interface{}

type Rows []Row

type Storage interface {
	Write(row Row) error
	Find(filters SearchFilter, page int32, limit int32) (Rows, error)
}
