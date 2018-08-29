package ic

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDebtAddItems(t *testing.T) {
	d := MakeDebt()
	require.Equal(t, 0, len(d.items()))

	item := MakeItem()
	item.SetName(testString1)
	d.AddItem(item)
	require.Equal(t, 1, len(d.items()))
	require.Equal(t, testString1, d.Items()[0].Name())
}

func TestDebtItemsHaveMandatoryFields(t *testing.T) {
	d := MakeDebt()
	assert.Nil(t, d.AssertItemsHaveFields([]fieldNamer{ItemName}))

	item := MakeItem()
	item.SetName(testString1)
	item.SetDescription(testString2)
	d.AddItem(item)
	assert.Nil(t, d.AssertItemsHaveFields([]fieldNamer{ItemName, ItemDescription}))
	assert.Nil(t, d.AssertItemsHaveFields([]fieldNamer{}))
	assert.NotNil(t, d.AssertItemsHaveFields([]fieldNamer{ItemVat}))

	d.AddItem(MakeItem())
	assert.NotNil(t, d.AssertItemsHaveFields([]fieldNamer{ItemName}))
}

func TestSetAttributes(t *testing.T) {
	d := MakeDebt()

	d.SetAttribute("k", "v")
	require.Equal(t, "v", d.getAttributes()["k"])
}

func TestDebtUnmarshalJSON(t *testing.T) {
	const json = `{
  "number": "1",
  "customerId": "0d3987e3-a6df-422c-8722-3fde26eec9a8",
  "type": "FT",
  "status": "PENDING",
  "date": "2018-02-02",
  "dueDate": "2019-01-02",
  "netTotal": 1000.0,
  "tax": 200.1,
  "grossTotal": null,
  "currency": "EUR",
  "items": [
    {
      "name": "name-1",
      "description": "a debt item description",
      "quantity": 3.0,
      "vat": 23.0,
      "price": 15.2
    },
    {
      "name": "name-2",
      "description": null
    }
  ],
  "attributes": {
    "name_1": "attribute_1",
    "name_2": "attribute_2"
  }
}`

	d := MakeDebt()
	err := d.UnmarshalJSON([]byte(json))
	require.Nil(t, err)

	assert.Equal(t, "1", d.Number())
	assert.InDelta(t, 200.1, d.Tax(), 0.001)
	assert.False(t, d.FieldExists(DebtGrossTotal))

	// dates
	y, m, day := d.Date().Date() // implicit type assertion test
	assert.Equal(t, 2018, y)
	assert.Equal(t, time.February, m)
	assert.Equal(t, 2, day)
	assert.Equal(t, 0, d.Date().Hour())
	assert.Equal(t, 0, d.Date().Minute())
	assert.Equal(t, 0, d.Date().Second())
	assert.Equal(t, 0, d.Date().Nanosecond())

	// attributes
	assert.Equal(t, "attribute_1", d.Attributes()["name_1"]) // also implicitly testing correct type conversion

	// items
	items := d.Items() //also implicitly testing for type conversion
	require.Equal(t, 2, len(items))
	assert.Equal(t, "name-1", items[0].Name())
	assert.InDelta(t, 15.2, items[0].Price(), 0.01)
	assert.Equal(t, "name-2", items[1].Name())
	_, ok := items[1].fields[string(ItemDescription)]
	assert.False(t, ok)
	assert.False(t, items[1].FieldExists(ItemVat))

}

func TestDebtMarshalJSON(t *testing.T) {
	d := MakeDebt()
	d.SetNumber("1")
	d.SetGrossTotal(100.1)
	d.SetDate(time.Date(2018, time.February, 5, 1, 1, 1, 1, time.UTC))
	d.SetAttribute("k", "v")
	d.SetFieldToNil(DebtDueDate)

	i := MakeItem()
	i.SetName("name-1")
	i.SetFieldToNil(ItemDescription)
	d.AddItem(i)

	j, err := d.MarshalJSON()
	require.Nil(t, err)

	model := make(map[string]interface{})
	err2 := json.Unmarshal(j, &model)
	require.Nil(t, err2)

	assert.Equal(t, "1", model[string(DebtNumber)].(string))
	assert.InDelta(t, 100.1, model[string(DebtGrossTotal)].(float64), 0.0001)
	assert.Nil(t, model[string(DebtDueDate)])
	assert.Equal(t, "2018-02-05", model[string(DebtDate)].(string))
	assert.Equal(t, "v", model[string(DebtAttributes)].(map[string]interface{})["k"].(string))

	item := model[string(DebtItems)].([]interface{})[0].(map[string]interface{}) // implicit test

	assert.Equal(t, "name-1", item[string(ItemName)])
	assert.Nil(t, item[string(ItemDescription)])
	_, ok := item[string(ItemVat)]
	assert.False(t, ok)

}
