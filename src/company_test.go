package ic

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

// check that Company implements various interfaces
var companyNil *Company = nil
var _ Modeler = companyNil
var _ json.Marshaler = companyNil
var _ json.Unmarshaler = companyNil

func TestEmptyZeroValues(t *testing.T) {
	company := MakeCompany()
	require.Equal(t, "", company.City())

	require.Equal(t, false, company.NotificationsEnabled())
}

func TestGettersAndSetters(t *testing.T) {
	company := MakeCompany()
	company.SetName(testString1)
	company.SetCountry(testString2)
	require.Equal(t, testString1, company.Name())
	require.Equal(t, testString2, company.Country())
}

func TestFieldExists(t *testing.T) {
	c := MakeCompany()
	require.False(t, c.FieldExists(CompanyName))

	c.SetName("Name")

	require.True(t, c.FieldExists(CompanyName))
}

func TestSetFieldToNil(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.SetFieldToNil(CompanyName)
	require.Equal(t, "", c.Name())
}

func TestUnsetField(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.UnsetField(CompanyName)
	require.Equal(t, "", c.Name())
}

func TestAssertHasFields(t *testing.T) {
	c := MakeCompany()
	require.Nil(t, c.AssertHasFields([]fieldNamer{}))

	c.SetName(testString1)
	require.Nil(t, c.AssertHasFields([]fieldNamer{CompanyName}))

	require.NotEqual(t, nil, c.AssertHasFields([]fieldNamer{CompanyName, CompanyAddress}))
}

func TestDeepCopy(t *testing.T) {
	c1 := MakeCompany()
	c1.SetName(testString1)

	c2 := c1.deepCopy()
	require.Equal(t, testString1, c2.Name())

	c1.SetName(testString2)
	require.Equal(t, testString1, c2.Name())

}

func TestCompanyMarshal(t *testing.T) {
	c := MakeCompany()
	c.SetName(testString1)
	c.SetAddress(testString2)
	c.SetFieldToNil(CompanyCity)
	c.fields[CompanyNotificationsEnabled.fieldName()] = true

	jsonBytes, err := json.Marshal(c)
	require.Nil(t, err)

	jsonString := string(jsonBytes)
	require.Contains(t, jsonString, testString1)
	require.Contains(t, jsonString, string(CompanyName))
	require.Contains(t, jsonString, testString2)
	require.Contains(t, jsonString, string(CompanyAddress))
	require.Contains(t, jsonString, "null")
	require.Contains(t, jsonString, string(CompanyCity))
	require.NotContains(t, jsonString, CompanyNotificationsEnabled.fieldName())
	require.NotContains(t, jsonString, "true")
}

func TestUnmarshal(t *testing.T) {
	const companyJson = `
{
	"name": "test-name",
	"city": null
}
`

	companyJsonBytes := []byte(companyJson)
	var c Company
	err := json.Unmarshal(companyJsonBytes, &c)
	require.Nil(t, err)

	require.Equal(t, c.Name(), "test-name")
	require.False(t, c.FieldExists(CompanyCity))
	require.False(t, c.FieldExists(CompanyId))
}
