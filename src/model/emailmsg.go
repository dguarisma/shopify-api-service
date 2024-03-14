package model

type EmailMsg struct {
	Fecha     string
	Orden     uint
	Bodega    string
	Name      string
	NIT       string
	Email     string
	Telefono  string
	Productos []Producto
	Subtotal  float32
	Impuesto  float32
	Descuento float32
	Total     float32
	To        string
}

type Producto struct {
	Name string
	Sku  string
	Ean  string

	Count     uint
	BasePrice float32
	Tax       float32
	Discount  float32
	Bonus     uint
	Subtotal  float32
	Total     float32
}
