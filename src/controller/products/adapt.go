package products

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/converts"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
	"fmt"
)

const (
	ErrFormat = "Formato incorrecto"
)

func ErrTypeFormat(field string, example, value interface{}) error {
	return fmt.Errorf(
		"%s: el campo %s deberia tener un formato %v en vez de %v",
		ErrFormat, field, example, value,
	)
}

type Adapt struct {
	ID              string
	PackID          string
	MakerID         string
	TrademarkID     string
	TypesProductID  string
	CategoryOneID   string
	CategoryTwoID   string
	CategoryThreeID string
	WarehouseIDS    string
	SubstitutesIDS  string
	SubstancesIDS   string
	IDFloorProduct  uint

	Name        string
	SKU         string `json:"Sku"`
	Ean         string
	Iva         string `json:"iva"`
	Img         string `json:"UrlImage"`
	Quantity    string
	Weight      string
	Width       string
	Height      string
	Pack        string `json:"PackInfo"`
	PackUnit    string
	Depth       string
	Keywords    string
	MakerUnit   string
	Variation   string
	Wrapper     string
	WrapperUnit string
	HandlesBog  string
	HandlesBaq  string
	Status      bool
	IsTaxed     bool `json:"Taxed"`
}

func checkFields(product Adapt) (model.Product, error) {
	campo := ""
	iva, err := converts.StringToFloat32(product.Iva)
	if err != nil {
		campo = "Iva"
		return model.Product{}, ErrTypeFormat(campo, "21.4", product.Iva)
	}

	weight, err := converts.StringToFloat32(product.Weight)
	if err != nil {
		campo = "Weight"
		return model.Product{}, ErrTypeFormat(campo, "21.4", product.Weight)
	}

	width, err := converts.StringToFloat32(product.Width)
	if err != nil {
		campo = "Width"
		return model.Product{}, ErrTypeFormat(campo, "31.4", product.Width)
	}

	height, err := converts.StringToFloat32(product.Height)
	if err != nil {
		campo = "Height"
		return model.Product{}, ErrTypeFormat(campo, "31.4", product.Height)
	}

	quantity, err := converts.StringToUint64(product.Quantity)
	if err != nil {
		campo = "Quantity"
		return model.Product{}, ErrTypeFormat(campo, "31", product.Quantity)
	}

	packUnit, err := converts.StringToUint64(product.PackUnit)
	if err != nil {
		campo = "PackUnit"
		return model.Product{}, ErrTypeFormat(campo, "100", product.PackUnit)
	}

	makeID, err := converts.StringToUint(product.MakerID)
	if err != nil {
		campo = "MakerID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.MakerID)
	}

	packID, err := converts.StringToUint(product.PackID)
	if err != nil {
		campo = "PackID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.PackID)
	}

	trademarkID, err := converts.StringToUint(product.TrademarkID)
	if err != nil {
		campo = "TrademarkID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.TrademarkID)
	}

	typesProductID, err := converts.StringToUint(product.TypesProductID)
	if err != nil {
		campo = "TypesProductID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.TypesProductID)
	}

	categoryOneID, err := converts.StringToUint(product.CategoryOneID)
	if err != nil {
		campo = "CategoryOneID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.CategoryOneID)
	}

	categoryTwoID, err := converts.StringToUint(product.CategoryTwoID)
	if err != nil {
		campo = "CategoryTwoID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.CategoryTwoID)
	}

	categoryThreeID, err := converts.StringToUint(product.CategoryThreeID)
	if err != nil {
		campo = "CategoryThreeID"
		return model.Product{}, ErrTypeFormat(campo, "100", product.CategoryThreeID)
	}

	/* makeUnit, err := converts.StringToUint64(product.MakerUnit)
	if err != nil {
		campo = "MakerUnit"
		return model.Product{}, ErrTypeFormat(campo, "100", product.MakerUnit)
	} */

	depth, err := converts.StringToUint64(product.Depth)
	if err != nil {
		campo = "Depth"
		return model.Product{}, ErrTypeFormat(campo, "100", product.Depth)
	}

	productDb := model.Product{
		Name:     product.Name,
		Sku:      product.SKU,
		Ean:      product.Ean,
		Taxed:    product.IsTaxed,
		Iva:      iva,
		UrlImage: product.Img,
		Quantity: quantity,
		Weight:   weight,
		Width:    width,
		Height:   height,
		PackUnit: packUnit,
		Depth:    depth,
		Keywords: product.Keywords,
		PackInfo: product.Pack,

		HandlesBog:      product.HandlesBog,
		HandlesBaq:      product.HandlesBaq,
		IDFloorProduct:  product.IDFloorProduct,
		MakerID:         makeID,
		PackID:          packID,
		TrademarkID:     trademarkID,
		TypesProductID:  typesProductID,
		CategoryOneID:   categoryOneID,
		CategoryTwoID:   categoryTwoID,
		CategoryThreeID: categoryThreeID,

		MakerUnit:   product.MakerUnit,
		Variation:   product.Variation,
		Wrapper:     product.Wrapper,
		WrapperUnit: product.WrapperUnit,
		Status:      product.Status,
		Warehouses:  []*model.Warehouse{},
		Substance:   []*model.Substance{},
		Substitutes: []*model.Product{},
	}

	if product.ID != "" && product.ID != "0" {
		id, err := converts.StringToUint(product.ID)
		if err != nil {
			campo = "ID"
			return model.Product{}, ErrTypeFormat(campo, "100", product.ID)
		}
		productDb.ID = *id
	}

	subtances, err := converts.StringToArrUint(product.SubstancesIDS)
	if err != nil {
		campo = "SubstancesIDS"
		return model.Product{}, ErrTypeFormat(campo, "11, 14, 3", product.SubstancesIDS)
	}

	warehouses, err := converts.StringToArrUint(product.WarehouseIDS)
	if err != nil {
		campo = "WarehouseIDS"
		return model.Product{}, ErrTypeFormat(campo, "11, 14, 3", product.WarehouseIDS)
	}
	substitues, err := converts.StringToArrUint(product.SubstitutesIDS)
	if err != nil {
		campo = "SubstitutesIDS"
		return model.Product{}, ErrTypeFormat(campo, "11, 14, 3", product.SubstitutesIDS)
	}

	for _, subtance := range subtances {
		sub := model.Substance{}
		sub.ID = subtance
		productDb.Substance = append(productDb.Substance, &sub)
	}

	for _, warehouse := range warehouses {
		ware := model.Warehouse{}
		ware.ID = warehouse
		productDb.Warehouses = append(productDb.Warehouses, &ware)

	}

	for _, subID := range substitues {
		subTemp := model.Product{}
		subTemp.ID = subID
		productDb.Substitutes = append(productDb.Substitutes, &subTemp)
	}

	return productDb, nil
}

func HandleNewElement(body []byte) (model.Product, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Product{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleNewElements(body []byte) (newElement, updateElement []model.Product, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return []model.Product{}, []model.Product{}, controller.ErrFormatIncorrect(err, example, body)
	}

	newElement = make([]model.Product, 0, len(adapts))
	updateElement = make([]model.Product, 0, len(adapts))

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

func HandleAdapt(body []byte) (updateElement []*model.Product, err error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		return nil, err
	}
	updateElement = make([]*model.Product, 0, len(adapts))

	for _, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		updateElement = append(updateElement, &element)
	}
	return updateElement, nil
}
