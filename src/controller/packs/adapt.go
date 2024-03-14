package packs

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

func checkFields(adapt Adapt) (element model.Pack, err error) {
	element = model.Pack{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Status:      adapt.Status,
	}
	return element, nil
}

func HandleNewElement(body []byte) (model.Pack, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Pack{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleAdapt(body []byte) (updateElement []*model.Pack, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Pack, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
