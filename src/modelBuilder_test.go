package ic

import (
	"encoding/json"
	"github.com/invisiblecloud/invisible-collector-go/internal"
	"time"
)

func buildJson(source interface{}) string {
	j, err := json.Marshal(source)
	if err != nil {
		panic(err)
	}

	return string(j)
}

type testModelBuilder struct {
	fields map[string]interface{}
}

func makeTestModelBuilder() *testModelBuilder {
	return &testModelBuilder{make(map[string]interface{})}
}

func (m *testModelBuilder) addField(key string, value interface{}) {
	m.fields[key] = value
}

func (m *testModelBuilder) buildJson() string {
	return buildJson(m.fields)
}

func (m *testModelBuilder) buildDebtJson() string {
	fieldsCopy := internal.MapCopy(m.fields)
	if date, ok := fieldsCopy[string(DebtDate)]; ok && date != nil {
		fieldsCopy[string(DebtDate)] = date.(time.Time).Format(dateFormat)
	}

	if date, ok := fieldsCopy[string(DebtDueDate)]; ok && date != nil {
		fieldsCopy[string(DebtDueDate)] = date.(time.Time).Format(dateFormat)
	}

	return buildJson(fieldsCopy)
}

func getJsonBits(m map[string]interface{}, excludeKeys ...string) []string {
	clone := internal.MapCopy(m)
	for _, key := range excludeKeys {
		delete(clone, key)
	}

	ss := make([]string, 0)
	for k, v := range clone {
		switch val := v.(type) {
		case time.Time:
			ss = append(ss, k, val.Format(dateFormat))
		default:
			vj, _ := json.Marshal(v)
			ss = append(ss, k, string(vj))
		}
	}

	return ss
}

func (m *testModelBuilder) getRequestJsonBits(excludeKeys ...string) []string {
	return getJsonBits(m.fields, excludeKeys...)
}

func (m *testModelBuilder) getDebtRequestJsonBits(excludeKeys ...string) []string {
	bits := make([]string, 0)
	clone := internal.MapCopy(m.fields)
	if attributes, ok := clone[string(DebtAttributes)]; ok && attributes != nil {
		for k, v := range attributes.(map[string]string) {
			bits = append(bits, k, v)
		}

		delete(clone, string(DebtAttributes))
	}

	if items, ok := clone[string(DebtItems)]; ok && items != nil {
		for _, v := range items.([]Item) {
			slice := getJsonBits(v.fields)
			bits = append(bits, slice...)
		}

		delete(clone, string(DebtItems))
	}

	return append(bits, getJsonBits(clone, excludeKeys...)...)
}

func (m *testModelBuilder) buildReturnModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	internal.MapRemoveNils(fieldsCopy)
	return model{fieldsCopy}
}

func (m *testModelBuilder) buildDebtReturnModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	internal.MapRemoveNils(fieldsCopy)

	if items, ok := m.fields[string(DebtItems)]; ok && items != nil {
		for _, item := range items.([]Item) {
			item.unsetNilFields()
		}
	}

	return model{fieldsCopy}
}

func (m *testModelBuilder) buildRequestModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	return model{fieldsCopy}
}

func buildTestCompanyModelBuilder() *testModelBuilder {
	m := makeTestModelBuilder()

	m.addField(string(CompanyName), "test-name")
	m.addField(string(CompanyVatNumber), "1234")
	m.addField(string(CompanyCity), nil)

	return m
}

func buildTestCustomerModelBuilder() (m *testModelBuilder, id string) {
	m = makeTestModelBuilder()

	id = "adad"
	m.addField(string(CustomerName), "test-name")
	m.addField(string(CustomerVatNumber), "1234")
	m.addField(string(CustomerCountry), "PT")
	m.addField(string(CustomerCity), nil)
	m.addField(string(CustomerId), id)

	return
}

func buildTestDebtModelBuilder() (m *testModelBuilder, id string) {
	m = makeTestModelBuilder()

	id = "adad"
	m.addField(string(DebtNumber), "1")
	m.addField(string(DebtId), id)
	m.addField(string(DebtCustomerId), "ffff")
	m.addField(string(DebtType), "FT")
	m.addField(string(DebtGrossTotal), 100.1)
	date := time.Date(2015, 02, 01, 0, 0, 0, 0, time.UTC)
	m.addField(string(DebtDate), date)
	m.addField(string(DebtDueDate), date.AddDate(1, 1, 1))
	m.addField(string(DebtTax), nil)

	item := MakeItem()
	item.SetName("name-1")
	item.SetFieldToNil(ItemDescription)
	m.addField(string(DebtItems), []Item{item})

	attributes := make(map[string]string)
	attributes["k"] = "v"
	m.addField(string(DebtAttributes), attributes)

	return
}
