package ic

// The FindCustomers model.
//
// Used to search for customers obeying specific restrictions
type FindCustomers struct {
	model
}

// Represent the FindCustomers fields.
//
// Can be used as an argument in various methods related to field manipulation.
const (
	FindCustomersExternalId modelField = "externalId"
	FindCustomersVatNumber  modelField = "vat"
	FindCustomersEmail      modelField = "email"
	FindCustomersPhone      modelField = "phone"
)

func MakeFindCustomer() FindCustomers {
	return FindCustomers{makeModel()}
}

func (c *FindCustomers) VatNumber() string {
	return c.getString(FindCustomersVatNumber)
}

func (c *FindCustomers) SetVatNumber(vatNumber string) {
	c.fields[string(FindCustomersVatNumber)] = vatNumber
}

func (c *FindCustomers) Phone() string {
	return c.getString(FindCustomersPhone)
}

func (c *FindCustomers) SetPhone(phone string) {
	c.fields[string(FindCustomersPhone)] = phone
}

func (c *FindCustomers) Email() string {
	return c.getString(FindCustomersEmail)
}

func (c *FindCustomers) SetEmail(email string) {
	c.fields[string(FindCustomersEmail)] = email
}

func (c *FindCustomers) SetExternalId(externalId string) {
	c.fields[string(FindCustomersExternalId)] = externalId
}

func (c *FindCustomers) ExternalId() string {
	return c.getString(FindCustomersExternalId)
}

func (c *FindCustomers) AsStringMap() map[string]string {
	ret := make(map[string]string)

	for k, v := range c.fields {
		ret[k] = v.(string)
	}

	return ret
}
