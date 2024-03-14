package receptions

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ErrFormat2 = "Formato incorrecto"
	ErrFormat  = "Formato incorrecto"
)

func ErrTypeFormat2(field string, example, value interface{}) error {
	return fmt.Errorf(
		"%s: el campo %s deberia tener un formato %v en vez de %v",
		ErrFormat, field, example, value,
	)
}

type Adapt struct {
	ID        uint `json:",omitempty"`
	ArticleID uint

	Count   uint
	Batch   string
	Missing uint  // faltantes
	Refund  uint  // devoluciones
	Reason  uint8 // [1, 2, 3]
	Date    time.Time
}

func checkFields2(adapt Adapt) (element model.ReceptionArt, err error) {
	element = model.ReceptionArt{
		CustomModel: model.CustomModel{ID: adapt.ID},
		ArticleID:   adapt.ArticleID,
		Count:       adapt.Count,
		Batch:       adapt.Batch,
		Missing:     adapt.Missing,
		Refund:      adapt.Refund,
		Reason:      adapt.Reason,
		Date:        adapt.Date,
	}

	if element.ArticleID == 0 {
		return element, fmt.Errorf("La reception necesita tener un article id")
	}

	if element.GetId() != 0 {
		return model.ReceptionArt{},
			fmt.Errorf("no se puede modificar una reception")
	}
	if element.Count == 0 {
		return model.ReceptionArt{}, fmt.Errorf("no se pueden agregar receptiones con `Count 0`")
	}

	return element, nil
}

func handleNewElement2(body []byte) (model.ReceptionArt, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.ReceptionArt{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields2(adapt)
}

func handleNewElements2(body []byte) ([]model.ReceptionArt, error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return nil, controller.ErrFormatIncorrect(err, example, body)
	}
	elements := make([]model.ReceptionArt, len(adapts))

	for i, adapt := range adapts {
		element, err := checkFields2(adapt)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}
	return elements, nil
}
