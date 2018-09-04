package ic

// Represent the Customer fields.
//
// Can be used as an argument in various methods related to field manipulation.
const (
	CustomerName       modelField = "name"
	CustomerId         modelField = "gid"
	CustomerExternalId modelField = "externalId"
	CustomerVatNumber  modelField = "vatNumber"
	CustomerAddress    modelField = "address"
	CustomerZipCode    modelField = "zipCode"
	CustomerCity       modelField = "city"
	CustomerCountry    modelField = "country"
	CustomerEmail      modelField = "email"
	CustomerPhone      modelField = "phone"
)

// The customer model
type Customer struct {
	model
}

// The customer constructor
func MakeCustomer() Customer {
	return Customer{makeModel()}
}

func (c Customer) MarshalJSON() ([]byte, error) {
	clone := c.deepCopy()
	clone.UnsetField(CustomerId)
	return clone.model.MarshalJSON()
}

func (c *Customer) SetName(name string) {
	c.fields[string(CustomerName)] = name
}

func (c *Customer) Name() string {
	return c.getString(CustomerName)
}

// The external ID of the customer.
//
// A convenience construct that you can set with your own IDs. Can be sometimes used to identify the customer.
func (c *Customer) SetExternalId(externalId string) {
	c.fields[string(CustomerExternalId)] = externalId
}

func (c *Customer) ExternalId() string {
	return c.getString(CustomerExternalId)
}

func (c *Customer) VatNumber() string {
	return c.getString(CustomerVatNumber)
}

func (c *Customer) SetVatNumber(vatNumber string) {
	c.fields[string(CustomerVatNumber)] = vatNumber
}

func (c *Customer) Address() string {
	return c.getString(CustomerAddress)
}

func (c *Customer) SetAddress(address string) {
	c.fields[string(CustomerAddress)] = address
}

func (c *Customer) ZipCode() string {
	return c.getString(CustomerZipCode)
}

func (c *Customer) SetZipCode(zipCode string) {
	c.fields[string(CustomerZipCode)] = zipCode
}

func (c *Customer) City() string {
	return c.getString(CustomerCity)
}

func (c *Customer) SetCity(city string) {
	c.fields[string(CustomerCity)] = city
}

func (c *Customer) Country() string {
	return c.getString(CustomerCountry)
}

// The country should be the ISO 3166-1 country code
func (c *Customer) SetCountry(country string) {
	c.fields[string(CustomerCountry)] = country
}

func (c *Customer) Email() string {
	return c.getString(CustomerEmail)
}

func (c *Customer) SetEmail(email string) {
	c.fields[string(CustomerEmail)] = email
}

func (c *Customer) Phone() string {
	return c.getString(CustomerPhone)
}

func (c *Customer) SetPhone(phone string) {
	c.fields[string(CustomerPhone)] = phone
}

func (c *Customer) Id() string {
	return c.getString(CustomerId)
}

// Set the customer ID.
//
// This value is not sent but is used simply to identify the customer in invisible collector's database.
// The "externalId" can only be used for this purpose in some circumstances (where you would use Customer.SetExternalId)
func (c *Customer) SetId(id string) {
	c.fields[string(CustomerId)] = id
}

// Represents the "routable" id of the customer
//
// Used for example in GET requests. It is the customer ID if it exists, if not it's the external id, if neither are present it's the string zero value.
// Used in various InvisibleCollector methods related to customers, you should check each method where a simple customer ID string is needed if external ID are admissible to see wheter you can use this method.
func (c *Customer) RoutableId() string {
	id := c.Id()
	if id == "" {
		return c.ExternalId()
	}

	return id
}

func (c *Customer) deepCopy() Customer {
	return Customer{c.shallowCopy()}
}
