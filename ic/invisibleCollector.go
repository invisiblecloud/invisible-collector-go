// The go library used to communicate with Invisible Collector services.
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
)

const (
	// The default Invisible Collector URL
	InvisibleCollectorUri = "https://api.invisiblecollector.com/"
)

// Returned by any InvisibleCollector methods dealing with companies.
//
// In case of a successful method request the Company field contains the up-to-date company information and Error is nil,
// in case of error Error is non nil and Company is in a invalid state.
type CompanyPair struct {
	Company Company
	Error   error
}

// Returned by any InvisibleCollector methods dealing with customers.
//
// In case of a successful method request the Customer field contains the up-to-date customer information and Error is nil,
// in case of error Error is non nil and Customer is in a invalid state.
type CustomerPair struct {
	Customer Customer
	Error    error
}

// Returned by any InvisibleCollector methods dealing with customer attributes.
//
// In case of a successful method request the Attributes field contains the up-to-date customer attributes and Error is nil,
// in case of error Error is non nil and Attributes is in a invalid state.
type AttributesPair struct {
	Attributes map[string]string
	Error      error
}

// Returned by any InvisibleCollector methods dealing with debts.
//
// In case of a successful method request the Debt field contains the up-to-date customer debt information and Error is nil,
// in case of error Error is non nil and Debt is in a invalid state.
type DebtPair struct {
	Debt  Debt
	Error error
}

// Returned by any InvisibleCollector.GetCustomerDebts()
//
// In case of a successful method request the Debts field contains the up-to-date list of customer debts and Error is nil,
// in case of error Error is non nil and Debt is in a invalid state.
type DebtListPair struct {
	Debts []Debt
	Error error
}

// Returned by any InvisibleCollector.GetFindCustomers()
//
// In case of a successful method request the Customers field contains the up-to-date list of customers and Error is nil,
// in case of error Error is non nil and Customers is in a invalid state.
type CustomerListPair struct {
	Customers []Customer
	Error     error
}

// The entry point for using this library.
// Contains all the implemented API request by this library that are used to communicate with Invisible Collector.
//
// Any methods API request methods of this struct only ever return a single "pair", with the updated model and error as nil in case of a successful request,
// or a zero value model (unusable) and a non nil error on any error.
type InvisibleCollector struct {
	apiRequest
}

// Constructor for this class.
//
// Use your company API key. You can use the InvisibleCollectorUri constant to use the default Invisible Collector URI.
// Only http and https URI schemes are supported.
func NewInvisibleCollector(apiKey string, apiUrl string) (*InvisibleCollector, error) {
	requests, err := newApiRequest(apiKey, apiUrl)
	if err != nil {
		return nil, err
	}

	return &InvisibleCollector{*requests}, nil
}

// Get the company information.
func (iC *InvisibleCollector) GetCompany(returnChannel chan<- CompanyPair) {

	iC.makeCompanyRequest(returnChannel, http.MethodGet, []string{companiesPath}, nil)
}

// Update company information.
//
// Only "address", "zipCode" and "city" can be changed. "name" and "vatNumber" are mandatory fields that are used for validation. You can use GetCompany first to get the validation fields.
func (iC *InvisibleCollector) SetCompany(returnChannel chan<- CompanyPair, companyUpdate Company) {

	iC.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath}, &companyUpdate)
}

// Enable or disable company notifications
//
// Set to true to enable, false to disable notification
func (iC *InvisibleCollector) SetCompanyNotifications(returnChannel chan<- CompanyPair, enableNotifications bool) {
	var notificationsPath string
	if enableNotifications {
		notificationsPath = "enableNotifications"
	} else {
		notificationsPath = "disableNotifications"
	}

	iC.makeCompanyRequest(returnChannel, http.MethodPut, []string{companiesPath, notificationsPath}, nil)
}

// Create a new customer
//
// "name", "vatNumber" and "country" are mandatory fields
func (iC *InvisibleCollector) SetNewCustomer(returnChannel chan<- CustomerPair, newCustomer Customer) {
	iC.makeCustomerRequest(returnChannel, http.MethodPost, []string{customersPath}, &newCustomer)
}

// Update customer.
//
// Id can be the customer ID ("gid") or the customer external ID ("externalId"). Either or both must be present in the
// updatedCustomer argument. If both are present the customer ID is used to identify the customer
//
// "country" is a mandatory field
func (iC *InvisibleCollector) SetCustomer(returnChannel chan<- CustomerPair, updatedCustomer Customer) {
	iC.makeCustomerRequest(returnChannel, http.MethodPut, []string{customersPath, updatedCustomer.RoutableId()}, &updatedCustomer)
}

