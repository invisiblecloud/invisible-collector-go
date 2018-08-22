package internal

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	jsonMime = "application/json"
)

type ApiRequest struct {
	apiKey string
	apiUrl url.URL
}

// it's thread-safe, recommended to reuse the client
var httpClient = http.Client{}

func NewApiRequests(apiKey string, apiUrl string) (*ApiRequest, error) {
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

	if isWhitespaceString(apiKey) {
		return nil, errors.New("invalid api key: " + apiKey)
	}

	return &ApiRequest{apiKey, *uri}, nil
}

func (api *ApiRequest) MakeJsonRequest(requestBody []byte, requestType string, pathSegments ...string) (returnBody []byte, err error) {
	request, requestErr := api.buildRequest(requestType, pathSegments, requestBody)
	if requestErr != nil {
		return nil, requestErr
	}

	response, clientErr := httpClient.Do(request)
	if clientErr != nil {
		return nil, clientErr
	}

	if response.StatusCode/100 != 2 {
		return nil, api.buildProtocolErrorMessage(response)
	}

	return readCloseableBuffer(response.Body)
}

func (api *ApiRequest) joinPathFragments(pathSegments []string) (string, error) {
	encodedPaths := make([]string, len(pathSegments))
	for i, pathSegment := range pathSegments {
		if isWhitespaceString(pathSegment) {
			return "", errors.New("Invalid uri path segment: " + pathSegment)
		}

		encodedPaths[i] = url.PathEscape(pathSegment)
	}

	return api.apiUrl.String() + "/" + strings.Join(encodedPaths, "/"), nil
}

func (api *ApiRequest) buildRequest(requestType string, pathSegments []string, requestBody []byte) (*http.Request, error) {
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
		request.Header.Set("Content-Type", jsonMime+"; charset=UTF-8")
	}

	return request, nil
}

func (api *ApiRequest) buildProtocolErrorMessage(response *http.Response) error {
	var httpError error = HttpStatusCodeError(response.StatusCode)

	// TODO check error json

	return httpError
}
