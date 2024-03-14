package categories

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type AdaptCategoryOne struct {
	ID     uint `json:",omitempty"`
	Name   string
	Status bool
}

func checkFieldsOne(adapt AdaptCategoryOne) (element model.CategoryOne, err error) {
	element = model.CategoryOne{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Status:      adapt.Status,
	}

	return element, nil
}

func HandleNewElementOne(body []byte) (model.CategoryOne, error) {
	adapt := AdaptCategoryOne{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.CategoryOne{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFieldsOne(adapt)
}

func HandleAdapt(body []byte) (updateElement []*model.CategoryOne, err error) {
	adapts := []AdaptCategoryOne{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.CategoryOne, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFieldsOne(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
