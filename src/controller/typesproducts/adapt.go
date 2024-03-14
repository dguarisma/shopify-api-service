package typesproducts

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type Adapt struct {
	ID     uint `json:",omitempty"`
	Name   string
	Status bool
}

func checkFields(adapt Adapt) (model.TypesProduct, error) {
	element := model.TypesProduct{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Status:      adapt.Status,
	}
	return element, nil
}

func HandleNewElement(body []byte) (model.TypesProduct, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.TypesProduct{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleNewElements(body []byte) (newElement, updateElement []model.TypesProduct, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return newElement, updateElement, controller.ErrFormatIncorrect(err, example, body)
	}

	newElement = make([]model.TypesProduct, 0, len(adapts))
	updateElement = make([]model.TypesProduct, 0, len(adapts))

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

func HandleAdapt(body []byte) (updateElement []*model.TypesProduct, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.TypesProduct, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
