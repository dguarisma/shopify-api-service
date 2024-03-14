package purchases

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"encoding/json"
	"time"
)

type ArticleAdapt struct {
	ID                 uint
	PurchaseID         uint
	ProductID          uint
	Count              uint
	BasePrice          float32
	Tax                float32
	Discount           float32
	DiscountAdditional float32
	Bonus              uint
	SubTotal           float32
	Total              float32
}

type Adapt struct {
	ID                   uint
	SupplierID           uint
	WarehouseID          uint
	Notes                string
	Discount             float32
	DiscountEarliyPay    float32
	SubTotal             float32
	DiscountGlobal       float32
	SubtotalWithDiscount float32
	Total                float32
	Status               uint8
	Tax                  float32
	Articles             []ArticleAdapt
	CreatedAt            time.Time

	ReceptionNotes    string
	InvoiceNumber     string // numero de factura
	ReceptionStatus   uint8  // [1, 2, 3] parcialmente entregado, cancelado, completado //
	DateExpireInvoice *time.Time
}

func checkFields(adapt Adapt) (element model.Purchase, err error) {
	if adapt.SupplierID == 0 {
		campo := "SupplierID"
		return element, controller.ErrTypeFormat(campo, "10", adapt.SupplierID)
	}
	if adapt.WarehouseID == 0 {
		campo := "WarehouseID"
		return element, controller.ErrTypeFormat(campo, "10", adapt.WarehouseID)
	}
	for _, art := range adapt.Articles {
		if art.ProductID == 0 {
			campo := "Article"
			return element, controller.ErrTypeFormat(campo, "10", adapt.Articles)
		}
	}

	element = model.Purchase{
		CustomModel:          model.CustomModel{ID: adapt.ID},
		SupplierID:           adapt.SupplierID,
		WarehouseID:          adapt.WarehouseID,
		Notes:                adapt.Notes,
		Discount:             adapt.Discount,
		DiscountEarliyPay:    adapt.DiscountEarliyPay,
		SubTotal:             adapt.SubTotal,
		SubtotalWithDiscount: adapt.SubtotalWithDiscount,
		DiscountGlobal:       adapt.DiscountGlobal,
		ReceptionNotes:       adapt.ReceptionNotes,
		Status:               adapt.Status,
		Total:                adapt.Total,
		Tax:                  adapt.Tax,
		Articles:             make([]model.Article, len(adapt.Articles)),
		InvoiceNumber:        adapt.InvoiceNumber,
		ReceptionStatus:      adapt.ReceptionStatus,
		DateExpireInvoice:    adapt.DateExpireInvoice,
	}

	for i, art := range adapt.Articles {
		element.Articles[i] = model.Article{
			CustomModel:        model.CustomModel{ID: art.ID},
			PurchaseID:         art.PurchaseID,
			ProductID:          art.ProductID,
			Count:              art.Count,
			BasePrice:          art.BasePrice,
			Tax:                art.Tax,
			Discount:           art.Discount,
			DiscountAdditional: art.DiscountAdditional,
			Bonus:              art.Bonus,
			SubTotal:           art.SubTotal,
			Total:              art.Total,
		}
	}
	return element, nil
}

func HandleNewElement(body []byte) (model.Purchase, error) {
	adapt := Adapt{}
	if err := json.Unmarshal(body, &adapt); err != nil {
		example, _ := json.MarshalIndent(adapt, "", "  ")
		return model.Purchase{}, controller.ErrFormatIncorrect(err, example, body)
	}
	return checkFields(adapt)
}

func HandleNewElements(body []byte) (*[]model.Purchase, error) {
	adapts := []Adapt{}
	if err := json.Unmarshal(body, &adapts); err != nil {
		example, _ := json.MarshalIndent(adapts, "", "  ")
		return nil, controller.ErrFormatIncorrect(err, example, body)
	}

	elements := make([]model.Purchase, len(adapts))
	for i, adapt := range adapts {
		element, err := checkFields(adapt)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}
	return &elements, nil
}
