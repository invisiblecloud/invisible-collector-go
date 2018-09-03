package ic

// Represent the Item fields.
//
// Can be used as an argument in various methods related to field manipulation.
const (
	ItemName        modelField = "name"
	ItemDescription modelField = "description"
	ItemQuantity    modelField = "quantity"
	ItemVat         modelField = "vat"
	ItemPrice       modelField = "price"
)

// The Item model.
//
// Item is a part of the Debt model
type Item struct {
	model
}

// The Item constructor
func MakeItem() Item {
	return Item{makeModel()}
}

func (i *Item) SetName(name string) {
	i.fields[string(ItemName)] = name
}

func (i *Item) Name() string {
	return i.getString(ItemName)
}

func (i *Item) SetDescription(description string) {
	i.fields[string(ItemDescription)] = description
}

func (i *Item) Description() string {
	return i.getString(ItemDescription)
}

func (i *Item) SetQuantity(quantity float64) {
	i.fields[string(ItemQuantity)] = quantity
}

func (i *Item) Quantity() float64 {
	return i.getFloat64(ItemQuantity)
}

// set the vat in percentage values (0 to 100.0)
func (i *Item) SetVat(vat float64) {
	i.fields[string(ItemVat)] = vat
}

func (i *Item) Vat() float64 {
	return i.getFloat64(ItemVat)
}

func (i *Item) SetPrice(price float64) {
	i.fields[string(ItemPrice)] = price
}

func (i *Item) Price() float64 {
	return i.getFloat64(ItemPrice)
}

func (i *Item) deepCopy() Item {
	return Item{i.shallowCopy()}
}
