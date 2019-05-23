package listeners

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
)

// =========================================
//  ZLIB + JSON EXTRACTOR implements Extractor
// =========================================

const EXTRACTOR_ZLIB_JSON = "ZLIB_JSON"

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

func init() {
	// register extractor
	AddExtractor(ZlibJsonExtractor{})
}
