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

// List of available Extractors
var Extractors = make(map[string]Extractor, 0)

// Get extractor by name
func GetExtractor(name string) (Extractor, error) {
	ext, ok := Extractors[name]
	if ok {
		return ext, nil
	} else {
		return nil, fmt.Errorf("extractor %s not found", name)
	}
}

// Register extractor
func AddExtractor(extractor Extractor) {
	name := extractor.GetName()
	Extractors[name] = extractor
}
