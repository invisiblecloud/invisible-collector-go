package invisibleCollector

import (
	"encoding/json"
	"github.com/invisiblecloud/invisible-collector-go/internal"
	"net/http"
)

const (
	companiesPath = "companies"
)

type modelPair struct {
	modeler
	error
}

type CompanyPair struct {
	*Company
	error
}

const IcAddress = "https://api.invisiblecollector.com/"

type InvisibleCollector struct {
	internal.ApiRequest
}

func NewInvisibleCollector(apiKey string, apiUrl string) (*InvisibleCollector, error) {
	requests, err := internal.NewApiRequests(apiKey, apiUrl)
	if err != nil {
		return nil, err
	}

	return &InvisibleCollector{*requests}, nil
}

func (ic *InvisibleCollector) GetCompany(returnChannel chan<- CompanyPair) {

	ic.makeCompanyRequest(returnChannel, http.MethodGet, []string{companiesPath}, nil, nil)
}

func (ic *InvisibleCollector) SetCompany(returnChannel chan<- CompanyPair, companyUpdate Company) {

	ic.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath}, &companyUpdate, []fieldNamer{CompanyName, CompanyVatNumber})
}

func (ic *InvisibleCollector) SetCompanyNotifications(returnChannel chan<- CompanyPair, enableNotifications bool) {
	var notificationsPath string
	if enableNotifications {
		notificationsPath = "enableNotifications"
	} else {
		notificationsPath = "disableNotifications"
	}

	ic.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath, notificationsPath}, nil, nil)
}

func (ic *InvisibleCollector) makeCompanyRequest(returnChannel chan<- CompanyPair, requestMethod string, pathFragments []string, requestModel modeler, mandatoryFields []fieldNamer) {

	company := Company{}
	err := ic.makeRequest(&company, requestMethod, pathFragments, requestModel, mandatoryFields)
	returnChannel <- CompanyPair{&company, err}
}

func (ic *InvisibleCollector) makeRequest(returnModel modeler, requestMethod string, pathFragments []string, requestModel modeler, mandatoryFields []fieldNamer) error {

	var requestBody []byte = nil
	if requestModel != nil {
		if fieldErr := requestModel.AssertHasFields(mandatoryFields); fieldErr != nil {
			return nil
		}

		requestJson, marshalErr := json.Marshal(requestModel)
		if marshalErr != nil {
			return nil
		}

		requestBody = requestJson
	}

	returnJson, requestErr := ic.MakeJsonRequest(requestBody, requestMethod, pathFragments)
	if requestErr != nil {
		return nil
	}

	return json.Unmarshal(returnJson, returnModel)
}
