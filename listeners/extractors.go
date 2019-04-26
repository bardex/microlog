package listeners

import (
	"fmt"
)

//  Extractor interface.

type Extractor interface {
	Extract([]byte) (map[string]interface{}, error)
	GetName() string
	GetDescription() string
}

// List of available extractors
var extractors = make(map[string]Extractor, 0)

// Get extractor by name
func GetExtractor(name string) (Extractor, error) {
	ext, ok := extractors[name]
	if ok {
		return ext, nil
	} else {
		return NilExtractor{}, fmt.Errorf("extractor %s not found", name)
	}
}

// Register extractor
func AddExtractor(extractor Extractor) {
	name := extractor.GetName()
	extractors[name] = extractor
}

type NilExtractor struct{}

func (e NilExtractor) Extract(buf []byte) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func (e NilExtractor) GetName() string {
	return "NIL"
}

func (e NilExtractor) GetDescription() string {
	return "Nil"
}
