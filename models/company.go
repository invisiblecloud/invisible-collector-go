package models

const (
	CompanyName                 modelField = "name"
	CompanyVatNumber            modelField = "vatNumber"
	CompanyAddress              modelField = "address"
	CompanyZipCode              modelField = "zipCode"
	CompanyCity                 modelField = "city"
	CompanyCountry              modelField = "country"
	CompanyGid                  modelField = "gid"
	CompanyNotificationsEnabled modelField = "notificationsEnabled"
)

type Company struct {
	model
}

func MakeCompany() Company {
	return Company{makeModel()}
}

func (c *Company) SetName(name string) {
	c.fields[string(CompanyName)] = name
}

func (c *Company) Name() string {
	return c.fields[string(CompanyName)].(string)
}

func (c *Company) SetVatNumber(vatNumber string) {
	c.fields[string(CompanyVatNumber)] = vatNumber
}

func (c *Company) VatNumber() string {
	return c.fields[string(CompanyVatNumber)].(string)
}

func (c *Company) SetAddress(address string) {
	c.fields[string(CompanyAddress)] = address
}

func (c *Company) Address() string {
	return c.fields[string(CompanyAddress)].(string)
}

func (c *Company) SetZipCode(ZipCode string) {
	c.fields[string(CompanyZipCode)] = ZipCode
}

func (c *Company) ZipCode() string {
	return c.fields[string(CompanyZipCode)].(string)
}

func (c *Company) SetCity(City string) {
	c.fields[string(CompanyCity)] = City
}

func (c *Company) City() string {
	return c.fields[string(CompanyCity)].(string)
}

func (c *Company) SetCountry(Country string) {
	c.fields[string(CompanyCountry)] = Country
}

func (c *Company) Country() string {
	return c.fields[string(CompanyCountry)].(string)
}

func (c *Company) SetGid(Gid string) {
	c.fields[string(CompanyGid)] = Gid
}

func (c *Company) Gid() string {
	return c.fields[string(CompanyGid)].(string)
}

func (c *Company) SetNotificationsEnabled(NotificationsEnabled bool) {
	c.fields[string(CompanyNotificationsEnabled)] = NotificationsEnabled
}

func (c *Company) NotificationsEnabled() bool {
	return c.fields[string(CompanyNotificationsEnabled)].(bool)
}
