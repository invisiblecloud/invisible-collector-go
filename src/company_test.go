package ic

import (
	"encoding/json"
	"strings"
	"testing"
)

const (
	testString1 = "test-value-1"
	testString2 = "test-value-2"
)

// check that Company implements various interfaces
var companyNil *Company = nil
var _ Modeler = companyNil
var _ json.Marshaler = companyNil
var _ json.Unmarshaler = companyNil

func TestEmptyZeroValues(t *testing.T) {
	company := MakeCompany()
	if company.City() != "" {
		t.Fatalf("Empty model must return default value")
	}

	if company.NotificationsEnabled() != false {
		t.Fatalf("Empty model must return default value")
	}
}

func TestGettersAndSetters(t *testing.T) {
	company := MakeCompany()
	company.SetName(testString1)
	company.SetCountry(testString2)
	if company.Name() != testString1 {
		t.Fail()
	}

	if company.Country() != testString2 {
		t.Fail()
	}
}

func TestFieldExists(t *testing.T) {
	c := MakeCompany()
	if c.FieldExists(CompanyName) {
		t.Fail()
	}

	c.SetName("Name")

	if !c.FieldExists(CompanyName) {
		t.Fail()
	}
}

func TestSetFieldToNil(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.SetFieldToNil(CompanyName)
	if c.Name() != "" {
		t.Fail()
	}
}

func TestUnsetField(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.UnsetField(CompanyName)
	if c.Name() != "" {
		t.Fail()
	}
}

func TestAssertHasFields(t *testing.T) {
	c := MakeCompany()
	if c.AssertHasFields([]fieldNamer{}) != nil {
		t.Fail()
	}

	c.SetName(testString1)
	if c.AssertHasFields([]fieldNamer{CompanyName}) != nil {
		t.Fail()
	}

	if c.AssertHasFields([]fieldNamer{CompanyName, CompanyAddress}) == nil {
		t.Fail()
	}
}

func TestDeepCopy(t *testing.T) {
	c1 := MakeCompany()
	c1.SetName(testString1)

	c2 := c1.DeepCopy()
	if c2.Name() != testString1 {
		t.Fail()
	}

	c1.SetName(testString2)
	if c2.Name() != testString1 {
		t.Fail()
	}

}

func TestMarshal(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.SetAddress(testString2)
	c.SetFieldToNil(CompanyCity)

	jsonBytes, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	json := string(jsonBytes)
	if !(strings.Contains(json, testString1) &&
		strings.Contains(json, string(CompanyName)) &&
		strings.Contains(json, testString2) &&
		strings.Contains(json, string(CompanyAddress)) &&
		strings.Contains(json, "null") &&
		strings.Contains(json, string(CompanyCity))) {
		t.Fail()
	}
}

func TestUnmarshal(t *testing.T) {
	var companyJson = `
{
	"name": "test-name",
	"city": null
}
`

	companyJsonBytes := []byte(companyJson)
	var c Company
	var m Modeler = &c
	err := json.Unmarshal(companyJsonBytes, m)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != "test-name" || c.FieldExists(CompanyCity) || c.FieldExists(CompanyId) {
		t.Fail()
	}

}
