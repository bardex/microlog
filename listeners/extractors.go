package listeners

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
)

const EXTRACTOR_ZLIB_JSON = "ZLIB_JSON"
const EXTRACTOR_JSON = "JSON"
const EXTRACTOR_STRING = "STRING"

var ExtractorsList = map[string]string{
	EXTRACTOR_ZLIB_JSON: "zlib + JSON (GrayLog)",
	EXTRACTOR_JSON:      "JSON",
	EXTRACTOR_STRING:    "String",
}

type Extractor interface {
	Extract([]byte) (map[string]interface{}, error)
}

func createExtractor(name string) Extractor {
	var e Extractor
	switch name {
	case EXTRACTOR_ZLIB_JSON:
		e = ZlibJsonExtractor{}
	case EXTRACTOR_JSON:
		e = JsonExtractor{}
	case EXTRACTOR_STRING:
		e = StringExtractor{}
	}
	return e
}

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

type StringExtractor struct{}

func (e StringExtractor) Extract(buf []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	msg["msg"] = string(buf)
	return msg, nil
}
