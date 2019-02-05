package ic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/invisiblecloud/invisible-collector-go/internal"
	"net/http"
	"net/url"
	"strings"
)

const (
	jsonMime = "application/json"
)

type apiRequest struct {
	apiKey string
	apiUrl url.URL
}

// it's thread-safe, recommended to reuse the client
var httpClient = http.Client{}

func newApiRequest(apiKey string, apiUrl string) (*apiRequest, error) {
	uri, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	if !uri.IsAbs() {
		return nil, errors.New("uri must be absolute")
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return nil, errors.New("invalid url scheme")
	}

	// TODO check if api key contains only ascii (for headers)

	if internal.IsWhitespaceString(apiKey) {
		return nil, errors.New("invalid api key: " + apiKey)
	}

	return &apiRequest{apiKey, *uri}, nil
}

// internal use
func (api *apiRequest) makeRequestExpectingJson(request *http.Request) (returnBody []byte, err error) {
	response, clientErr := httpClient.Do(request)
	if clientErr != nil {
		return nil, clientErr
	}

	family := response.StatusCode / 100
	if family != 2 && family != 3 && family != 1 {
		return nil, api.buildProtocolErrorMessage(response)
	}

	if !internal.JsonContentType(&response.Header) {
		return nil, errors.New("returned content-type isn't json")
	}

	return internal.ReadCloseableBuffer(response.Body)
}

func (api *apiRequest) makeJsonRequest(requestBody []byte, requestType string, pathSegments ...string) (returnBody []byte, err error) {
	request, requestErr := api.buildRequest(requestType, pathSegments, requestBody)
	if requestErr != nil {
		return nil, requestErr
	}

	return api.makeRequestExpectingJson(request)
}

func (api *apiRequest) makeUrlEncodedRequest(queryParams map[string]string, requestType string, pathSegments ...string) (returnBody []byte, err error) {
	request, requestErr := api.buildRequest(requestType, pathSegments, nil)
	if requestErr != nil {
		return nil, requestErr
	}

	if queryParams != nil && len(queryParams) != 0 {
		query := request.URL.Query()
		for k, v := range queryParams {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	return api.makeRequestExpectingJson(request)
}

func (api *apiRequest) joinPathFragments(pathSegments []string) (string, error) {
	encodedPaths := make([]string, len(pathSegments))
	for i, pathSegment := range pathSegments {
		if internal.IsWhitespaceString(pathSegment) {
			return "", errors.New("Invalid uri path segment: " + pathSegment)
		}

		encodedPaths[i] = url.PathEscape(pathSegment)
	}

	return api.apiUrl.String() + "/" + strings.Join(encodedPaths, "/"), nil
}

func (api *apiRequest) buildRequest(requestType string, pathSegments []string, requestBody []byte) (*http.Request, error) {
	if requestType != http.MethodGet && requestType != http.MethodPost && requestType != http.MethodPut && requestType != http.MethodDelete {
		panic("Internal error: invalid http request method.")
	}

	requestUri, pathErr := api.joinPathFragments(pathSegments)
	if pathErr != nil {
		return nil, pathErr
	}

	request, requestErr := http.NewRequest(requestType, requestUri, bytes.NewBuffer(requestBody))
	if requestErr != nil {
		return nil, requestErr
	}

	//request headers
	request.Header.Set("Host", api.apiUrl.Host)
	request.Header.Set("Accept", jsonMime)
	request.Header.Set("Authorization", "Bearer "+api.apiKey)
	if requestBody != nil && len(requestBody) != 0 {
		request.Header.Set("Content-Type", jsonMime+"; charset=utf-8")
	}

	return request, nil
}

func (api *apiRequest) buildProtocolErrorMessage(response *http.Response) error {
	var httpError error = HttpStatusCodeError(response.StatusCode)

	if internal.JsonContentType(&response.Header) {
		if m, err := internal.BufferToMap(response.Body); err == nil {
			statusCode := internal.MapGetValue(m, "code")
			msg := internal.MapGetValue(m, "message")
			id := internal.SliceFirstNonNil(internal.MapGetValue(m, "id"), internal.MapGetValue(m, "gid"))
			if msg != nil && statusCode != nil {
				errMsg := fmt.Sprintf("%v: %v", statusCode, msg)
				if id != nil && response.StatusCode == 409 {
					return ConflictError{errMsg, id.(string)}
				} else {
					return ServerError(errMsg)
				}
			}
		}
	}

	return httpError
}
