package warehouses

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/controller/products"
	"desarrollosmoyan/lambda/src/controller/purchases"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type Adapt struct {
	ID         uint `json:",omitempty"`
	Name       string
	Department string
	City       string
	Location   string
	Status     bool
	Products   []products.Adapt
	Purchase   []purchases.Adapt
}

func checkFields(adapt Adapt) (element model.Warehouse, err error) {
	element = model.Warehouse{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Department:  adapt.Department,
		City:        adapt.City,
		Location:    adapt.Location,
		Status:      adapt.Status,
	}

	return element, nil
}

func HandleNewElement(body []byte) (model.Warehouse, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Warehouse{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleNewElements(body []byte) (newElement, updateElement []model.Warehouse, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return newElement, updateElement, controller.ErrFormatIncorrect(err, example, body)
	}

	newElement = make([]model.Warehouse, 0, len(adapts))
	updateElement = make([]model.Warehouse, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return newElement, updateElement, err
		}

		if element.ID == 0 {
			newElement = append(newElement, element)
			continue
		}

		updateElement = append(updateElement, element)
	}
	return newElement, updateElement, nil
}

func HandleAdapt(body []byte) (updateElement []*model.Warehouse, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Warehouse, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}