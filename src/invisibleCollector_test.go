package ic

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testApiKey = "aded"
)

func TestInvalidUri(t *testing.T) {
	_, err := NewInvisibleCollector(testApiKey, "ftp://123.23.23.23")
	require.NotNil(t, err)
}

func TestInvalidApiKey(t *testing.T) {
	_, err := NewInvisibleCollector("  \t\n", InvisibleCollectorUri)
	require.NotNil(t, err)
}

func TestNew(t *testing.T) {
	_, err := NewInvisibleCollector(testApiKey, InvisibleCollectorUri)
	require.Nil(t, err)
}
