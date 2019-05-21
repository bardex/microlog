package listeners

type Storage interface {
	// write message
	Write(map[string]interface{}) error

	// search messages
	//Search(SearchFilter) []map[string]interface{}
}
