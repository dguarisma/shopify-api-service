package model

type Article struct {
	CustomModel
	PurchaseID         uint
	ProductID          uint
	Count              uint
	BasePrice          float32 // tomo para el calculo
	Tax                float32
	Discount           float32
	DiscountAdditional float32
	Bonus              uint
	SubTotal           float32
	Total              float32
	ReceptionInfo      []ReceptionArt
}

func (art Article) GetPurchaseId() string  { return "purchase_id" }
func (art Article) GetProductId() string   { return "product_id" }
func (art Article) GetWarehouseId() string { return "warehouse_id" }
