package ic

import (
	"encoding/json"
	"net/http"
)

const (
	companiesPath = "companies"
	customersPath = "customers"
)

type CompanyPair struct {
	Company Company
	Error   error
}

type CustomerPair struct {
	Customer Customer
	Error    error
}

const InvisibleCollectorUri = "https://api.invisiblecollector.com/"

type InvisibleCollector struct {
	apiRequest
}

func NewInvisibleCollector(apiKey string, apiUrl string) (*InvisibleCollector, error) {
	requests, err := newApiRequest(apiKey, apiUrl)
	if err != nil {
		return nil, err
	}

	return &InvisibleCollector{*requests}, nil
}

func (iC *InvisibleCollector) GetCompany(returnChannel chan<- CompanyPair) {

	iC.makeCompanyRequest(returnChannel, http.MethodGet, []string{companiesPath}, nil, nil)
}

func (iC *InvisibleCollector) SetCompany(returnChannel chan<- CompanyPair, companyUpdate Company) {

	iC.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath}, &companyUpdate, []fieldNamer{CompanyName, CompanyVatNumber})
}

func (iC *InvisibleCollector) SetCompanyNotifications(returnChannel chan<- CompanyPair, enableNotifications bool) {
	var notificationsPath string
	if enableNotifications {
		notificationsPath = "enableNotifications"
	} else {
		notificationsPath = "disableNotifications"
	}

	iC.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath, notificationsPath}, nil, nil)
}

func (iC *InvisibleCollector) SetNewCustomer(returnChannel chan<- CustomerPair, newCustomer Customer) {
	iC.makeCustomerRequest(returnChannel, http.MethodPost, []string{customersPath}, &newCustomer, []fieldNamer{CustomerName, CustomerVatNumber, CustomerCountry})
}

func (iC *InvisibleCollector) SetCustomer(returnChannel chan<- CustomerPair, updatedCustomer Customer) {
	iC.makeCustomerRequest(returnChannel, http.MethodPut, []string{customersPath, updatedCustomer.RoutableId()}, &updatedCustomer, []fieldNamer{CustomerCountry})
}

func (iC *InvisibleCollector) GetCustomer(returnChannel chan<- CustomerPair, customerId string) {
	iC.makeCustomerRequest(returnChannel, http.MethodGet, []string{customersPath, customerId}, nil, nil)
}

func (iC *InvisibleCollector) makeCustomerRequest(returnChannel chan<- CustomerPair, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) {

	customer := Customer{}
	err := iC.makeRequest(&customer, requestMethod, pathFragments, requestModel, mandatoryFields)
	returnChannel <- CustomerPair{customer, err}
}

func (iC *InvisibleCollector) makeCompanyRequest(returnChannel chan<- CompanyPair, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) {

	company := Company{}
	err := iC.makeRequest(&company, requestMethod, pathFragments, requestModel, mandatoryFields)
	returnChannel <- CompanyPair{company, err}
}

func (iC *InvisibleCollector) makeRequest(returnModel Modeler, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) error {

	var requestBody []byte = nil
	if requestModel != nil {
		if fieldErr := requestModel.AssertHasFields(mandatoryFields); fieldErr != nil {
			return fieldErr
		}

		requestJson, marshalErr := json.Marshal(requestModel)
		if marshalErr != nil {
			return marshalErr
		}

		requestBody = requestJson
	}

	returnJson, requestErr := iC.makeJsonRequest(requestBody, requestMethod, pathFragments)
	if requestErr != nil {
		return requestErr
	}

	return json.Unmarshal(returnJson, returnModel)
}
