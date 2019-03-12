package listeners

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"testing"
)

func TestStringExtract(t *testing.T) {
	extractor := StringExtractor{}
	input := "ABCDEFGH"
	buffer := []byte(input)
	output, err := extractor.Extract(buffer)

	if err != nil {
		t.Fatal(err)
	}

	if output["msg"] != input {
		t.Fatalf("Output (%s) not equal input (%s).", output, input)
	}
}

func TestJsonExtract(t *testing.T) {
	in := make(map[string]interface{})
	in["facility"] = "test"
	in["level"] = 1
	in["price"] = 10.5
	in["date"] = "2018-10-20"

	data, _ := json.Marshal(in)
	extractor := JsonExtractor{}
	out, err := extractor.Extract(data)

	if err != nil {
		t.Fatal(err)
	}

	for k, inv := range in {
		if outv, ok := out[k]; ok {
			if fmt.Sprintf("%v", outv) != fmt.Sprintf("%v", inv) {
				t.Fatalf("Output (%v) not equal input (%v) for key %s", outv, inv, k)
			}
		} else {
			t.Fatalf("Output key (%s) not found.", k)
		}
	}
}

func TestZlibJsonExtract(t *testing.T) {
	in := make(map[string]interface{})
	in["facility"] = "test"
	in["level"] = 1
	in["price"] = 10.5
	in["date"] = "2018-10-20"

	data, _ := json.Marshal(in)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()

	data = []byte(b.String())

	extractor := ZlibJsonExtractor{}
	out, err := extractor.Extract(data)

	if err != nil {
		t.Fatal(err)
	}

	for k, inv := range in {
		if outv, ok := out[k]; ok {
			if fmt.Sprintf("%v", outv) != fmt.Sprintf("%v", inv) {
				t.Fatalf("Output (%v) not equal input (%v) for key %s", outv, inv, k)
			}
		} else {
			t.Fatalf("Output key (%s) not found.", k)
		}
	}
}

func TestExtractorFactory(t *testing.T) {
	jsonExt, err := createExtractor("JSON")
	if err != nil {
		t.Fatal(err)
	}

	if jsonExt.GetName() != "JSON" {
		t.Fatal("Not JSON extractor")
	}

	zlibJsonExt, err := createExtractor("ZLIB_JSON")
	if err != nil {
		t.Fatal(err)
	}

	if zlibJsonExt.GetName() != "ZLIB_JSON" {
		t.Fatal("Not ZLIB_JSON extractor")
	}

	stringExt, err := createExtractor("STRING")
	if err != nil {
		t.Fatal(err)
	}

	if stringExt.GetName() != "STRING" {
		t.Fatal("Not STRING extractor")
	}

	_, notNilError := createExtractor("any")
	if notNilError == nil {
		t.Fatal("Expected error")
	}
}
