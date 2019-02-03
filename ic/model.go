package ic

import (
	"encoding/json"
	"errors"
	"github.com/invisiblecloud/invisible-collector-go/internal"
	"time"
)

type modelField string

func (m modelField) fieldName() string {
	return string(m)
}

type fieldNamer interface {
	fieldName() string
}

// An interface to box any model
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

// Check whether the specified field exists.
//
// The field exists if it is set and isn't nil.
func (m *model) FieldExists(field fieldNamer) bool {
	return m.getField(field.(modelField)) != nil
}

// Check wheter the specified field has been set
//
// return true if field contains a value or nil
func (m *model) ContainsField(field fieldNamer) bool {
	_, ok := m.fields[field.fieldName()]
	return ok
}

// Set the specified field to nil.
//
// The field if set to nil will be set if it isn't set. A "null" json value will be sent for this field in any seding api request.
func (m *model) SetFieldToNil(field fieldNamer) {
	m.fields[field.fieldName()] = nil
}

// Unset the specified field.
//
// An unset field is neither set or nil. No field or value will be sent corresponding to the unset field.
// All fields in an empty model/newly constructed model are unset.
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
	m.unsetNilFields()
	return err
}

// Assert the model has all the specified fields, returning an error in case any is missing.
//
// Mostly used internaly.
func (m *model) AssertHasFields(requiredFields []fieldNamer) error {
	for _, requiredField := range requiredFields {
		fieldName := requiredField.fieldName()
		if _, ok := m.fields[fieldName]; !ok {
			return errors.New("Missing field: " + fieldName)
		}
	}

	return nil
}

func (m *model) unsetNilFields() {
	internal.MapRemoveNils(m.fields)
}

func (m *model) getField(fieldName modelField) interface{} {
	return internal.MapGetValue(m.fields, string(fieldName))
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

func (m *model) getFloat64(fieldName modelField) float64 {
	if v := m.getField(fieldName); v != nil {
		return v.(float64)
	}

	return 0.0
}

func (m *model) getDate(fieldName modelField) time.Time {
	if v := m.getField(fieldName); v != nil {
		return v.(time.Time)
	}

	return time.Time{}
}

func makeModel() model {
	return model{make(map[string]interface{})}
}

func (m *model) shallowCopy() model {
	clone := makeModel()
	for k, v := range m.fields {
		clone.fields[k] = v
	}

	return clone
}
