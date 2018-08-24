package ic

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	testApiKey = "aded"
)

type requestPair struct {
	Error error
	Model model
}

func (p requestPair) buildFromCompanyPair(cp CompanyPair) requestPair {
	return requestPair{cp.Error, cp.Company.model}
}

func buildAssertingTestServerRequest(t *testing.T, returnJson string, expectedMethod string, expectedUriPath string, expectedJson interface{}) *httptest.Server {
	const jsonMime = "application/json"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUriPath, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)

		assert.Contains(t, r.Header.Get("Accept"), jsonMime)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
		assert.Contains(t, r.Header.Get("Authorization"), testApiKey)

		if expectedJson != nil {

			assert.Contains(t, r.Header.Get("Content-Type"), jsonMime)
			assert.Contains(t, r.Header.Get("Content-Type"), "utf-8")
		} else {
			assert.NotContains(t, r.Header.Get("Content-Type"), jsonMime)
		}

		if returnJson != "" {
			w.Header().Set("Content-Type", jsonMime)
			io.WriteString(w, returnJson)
		}
	})

	return httptest.NewServer(handler)
}

func buildCollector(t *testing.T, baseUri string) *InvisibleCollector {
	collector, err := NewInvisibleCollector(testApiKey, baseUri)
	require.Nil(t, err)

	return collector
}

func assertCorrectReturnedModel(t *testing.T, expected model, actual model, returnedErr error) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Models must be the same: \n expected: %v \n actual: %v", expected, actual)
	}

	require.Nil(t, returnedErr)
}

func assertCompanyRequest(t *testing.T, baseUri string, expectedModel model, method func(collector *InvisibleCollector, ch chan CompanyPair)) {

	collector := buildCollector(t, baseUri)
	ch := make(chan CompanyPair)
	go method(collector, ch)
	p := <-ch

	assertCorrectReturnedModel(t, expectedModel, p.Company.model, p.Error)
}

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

func TestGetCompany(t *testing.T) {

	builder := buildTestompanyModelBuilder()
	jsonStr := builder.buildJson()
	expectedModel := builder.buildModel()

	ts := buildAssertingTestServerRequest(t, jsonStr, "GET", "/companies", nil)
	defer ts.Close()

	assertCompanyRequest(t, ts.URL, expectedModel,
		func(collector *InvisibleCollector, ch chan CompanyPair) { collector.GetCompany(ch) })
}
