package listeners

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"microlog/storage"
)

const ExtractorJson = "JSON"
const ExtractorString = "STRING"
const ExtractorZlibJson = "ZLIB_JSON"

type Extractor interface {
	Extract([]byte) (storage.Message, error)
}

type JsonExtractor struct{}

func (e JsonExtractor) Extract(buf []byte) (storage.Message, error) {
	msg := make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewReader(buf))
	err := dec.Decode(&msg)
	if err != nil && err != io.EOF {
		return msg, err
	}
	return msg, nil
}

type StringExtractor struct{}

func (e StringExtractor) Extract(buf []byte) (storage.Message, error) {
	msg := make(map[string]interface{})
	msg["msg"] = string(buf)
	return msg, nil
}

type ZlibJsonExtractor struct{}

func (e ZlibJsonExtractor) Extract(buf []byte) (storage.Message, error) {
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

// Get extractor by name
func GetExtractor(name string) (Extractor, error) {
	switch name {
	case ExtractorJson:
		return JsonExtractor{}, nil
	case ExtractorString:
		return StringExtractor{}, nil
	case ExtractorZlibJson:
		return ZlibJsonExtractor{}, nil
	default:
		return nil, fmt.Errorf("extractor %s not found", name)
	}
}
