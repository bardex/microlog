package settings

import "testing"

func TestAll(t *testing.T) {
	_, err1 := getDb()
	if err1 != nil {
		t.Fatalf("%s", err1)
	}

	err2 := closeDb()
	if err2 != nil {
		t.Fatalf("%s", err2)
	}
}
