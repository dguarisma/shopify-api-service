package categories

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type AdaptCategoryThree struct {
	ID            uint `json:",omitempty"`
	Name          string
	Status        bool
	CategoryOneID uint
	CategoryTwoID uint
}

func checkFieldsThree(adapt AdaptCategoryThree) (element model.CategoryThree, err error) {
	element = model.CategoryThree{
		CustomModel:   model.CustomModel{ID: adapt.ID},
		Name:          adapt.Name,
		Status:        adapt.Status,
		CategoryOneID: adapt.CategoryOneID,
		CategoryTwoID: adapt.CategoryTwoID,
	}
	return element, nil
}

func HandleNewElementThree(body []byte) (model.CategoryThree, error) {
	adapt := AdaptCategoryThree{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.CategoryThree{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFieldsThree(adapt)
}

func HandleAdaptThree(body []byte) (updateElement []*model.CategoryThree, err error) {
	adapts := []AdaptCategoryThree{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.CategoryThree, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFieldsThree(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
