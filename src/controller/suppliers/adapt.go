package suppliers

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/controller/purchases"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type Adapt struct {
	ID           uint
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

	Purchases []purchases.Adapt
}

func checkFields(adapt Adapt) (element model.Supplier, err error) {
	element = model.Supplier{
		CustomModel:  model.CustomModel{ID: adapt.ID},
		BusinessName: adapt.BusinessName,
		DaysPayment:  adapt.DaysPayment,
		Nit:          adapt.Nit,
		PaymenTerm:   adapt.PaymenTerm,
		Cupo:         adapt.Cupo,
		Discount:     adapt.Discount,
		LeadTimeBaq:  adapt.LeadTimeBaq,
		LeadTimeBog:  adapt.LeadTimeBog,
		NameContact:  adapt.NameContact,
		EmailContact: adapt.EmailContact,
		PhoneContact: adapt.PhoneContact,
		Status:       adapt.Status,
		Location:     adapt.Location,
	}
	return element, nil
}

func HandleNewElement(body []byte) (model.Supplier, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Supplier{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleAdapt(body []byte) (updateElement []*model.Supplier, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Supplier, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
