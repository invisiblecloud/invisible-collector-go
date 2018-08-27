package ic

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

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
