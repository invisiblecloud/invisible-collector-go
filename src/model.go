package ic

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

type Modeler interface {
	AssertHasFields(requiredFields []fieldNamer) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(jsonString []byte) (err error)
	UnsetField(field fieldNamer) bool
	FieldExists(field fieldNamer) bool
	SetFieldToNil(field fieldNamer)
}

type model struct {
	fields map[string]interface{}
}

func makeModel() model {
	return model{make(map[string]interface{})}
}

func (m *model) deepCopy() model {
	copy := makeModel()
	for k, v := range m.fields {
		copy.fields[k] = v
	}

	return copy
}

func (m *model) FieldExists(field fieldNamer) bool {
	return m.getField(field.(modelField)) != nil
}

func (m *model) SetFieldToNil(field fieldNamer) {
	m.fields[field.fieldName()] = nil
}

func (m *model) UnsetField(field fieldNamer) bool {
	fieldName := field.fieldName()
	if _, ok := m.fields[string(fieldName)]; ok {
		delete(m.fields, fieldName)
		return true
	}

	return false
}

func (m model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.fields)
}

func (m *model) UnmarshalJSON(jsonString []byte) (err error) {
	err = json.Unmarshal(jsonString, &m.fields)
	m.unsetNullFields()
	return err
}

func (m *model) AssertHasFields(requiredFields []fieldNamer) error {
	for _, requiredField := range requiredFields {
		fieldName := requiredField.fieldName()
		if _, ok := m.fields[fieldName]; !ok {
			return errors.New("Missing field: " + fieldName)
		}
	}

	return nil
}

func (m *model) unsetNullFields() {
	for k, v := range m.fields {
		if v == nil {
			delete(m.fields, k)
		}
	}
}

func (m *model) getField(fieldName modelField) interface{} {
	if v, ok := m.fields[string(fieldName)]; ok {
		return v
	}

	return nil
}

func (m *model) getString(fieldName modelField) string {
	if v := m.getField(fieldName); v != nil {
		return v.(string)
	}

	return ""
}

func (m *model) getBool(fieldName modelField) bool {
	if v := m.getField(fieldName); v != nil {
		return v.(bool)
	}

	return false
}
