package model

type Supplier struct {
	CustomModel
	BusinessName string
	DaysPayment  string
	Nit          string
	PaymenTerm   string

	Cupo     uint
	Discount uint

	LeadTimeBaq uint
	LeadTimeBog uint

	NameContact  string
	EmailContact string
	PhoneContact string

	Status   bool
	Location string

	Purchases []Purchase
}
