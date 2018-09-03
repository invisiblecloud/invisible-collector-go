package ic

import (
	"encoding/json"
	"time"
)

// Represent the Debt fields.
//
// Can be used as an argument in various methods related to field manipulation.
const (
	DebtId         modelField = "gid"
	DebtNumber     modelField = "number"
	DebtCustomerId modelField = "customerId"
	DebtType       modelField = "type"
	DebtStatus     modelField = "status"
	DebtDate       modelField = "date"
	DebtDueDate    modelField = "dueDate"
	DebtNetTotal   modelField = "netTotal"
	DebtTax        modelField = "tax"
	DebtGrossTotal modelField = "grossTotal"
	DebtCurrency   modelField = "currency"
	DebtItems      modelField = "items"
	DebtAttributes modelField = "attributes"
)

const (
	dateFormat = "2006-01-02"
)

// The debt model
type Debt struct {
	model
}

// The model constructor
func MakeDebt() Debt {
	return Debt{makeModel()}
}

func (d Debt) MarshalJSON() ([]byte, error) {
	clone := Debt{d.shallowCopy()}
	clone.UnsetField(CustomerId)

	clone.tryFormatDateString(DebtDate)
	clone.tryFormatDateString(DebtDueDate)

	return clone.model.MarshalJSON()
}

func (d *Debt) UnmarshalJSON(jsonString []byte) error {
	if err := json.Unmarshal(jsonString, &d.fields); err != nil {
		return err
	}

	d.unsetNilFields()

	if err := d.tryUnformatDateString(DebtDate); err != nil {
		return err
	}

	if err := d.tryUnformatDateString(DebtDueDate); err != nil {
		return err
	}

	if d.FieldExists(DebtItems) {
		rawItems := d.fields[string(DebtItems)].([]interface{})
		items := make([]Item, len(rawItems))
		for i, v := range rawItems {
			item := Item{model{v.(map[string]interface{})}}
			item.unsetNilFields()
			items[i] = item
		}

		d.fields[string(DebtItems)] = items
	}

	if d.FieldExists(DebtAttributes) {
		rawAttributes := d.fields[string(DebtAttributes)].(map[string]interface{})
		attributes := make(map[string]string)
		for k, v := range rawAttributes {
			attributes[k] = v.(string)
		}

		d.fields[string(DebtAttributes)] = attributes
	}

	return nil
}

// Set the debt ID.
//
// This value is not sent but is used simply to identify the debt in invisible collector's database.
func (d *Debt) SetId(id string) {
	d.fields[string(DebtId)] = id
}

func (d *Debt) Id() string {
	return d.getString(DebtId)
}

func (d *Debt) SetNumber(number string) {
	d.fields[string(DebtNumber)] = number
}

func (d *Debt) Number() string {
	return d.getString(DebtNumber)
}

// The customer to whom the debt is issued
func (d *Debt) SetCustomerId(customerId string) {
	d.fields[string(DebtCustomerId)] = customerId
}

func (d *Debt) CustomerId() string {
	return d.getString(DebtCustomerId)
}

// Can be one of:
// "FT" - Normal invoice;
// "FS" - Simplified invoice;
// "SD" - Standard debt;
func (d *Debt) SetType(debtType string) {
	d.fields[string(DebtType)] = debtType
}

func (d *Debt) Type() string {
	return d.getString(DebtType)
}

// Can be one of:
// "PENDING" - the default value;
// "PAID";
// "CANCELLED";
func (d *Debt) SetStatus(status string) {
	d.fields[string(DebtStatus)] = status
}

func (d *Debt) Status() string {
	return d.getString(DebtStatus)
}

// The date where only the year, month and day is considered (YYYY-MM-DD)
func (d *Debt) SetDate(date time.Time) {
	d.fields[string(DebtDate)] = date
}

func (d *Debt) Date() time.Time {
	return d.getDate(DebtDate)
}

// The due date where only the year, month and day is considered (YYYY-MM-DD)
func (d *Debt) SetDueDate(dueDate time.Time) {
	d.fields[string(DebtDueDate)] = dueDate
}

func (d *Debt) DueDate() time.Time {
	return d.getDate(DebtDueDate)
}

func (d *Debt) SetNetTotal(netTotal float64) {
	d.fields[string(DebtNetTotal)] = netTotal
}

func (d *Debt) NetTotal() float64 {
	return d.getFloat64(DebtNetTotal)
}

func (d *Debt) SetTax(tax float64) {
	d.fields[string(DebtTax)] = tax
}

func (d *Debt) Tax() float64 {
	return d.getFloat64(DebtTax)
}

func (d *Debt) SetGrossTotal(grossTotal float64) {
	d.fields[string(DebtGrossTotal)] = grossTotal
}

func (d *Debt) GrossTotal() float64 {
	return d.getFloat64(DebtGrossTotal)
}

// currency must be in the ISO 4217 currency code format.
func (d *Debt) SetCurrency(currency string) {
	d.fields[string(DebtCurrency)] = currency
}

func (d *Debt) Currency() string {
	return d.getString(DebtCurrency)
}

// Add an Item model to the end of the items list
func (d *Debt) AddItem(item Item) {
	d.fields[string(DebtItems)] = append(d.items(), item.deepCopy())
}

// get all the items of this debt
func (d *Debt) Items() []Item {
	items := d.items()
	clone := make([]Item, len(items))

	for i, v := range items {
		clone[i] = v.deepCopy()
	}

	return clone
}

func (d *Debt) getAttributes() map[string]string {
	if v := d.getField(DebtAttributes); v != nil {
		return v.(map[string]string)
	}

	attributes := make(map[string]string)
	d.fields[string(DebtAttributes)] = attributes
	return attributes
}

// Set an attribute
func (d *Debt) SetAttribute(key string, value string) {
	d.getAttributes()[key] = value
}

// Get all the debt's attributes
func (d *Debt) Attributes() map[string]string {
	clone := make(map[string]string)

	for k, v := range d.getAttributes() {
		clone[k] = v
	}

	return clone
}

// Assert that the items have the specified field names.
//
// Used mostly internally.
func (d *Debt) AssertItemsHaveFields(requiredFields []fieldNamer) error {
	for _, v := range d.items() {
		if err := v.AssertHasFields(requiredFields); err != nil {
			return err
		}
	}

	return nil
}

func (d *Debt) items() []Item {
	if v := d.getField(DebtItems); v != nil {
		return v.([]Item)
	}

	return make([]Item, 0)
}

func (d *Debt) tryUnformatDateString(field modelField) (err error) {
	if d.FieldExists(field) {
		d.fields[string(field)], err = time.Parse(dateFormat, d.getString(field))
	}

	return
}

func (d *Debt) tryFormatDateString(field modelField) {
	if d.FieldExists(field) {
		d.fields[string(field)] = d.getDate(field).Format(dateFormat)
	}
}
