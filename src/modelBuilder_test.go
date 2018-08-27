package ic

import (
	"encoding/json"
	"github.com/invisiblecloud/invisible-collector-go/internal"
)

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
	j, err := json.Marshal(m.fields)
	if err != nil {
		panic(err)
	}

	return string(j)
}

func (m *testModelBuilder) getRequestJsonBits(excludeKeys ...string) []string {

	clone := internal.MapCopy(m.fields)
	for _, key := range excludeKeys {
		delete(clone, key)
	}

	ss := make([]string, 0)
	for k, v := range clone {
		vj, _ := json.Marshal(v)
		ss = append(ss, k, string(vj))
	}

	return ss
}

func (m *testModelBuilder) buildReturnModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	internal.MapRemoveNils(fieldsCopy)
	return model{fieldsCopy}
}

func (m *testModelBuilder) buildRequestModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	return model{fieldsCopy}
}

func buildTestCompanyModelBuilder() *testModelBuilder {
	m := makeTestModelBuilder()

	m.addField("name", "test-name")
	m.addField("vatNumber", "1234")
	m.addField("city", nil)

	return m
}

func buildTestCustomerModelBuilder() (m *testModelBuilder, id string) {
	m = makeTestModelBuilder()

	id = "adad"
	m.addField("name", "test-name")
	m.addField("vatNumber", "1234")
	m.addField("country", "PT")
	m.addField("city", nil)
	m.addField("gid", id)

	return
}
