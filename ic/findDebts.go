package ic

import "time"

// Represent the FindDebts fields.
//
// Can be used as an argument in various methods related to field manipulation.
const (
	FindDebtsNumber      modelField = "number"
	FindDebtsFromDate    modelField = "from_date"
	FindDebtsToDate      modelField = "to_date"
	FindDebtsFromDueDate modelField = "from_duedate"
	FindDebtsToDueDate   modelField = "to_duedate"
)

// The FindDebts model.
//
// Used to search for debts obeying specific restrictions
type FindDebts struct {
	model
}

// The FindDebts constructor
func MakeFindDebts() FindDebts {
	return FindDebts{makeModel()}
}

func (d *FindDebts) SetNumber(number string) {
	d.fields[string(FindDebtsNumber)] = number
}

func (d *FindDebts) Number() string {
	return d.getString(FindDebtsNumber)
}

// The date where only the year, month and day is considered (YYYY-MM-DD)
func (d *FindDebts) SetToDate(date time.Time) {
	d.fields[string(FindDebtsToDate)] = date
}

func (d *FindDebts) ToDate() time.Time {
	return d.getDate(FindDebtsToDate)
}

// The date where only the year, month and day is considered (YYYY-MM-DD)
func (d *FindDebts) SetFromDate(date time.Time) {
	d.fields[string(FindDebtsFromDate)] = date
}

func (d *FindDebts) FromDate() time.Time {
	return d.getDate(FindDebtsFromDate)
}

// The date where only the year, month and day is considered (YYYY-MM-DD)
func (d *FindDebts) SetToDueDate(date time.Time) {
	d.fields[string(FindDebtsToDueDate)] = date
}

func (d *FindDebts) ToDueDate() time.Time {
	return d.getDate(FindDebtsToDueDate)
}

// The date where only the year, month and day is considered (YYYY-MM-DD)
func (d *FindDebts) SetFromDueDate(date time.Time) {
	d.fields[string(FindDebtsFromDueDate)] = date
}

func (d *FindDebts) FromDueDate() time.Time {
	return d.getDate(FindDebtsFromDueDate)
}

func (d *FindDebts) QueryParams() map[string]string {
	queries := make(map[string]string)

	if d.FieldExists(FindDebtsNumber) {
		queries[string(FindDebtsNumber)] = d.Number()
	}

	for _, field := range []modelField{FindDebtsFromDate, FindDebtsToDate, FindDebtsFromDueDate, FindDebtsToDueDate} {
		if d.FieldExists(field) {
			date := d.getDate(field)
			queries[string(field)] = date.Format(dateFormat)
		}
	}

	return queries
}
