package categories

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type AdaptCategoryTwo struct {
	ID            uint `json:",omitempty"`
	Name          string
	Status        bool
	CategoryOneID uint
}

func checkFieldsTwo(adapt AdaptCategoryTwo) (element model.CategoryTwo, err error) {
	element = model.CategoryTwo{
		CustomModel:   model.CustomModel{ID: adapt.ID},
		Name:          adapt.Name,
		Status:        adapt.Status,
		CategoryOneID: adapt.CategoryOneID,
	}
	return element, nil
}

func HandleNewElementTwo(body []byte) (model.CategoryTwo, error) {
	adapt := AdaptCategoryTwo{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.CategoryTwo{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFieldsTwo(adapt)
}

func HandleAdaptTwo(body []byte) (updateElement []*model.CategoryTwo, err error) {
	adapts := []AdaptCategoryTwo{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.CategoryTwo, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFieldsTwo(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
