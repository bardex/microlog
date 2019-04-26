package listeners

import (
	"encoding/json"
	"bytes"
	"io"
)

// =========================================
//  JSON EXTRACTOR implements Extractor
// =========================================

const EXTRACTOR_JSON = "JSON"

type JsonExtractor struct{}

func (e JsonExtractor) Extract(buf []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewReader(buf))
	err := dec.Decode(&msg)
	if err != nil && err != io.EOF {
		return msg, err
	}
	return msg, nil
}

func (e JsonExtractor) GetName() string {
	return EXTRACTOR_JSON
}

func (e JsonExtractor) GetDescription() string {
	return "JSON"
}

func init() {
	// register extractor
	AddExtractor(JsonExtractor{})
}
