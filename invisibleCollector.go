package ic

import (
	"encoding/json"
	"github.com/invisiblecloud/invisible-collector-go/internal"
	"github.com/invisiblecloud/invisible-collector-go/models"
	"net/http"
)

const (
	companiesPath = "companies"
)

type CompanyPair struct {
	*models.Company
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

func (ic *InvisibleCollector) GetCompany(companies chan<- CompanyPair) {
	returnJson, requestErr := ic.MakeJsonRequest(nil, http.MethodGet, companiesPath)

	if requestErr != nil {
		companies <- CompanyPair{nil, requestErr}
		return
	}

	var company = models.Company{}
	err := json.Unmarshal(returnJson, &company)
	companies <- CompanyPair{&company, err}
}

func (ic *InvisibleCollector) SetCompany(companyUpdate models.Company, companies chan<- CompanyPair) {

	if fieldErr := companyUpdate.AssertHasFields(models.CompanyName, models.CompanyVatNumber); fieldErr != nil {
		companies <- CompanyPair{nil, fieldErr}
		return
	}

	requestJson, marshalErr := json.Marshal(companyUpdate)
	if marshalErr != nil {
		companies <- CompanyPair{nil, marshalErr}
		return
	}

	returnJson, requestErr := ic.MakeJsonRequest(requestJson, http.MethodPut, companiesPath)
	if requestErr != nil {
		companies <- CompanyPair{nil, requestErr}
		return
	}

	var company = models.Company{}
	err := json.Unmarshal(returnJson, &company)
	companies <- CompanyPair{&company, err}
}
