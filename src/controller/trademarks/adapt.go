package trademarks

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
)

type Adapt struct {
	ID      uint `json:",omitempty"`
	MakerID uint
	Name    string
	Status  bool
}

func checkFields(adapt Adapt) (model.Trademark, error) {
	element := model.Trademark{
		CustomModel: model.CustomModel{ID: adapt.ID},
		Name:        adapt.Name,
		Status:      adapt.Status,
		MakerID:     adapt.MakerID,
	}
	if element.MakerID == 0 {
		return element, controller.ErrTypeFormat("MakerID", "71", "''")
	}

	return element, nil
}

func HandleNewElement(body []byte) (model.Trademark, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Trademark{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleNewElements(body []byte) (newElement, updateElement []model.Trademark, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return newElement, updateElement, controller.ErrFormatIncorrect(err, example, body)
	}

	newElement = make([]model.Trademark, 0, len(adapts))
	updateElement = make([]model.Trademark, 0, len(adapts))

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

func HandleAdapt(body []byte) (updateElement []*model.Trademark, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Trademark, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
