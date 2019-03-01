package listeners

import (
	"bytes"
	"compress/zlib"
	"io"
	"fmt"
	"encoding/json"
)

const EXTRACTOR_ZLIB_JSON = "ZLIB_JSON"
const EXTRACTOR_JSON = "JSON"
const EXTRACTOR_STRING = "STRING"

type Message map[string]string

type Extractor interface {
	Extract([]byte) Message
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

type ZlibJsonExtractor struct {}

func (e ZlibJsonExtractor) Extract(buf []byte) Message {
	r, err := zlib.NewReader(bytes.NewReader(buf))
	defer r.Close()

	if err == nil {
		var data interface{}
		dec := json.NewDecoder(r)
		err = dec.Decode(&data)
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmt.Printf("%#v \n\n", data)
	}

	return Message{}
}

type JsonExtractor struct{}

func (e JsonExtractor) Extract(buf []byte) Message {
	var data interface{}
	dec := json.NewDecoder(bytes.NewReader(buf))
	err := dec.Decode(&data)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("%#v \n\n", data)

	return Message{}
}

type StringExtractor struct{}

func (e StringExtractor) Extract(buf []byte) Message {
	data := string(buf)

	fmt.Printf("%#v \n\n", data)

	return Message{}
}