package ic

import "testing"

const (
	testApiKey = "aded"
)

func TestInvalidUri(t *testing.T) {
	if _, err := NewInvisibleCollector(testApiKey, "ftp://123.23.23.23"); err == nil {
		t.Fail()
	}

}
