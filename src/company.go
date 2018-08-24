package ic

import "github.com/invisiblecloud/invisible-collector-go/internal"

const (
	CompanyName                 modelField = "name"
	CompanyVatNumber            modelField = "vatNumber"
	CompanyAddress              modelField = "address"
	CompanyZipCode              modelField = "zipCode"
	CompanyCity                 modelField = "city"
	CompanyCountry              modelField = "country"
	CompanyId                   modelField = "gid"
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
	return c.getString(CompanyName)
}

func (c *Company) SetVatNumber(vatNumber string) {
	c.fields[string(CompanyVatNumber)] = vatNumber
}

func (c *Company) VatNumber() string {
	return c.getString(CompanyVatNumber)
}

func (c *Company) SetAddress(address string) {
	c.fields[string(CompanyAddress)] = address
}

func (c *Company) Address() string {
	return c.getString(CompanyAddress)
}

func (c *Company) SetZipCode(zipCode string) {
	c.fields[string(CompanyZipCode)] = zipCode
}

func (c *Company) ZipCode() string {
	return c.getString(CompanyZipCode)
}

func (c *Company) SetCity(city string) {
	c.fields[string(CompanyCity)] = city
}

func (c *Company) City() string {
	return c.getString(CompanyCity)
}

func (c *Company) SetCountry(country string) {
	c.fields[string(CompanyCountry)] = country
}

func (c *Company) Country() string {
	return c.getString(CompanyCountry)
}

func (c *Company) Id() string {
	return c.getString(CompanyId)
}

func (c *Company) NotificationsEnabled() bool {
	return c.getBool(CompanyNotificationsEnabled)
}

func (c *Company) DeepCopy() Company {
	return Company{c.deepCopy()}
}

func (c *Company) MarshalJSON() ([]byte, error) {
	m := internal.MapSubmap(c.fields, string(CompanyName), string(CompanyVatNumber), string(CompanyAddress), string(CompanyZipCode), string(CompanyCity))
	return model{m}.MarshalJSON()
}
