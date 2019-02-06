package ic

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// check that Company implements various interfaces
var findDebtsNil *Company = nil
var _ Modeler = findDebtsNil
var _ json.Marshaler = findDebtsNil
var _ json.Unmarshaler = findDebtsNil

func TestQueryParams(t *testing.T) {
	findDebts := MakeFindDebts()
	findDebts.SetNumber("1234")
	findDebts.SetFieldToNil(FindDebtsToDueDate)
	time := time.Date(
		2016, 1, 1, 20, 34, 58, 651387237, time.UTC)
	findDebts.SetToDate(time)

	queries := map[string]string{
		"number":     "1234",
		"to_duedate": "",
		"to_date":    "2016-01-01",
	}

	assert.Equal(t, queries, findDebts.QueryParams())
}
