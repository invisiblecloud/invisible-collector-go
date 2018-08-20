package Ic

import "net/url"

type invisibleCollector struct {
	apiKey string
	apiUrl url.URL
}

const IcAddress = "https://api.invisiblecollector.com/"

func NewInvisibleCollector(apiKey string, apiUrl string) (*invisibleCollector, error) {
	uri, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	return &invisibleCollector{apiKey, *uri}, nil
}

func (ic invisibleCollector) Hello() string {
	return "hi"
}
