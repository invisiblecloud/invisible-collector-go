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

func TestInvalidApiKey(t *testing.T) {
	if _, err := NewInvisibleCollector("  \t\n", IcAddress); err == nil {
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	if _, err := NewInvisibleCollector(testApiKey, IcAddress); err != nil {
		t.Fail()
	}
}
