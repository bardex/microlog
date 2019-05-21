package listeners

// =========================================
//  STRING EXTRACTOR implements Extractor
// =========================================

const EXTRACTOR_STRING = "STRING"

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

func init() {
	// register extractor
	AddExtractor(StringExtractor{})
}
