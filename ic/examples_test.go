package ic_test

import (
	"fmt"
	"github.com/invisiblecloud/invisible-collector-go/ic"
)

const (
	exampleKey = "126618db3b74bc57b03ed22424933dd6a0732176c0eedd8e56f0937515eba826"
)

func Example() {
	iC, err := ic.NewInvisibleCollector(exampleKey, ic.InvisibleCollectorUri)
	if err != nil {
		panic(err)
	}

	// make a request
	var companies = make(chan ic.CompanyPair)
	go iC.GetCompany(companies)
	p := <-companies

	// check no errors occurred
	if p.Error != nil {
		panic(p.Error)
	}

	fmt.Println(p.Company)
}
