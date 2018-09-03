// The go library used to communicate with Invisible Collector services
package ic

import (
	"encoding/json"
	"net/http"
)

const (
	companiesPath          = "companies"
	customersPath          = "customers"
	customerAttributesPath = "attributes"
	customerDebtsPath      = "debts"
	debtsPath              = "debts"
	InvisibleCollectorUri  = "https://api.invisiblecollector.com/"
)

type CompanyPair struct {
	Company Company
	Error   error
}

type CustomerPair struct {
	Customer Customer
	Error    error
}

type AttributesPair struct {
	Attributes map[string]string
	Error      error
}

type DebtPair struct {
	Debt  Debt
	Error error
}

type DebtListPair struct {
	Debts []Debt
	Error error
}

// The entry point for using this library.
// Contains all the implemented API request by this library that are used to communicate with Invisible Collector
type InvisibleCollector struct {
	apiRequest
}

// Constructor for this class.
//
// Use your company API key. You can use the InvisibleCollectorUri constant to use the default Invisible Collector URI.
func NewInvisibleCollector(apiKey string, apiUrl string) (*InvisibleCollector, error) {
	requests, err := newApiRequest(apiKey, apiUrl)
	if err != nil {
		return nil, err
	}

	return &InvisibleCollector{*requests}, nil
}

// Get the company information.
func (iC *InvisibleCollector) GetCompany(returnChannel chan<- CompanyPair) {

	iC.makeCompanyRequest(returnChannel, http.MethodGet, []string{companiesPath}, nil, nil)
}

// Update company information.
//
// Only "address", "zipCode" and "city" can be changed. "name" and "vatNumber" are mandatory fields that are used for validation.
func (iC *InvisibleCollector) SetCompany(returnChannel chan<- CompanyPair, companyUpdate Company) {

	iC.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath}, &companyUpdate, []fieldNamer{CompanyName, CompanyVatNumber})
}

// a
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

func (iC *InvisibleCollector) SetCustomerAttributes(returnChannel chan<- AttributesPair, customerId string, attributes map[string]string) {

	iC.makeAttributesRequest(returnChannel, http.MethodPost, []string{customersPath, customerId, customerAttributesPath}, attributes)
}

func (iC *InvisibleCollector) GetCustomerAttributes(returnChannel chan<- AttributesPair, customerId string) {
	iC.makeAttributesRequest(returnChannel, http.MethodGet, []string{customersPath, customerId, customerAttributesPath}, nil)
}

func (iC *InvisibleCollector) SetNewDebt(returnChannel chan<- DebtPair, newDebt Debt) {
	iC.makeDebtRequest(returnChannel, http.MethodPost, []string{debtsPath}, &newDebt, []fieldNamer{DebtNumber, DebtCustomerId, DebtType, DebtDate, DebtDueDate}, []fieldNamer{ItemName})
}

func (iC *InvisibleCollector) GetDebt(returnChannel chan<- DebtPair, debtId string) {
	iC.makeDebtRequest(returnChannel, http.MethodGet, []string{debtsPath, debtId}, nil, nil, nil)
}

func (iC *InvisibleCollector) GetCustomerDebts(returnChannel chan<- DebtListPair, customerId string) {

	debts := make([]Debt, 0)
	err := iC.makeRequest(&debts, http.MethodGet, []string{customersPath, customerId, customerDebtsPath}, nil)
	returnChannel <- DebtListPair{debts, err}
}

func (iC *InvisibleCollector) makeDebtRequest(returnChannel chan<- DebtPair, requestMethod string, pathFragments []string, requestDebt *Debt, debtMandatoryFields []fieldNamer, itemsMandatoryFields []fieldNamer) {

	if requestDebt != nil && len(itemsMandatoryFields) != 0 {
		if err := requestDebt.AssertItemsHaveFields(itemsMandatoryFields); err != nil {
			returnChannel <- DebtPair{Debt{}, nil}
			return
		}
	}

	var requestModel Modeler = nil
	if requestDebt != nil {
		requestModel = requestDebt
	}

	debt := Debt{}
	err := iC.makeModelRequest(&debt, requestMethod, pathFragments, requestModel, debtMandatoryFields)
	returnChannel <- DebtPair{debt, err}
}

func (iC *InvisibleCollector) makeAttributesRequest(returnChannel chan<- AttributesPair, requestMethod string, pathFragments []string, requestAttributes map[string]string) {
	attributes := make(map[string]string)
	var err error
	if len(requestAttributes) == 0 {
		err = iC.makeRequest(&attributes, requestMethod, pathFragments, nil)
	} else {
		err = iC.makeRequest(&attributes, requestMethod, pathFragments, requestAttributes)
	}

	returnChannel <- AttributesPair{attributes, err}
}

func (iC *InvisibleCollector) makeCustomerRequest(returnChannel chan<- CustomerPair, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) {

	customer := Customer{}
	err := iC.makeModelRequest(&customer, requestMethod, pathFragments, requestModel, mandatoryFields)
	returnChannel <- CustomerPair{customer, err}
}

func (iC *InvisibleCollector) makeCompanyRequest(returnChannel chan<- CompanyPair, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) {

	company := Company{}
	err := iC.makeModelRequest(&company, requestMethod, pathFragments, requestModel, mandatoryFields)
	returnChannel <- CompanyPair{company, err}
}

func (iC *InvisibleCollector) makeModelRequest(returnModel Modeler, requestMethod string, pathFragments []string, requestModel Modeler, mandatoryFields []fieldNamer) error {

	if requestModel != nil {
		if fieldErr := requestModel.AssertHasFields(mandatoryFields); fieldErr != nil {
			return fieldErr
		}
	}

	return iC.makeRequest(returnModel, requestMethod, pathFragments, requestModel)
}

func (iC *InvisibleCollector) makeRequest(returnData interface{}, requestMethod string, pathFragments []string, requestData interface{}) error {

	var requestBody []byte = nil
	if requestData != nil {

		requestJson, marshalErr := json.Marshal(requestData)
		if marshalErr != nil {
			return marshalErr
		}

		requestBody = requestJson
	}

	returnJson, requestErr := iC.makeJsonRequest(requestBody, requestMethod, pathFragments...)
	if requestErr != nil {
		return requestErr
	}

	return json.Unmarshal(returnJson, returnData)
}
