package ic

import (
	"github.com/invisiblecloud/invisible-collector-go/internal"
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

func buildBarebonesTestServerRequest(statusCode int, useContentTypeHeader bool, returnJson string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if useContentTypeHeader {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(statusCode)
		io.WriteString(w, returnJson)
	})

	return httptest.NewServer(handler)
}

func buildAssertingTestServerRequest(t *testing.T, returnJson string, expectedMethod string, expectedUriPath string, expectedJsonBits []string) *httptest.Server {
	const jsonMime = "application/json"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUriPath, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)

		assert.Contains(t, r.Header.Get("Accept"), jsonMime)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
		assert.Contains(t, r.Header.Get("Authorization"), testApiKey)

		if expectedJsonBits != nil {
			assert.Contains(t, r.Header.Get("Content-Type"), jsonMime)
			assert.Contains(t, r.Header.Get("Content-Type"), "utf-8")

			bodyBytes, _ := internal.ReadCloseableBuffer(r.Body)
			bodyString := string(bodyBytes)
			for _, js := range expectedJsonBits {
				assert.Contains(t, bodyString, js)
			}

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

func assertCustomerRequest(t *testing.T, baseUri string, expectedModel model, method func(collector *InvisibleCollector, ch chan CustomerPair)) {

	collector := buildCollector(t, baseUri)
	ch := make(chan CustomerPair)
	go method(collector, ch)
	p := <-ch

	assertCorrectReturnedModel(t, expectedModel, p.Customer.model, p.Error)
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

	builder := buildTestCompanyModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()

	ts := buildAssertingTestServerRequest(t, jsonStr, "GET", "/companies", nil)
	defer ts.Close()

	assertCompanyRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CompanyPair) { collector.GetCompany(ch) })
}

func TestSetCompany(t *testing.T) {

	builder := buildTestCompanyModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()
	requestModel := Company{builder.buildRequestModel()}
	fragments := builder.getRequestJsonBits()

	ts := buildAssertingTestServerRequest(t, jsonStr, "PUT", "/companies", fragments)
	defer ts.Close()

	assertCompanyRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CompanyPair) { collector.SetCompany(ch, requestModel) })
}

func TestSetCompanyNotificationsEnable(t *testing.T) {

	builder := buildTestCompanyModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()

	ts := buildAssertingTestServerRequest(t, jsonStr, "PUT", "/companies/enableNotifications", nil)
	defer ts.Close()

	assertCompanyRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CompanyPair) { collector.SetCompanyNotifications(ch, true) })
}

func TestSetCompanyNotificationsDisable(t *testing.T) {

	builder := buildTestCompanyModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()

	ts := buildAssertingTestServerRequest(t, jsonStr, "PUT", "/companies/disableNotifications", nil)
	defer ts.Close()

	assertCompanyRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CompanyPair) { collector.SetCompanyNotifications(ch, false) })
}

func TestErrorNotJsonResponse(t *testing.T) {
	json := "{}"

	ts := buildBarebonesTestServerRequest(200, false, json)
	defer ts.Close()

	collector := buildCollector(t, ts.URL)
	ch := make(chan CompanyPair)
	go collector.GetCompany(ch)
	p := <-ch

	require.NotNil(t, p.Error)
}

func TestErrorBadJsonResponse(t *testing.T) {
	json := "/.;'\\"

	ts := buildBarebonesTestServerRequest(200, true, json)
	defer ts.Close()

	collector := buildCollector(t, ts.URL)
	ch := make(chan CompanyPair)
	go collector.GetCompany(ch)
	p := <-ch

	require.NotNil(t, p.Error)
}

func TestErrorNot200NotJson(t *testing.T) {
	msg := "an error occured"

	ts := buildBarebonesTestServerRequest(401, false, msg)
	defer ts.Close()

	collector := buildCollector(t, ts.URL)
	ch := make(chan CompanyPair)
	go collector.GetCompany(ch)
	p := <-ch

	require.NotNil(t, p.Error)
}

func TestErrorNot200JsonError(t *testing.T) {
	const json = `{
	"code": "401",
	"message": "testing"
}`

	ts := buildBarebonesTestServerRequest(401, true, json)
	defer ts.Close()

	collector := buildCollector(t, ts.URL)
	ch := make(chan CompanyPair)
	go collector.GetCompany(ch)
	p := <-ch

	require.NotNil(t, p.Error)
	assert.Contains(t, p.Error.Error(), "401")
	assert.Contains(t, p.Error.Error(), "testing")
}

func TestErrorNot200ConflictError(t *testing.T) {
	const json = `{
	"code": "409",
	"message": "testing",
	"gid": "conflict"
}`

	ts := buildBarebonesTestServerRequest(409, true, json)
	defer ts.Close()

	collector := buildCollector(t, ts.URL)
	ch := make(chan CompanyPair)
	go collector.GetCompany(ch)
	p := <-ch

	require.NotNil(t, p.Error)
	assert.Contains(t, p.Error.Error(), "conflict")

}

func TestGetCustomer(t *testing.T) {

	builder, id := buildTestCustomerModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()

	ts := buildAssertingTestServerRequest(t, jsonStr, "GET", "/customers/"+id, nil)
	defer ts.Close()

	assertCustomerRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CustomerPair) { collector.GetCustomer(ch, id) })
}

func TestSetNewCustomer(t *testing.T) {

	builder, _ := buildTestCustomerModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()
	requestModel := Customer{builder.buildRequestModel()}
	fragments := builder.getRequestJsonBits("gid")

	ts := buildAssertingTestServerRequest(t, jsonStr, "POST", "/customers", fragments)
	defer ts.Close()

	assertCustomerRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CustomerPair) { collector.SetNewCustomer(ch, requestModel) })
}

func TestSetCustomer(t *testing.T) {

	builder, id := buildTestCustomerModelBuilder()
	jsonStr := builder.buildJson()
	returnModel := builder.buildReturnModel()
	requestModel := Customer{builder.buildRequestModel()}
	fragments := builder.getRequestJsonBits("gid")

	ts := buildAssertingTestServerRequest(t, jsonStr, "PUT", "/customers/"+id, fragments)
	defer ts.Close()

	assertCustomerRequest(t, ts.URL, returnModel,
		func(collector *InvisibleCollector, ch chan CustomerPair) { collector.SetCustomer(ch, requestModel) })
}
