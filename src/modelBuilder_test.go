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

func (m *testModelBuilder) buildModel() model {
	fieldsCopy := internal.MapCopy(m.fields)
	internal.MapRemoveNils(fieldsCopy)
	return model{fieldsCopy}
}

func buildTestompanyModelBuilder() *testModelBuilder {
	m := makeTestModelBuilder()

	m.addField("name", "test-name")
	m.addField("vatNumber", "1234")
	m.addField("city", nil)

	return m
}
