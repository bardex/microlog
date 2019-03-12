package listeners

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
)

const EXTRACTOR_ZLIB_JSON = "ZLIB_JSON"
const EXTRACTOR_JSON = "JSON"
const EXTRACTOR_STRING = "STRING"

var Extractors = []Extractor{
	ZlibJsonExtractor{},
	JsonExtractor{},
	StringExtractor{},
}

type Extractor interface {
	Extract([]byte) (map[string]interface{}, error)
	GetName() string
	GetDescription() string
}

func createExtractor(name string) (Extractor, error) {
	for _, v := range Extractors {
		if v.GetName() == name {
			ext := v
			return ext, nil
		}
	}
	return StringExtractor{}, fmt.Errorf("extractor %s not found", name)
}

// ------  ZLIB_JSON ----------------------------------------
type ZlibJsonExtractor struct{}

func (e ZlibJsonExtractor) Extract(buf []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	r, err := zlib.NewReader(bytes.NewReader(buf))

	if err != nil {
		return msg, err
	}

	defer r.Close()

	dec := json.NewDecoder(r)
	err = dec.Decode(&msg)
	if err != nil && err != io.EOF {
		return msg, err
	}
	return msg, nil
}

func (e ZlibJsonExtractor) GetName() string {
	return EXTRACTOR_ZLIB_JSON
}

func (e ZlibJsonExtractor) GetDescription() string {
	return "zlib + JSON (GrayLog)"
}

// ------ JSON -----------------------------------------------
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

// ------ STRING ---------------------------------------------
type StringExtractor struct{}

func (e StringExtractor) Extract(buf []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	msg["msg"] = string(buf)
	return msg, nil
}

func (e StringExtractor) GetName() string {
	return EXTRACTOR_STRING
}

func (e StringExtractor) GetDescription() string {
	return "String"
}
