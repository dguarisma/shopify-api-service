package substances

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type Adapt struct {
	ID        uint `json:",omitempty"`
	ProductID uint
	Name      string
	Status    bool
}

func checkFields(adapt Adapt) (model.Substance, error) {
	element := model.Substance{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Status:      adapt.Status,
	}

	if adapt.ProductID != 0 {
		element.ProductID = &adapt.ProductID
	}
	return element, nil
}

func HandleNewElement(body []byte) (model.Substance, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Substance{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

/*
	func HandleNewElements(body []byte) (newElement, updateElement []model.Substance, err error) {
		adapts := []Adapt{}
		if err := json.Unmarshal(body, &adapts); err != nil {
			example, _ := json.MarshalIndent(adapts, "", "  ")
			return newElement, updateElement, controller.ErrFormatIncorrect(err, example, body)
		}

		newElement = make([]model.Substance, 0, len(adapts))
		updateElement = make([]model.Substance, 0, len(adapts))

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
*/
func HandleAdapt(body []byte) (updateElement []*model.Substance, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Substance, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
