package invisibleCollector

import (
	"encoding/json"
	"errors"
)

type modelField string

func (m modelField) fieldName() string {
	return string(m)
}

type fieldNamer interface {
	fieldName() string
}

type model struct {
	fields map[string]interface{}
}

func makeModel() model {
	return model{make(map[string]interface{})}
}

func (m *model) UnsetAllFields() {
	m.fields = make(map[string]interface{})
}

func (m *model) UnsetNullFields() {
	for k, v := range m.fields {
		if v == nil {
			delete(m.fields, k)
		}
	}
}

func (m *model) UnsetField(field fieldNamer) bool {
	fieldName := field.fieldName()
	if _, ok := m.fields[string(fieldName)]; ok {
		delete(m.fields, fieldName)
		return true
	}

	return false
}

func (m *model) SetFieldToNull(field fieldNamer) {
	m.fields[field.fieldName()] = nil
}

func (m model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.fields)
}

func (m *model) UnmarshalJSON(jsonString []byte) (err error) {
	err = json.Unmarshal(jsonString, &m.fields)
	m.UnsetNullFields()
	return err
}

func (m *model) AssertHasFields(requiredFields ...fieldNamer) error {
	for _, requiredField := range requiredFields {
		fieldName := requiredField.fieldName()
		if _, ok := m.fields[fieldName]; !ok {
			return errors.New("Missing field: " + fieldName)
		}
	}

	return nil
}
