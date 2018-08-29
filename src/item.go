package ic

const (
	ItemName        modelField = "name"
	ItemDescription modelField = "description"
	ItemQuantity    modelField = "quantity"
	ItemVat         modelField = "vat"
	ItemPrice       modelField = "price"
)

type Item struct {
	model
}

func MakeItem() Item {
	return Item{makeModel()}
}

func (i *Item) DeepCopy() Item {
	return Item{i.shallowCopy()}
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
