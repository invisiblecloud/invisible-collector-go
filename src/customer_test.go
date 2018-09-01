package ic

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

var customerNil *Customer = nil
var _ Modeler = customerNil
var _ json.Marshaler = customerNil
var _ json.Unmarshaler = customerNil

func TestCustomerMarshal(t *testing.T) {
	c := MakeCustomer()
	c.SetName(testString1)
	c.SetId(testString2)

	jsonBytes, err := json.Marshal(c)
	require.Nil(t, err)

	jsonString := string(jsonBytes)
	require.Contains(t, jsonString, testString1)
	require.Contains(t, jsonString, string(CustomerName))
	require.NotContains(t, jsonString, testString2)
	require.NotContains(t, jsonString, string(CustomerId))
}

func TestCustomerRoutableId(t *testing.T) {
	c := MakeCustomer()

	c.SetId("1")
	require.Equal(t, "1", c.RoutableId())

	c.SetExternalId("2")
	require.Equal(t, "1", c.RoutableId())

	c.UnsetField(CustomerId)
	require.Equal(t, "2", c.RoutableId())

}