// Get customer information
//
// customerId can be either the ID ("gid") or external ID ("externalId").
func (iC *InvisibleCollector) GetCustomer(returnChannel chan<- CustomerPair, customerId string) {
	iC.makeCustomerRequest(returnChannel, http.MethodGet, []string{customersPath, customerId}, nil)
}

// Set the customer's attributes.
//
// New attributes are created, existing attributes are overwritten.
//
// customerId must be the customer ID.
func (iC *InvisibleCollector) SetCustomerAttributes(returnChannel chan<- AttributesPair, customerId string, attributes map[string]string) {

	iC.makeAttributesRequest(returnChannel, http.MethodPost, []string{customersPath, customerId, customerAttributesPath}, attributes)
}

// Get the customer's attributes
//
// customerId must be the customer ID.
func (iC *InvisibleCollector) GetCustomerAttributes(returnChannel chan<- AttributesPair, customerId string) {
	iC.makeAttributesRequest(returnChannel, http.MethodGet, []string{customersPath, customerId, customerAttributesPath}, nil)
}

// Create a new customer Debt.
//
// newDebt's "number", "customerId", "type", "date" and "dueDate" are mandatory. If the debt contains Items their "name" field is mandatory.
func (iC *InvisibleCollector) SetNewDebt(returnChannel chan<- DebtPair, newDebt Debt) {
	iC.makeDebtRequest(returnChannel, http.MethodPost, []string{debtsPath}, &newDebt)
}

// Get a debt by the debt's id.
func (iC *InvisibleCollector) GetDebt(returnChannel chan<- DebtPair, debtId string) {
	iC.makeDebtRequest(returnChannel, http.MethodGet, []string{debtsPath, debtId}, nil)
}

// Get all of the customer's debts.
func (iC *InvisibleCollector) GetCustomerDebts(returnChannel chan<- DebtListPair, customerId string) {

	debts := make([]Debt, 0)
	err := iC.makeRequest(&debts, http.MethodGet, []string{customersPath, customerId, customerDebtsPath}, nil)
	returnChannel <- DebtListPair{debts, err}
}

// Search for debts
func (iC *InvisibleCollector) GetFindDebts(returnChannel chan<- DebtListPair, findDebts FindDebts) {
	debts := make([]Debt, 0)
	err := iC.makeQueryRequest(&debts, http.MethodGet, []string{debtsPath, "find"}, findDebts.QueryParams())
	returnChannel <- DebtListPair{debts, err}
}

// Search for customers
func (iC *InvisibleCollector) GetFindCustomers(returnChannel chan<- CustomerListPair, findCustomers FindCustomers) {
	customers := make([]Customer, 0)
	err := iC.makeQueryRequest(&customers, http.MethodGet, []string{customersPath, "find"}, findCustomers.AsStringMap())
	returnChannel <- CustomerListPair{customers, err}
}

func (iC *InvisibleCollector) makeDebtRequest(returnChannel chan<- DebtPair, requestMethod string, pathFragments []string, requestDebt *Debt) {

	var requestModel Modeler = nil
	if requestDebt != nil {
		requestModel = requestDebt
	}

	debt := Debt{}
	err := iC.makeRequest(&debt, requestMethod, pathFragments, requestModel)
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

func (iC *InvisibleCollector) makeCustomerRequest(returnChannel chan<- CustomerPair, requestMethod string, pathFragments []string, requestModel Modeler) {

	customer := Customer{}
	err := iC.makeRequest(&customer, requestMethod, pathFragments, requestModel)
	returnChannel <- CustomerPair{customer, err}
}

func (iC *InvisibleCollector) makeCompanyRequest(returnChannel chan<- CompanyPair, requestMethod string, pathFragments []string, requestModel Modeler) {

	company := Company{}
	err := iC.makeRequest(&company, requestMethod, pathFragments, requestModel)
	returnChannel <- CompanyPair{company, err}
}

func (iC *InvisibleCollector) makeQueryRequest(returnData interface{}, requestMethod string, pathFragments []string, queryParams map[string]string) error {

	returnJson, requestErr := iC.makeUrlEncodedRequest(queryParams, requestMethod, pathFragments...)
	if requestErr != nil {
		return requestErr
	}

	return json.Unmarshal(returnJson, returnData)
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
