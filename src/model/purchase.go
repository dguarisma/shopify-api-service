package model

import (
	"time"
)

type Purchase struct {
	CustomModel
	SupplierID           uint
	WarehouseID          uint
	Notes                string
	Discount             float32
	DiscountEarliyPay    float32
	DiscountGlobal       float32
	SubtotalWithDiscount float32
	SubTotal             float32
	Total                float32
	Tax                  float32
	Status               uint8  // 0 = new | 1 = send | 2 = cancelled
	InvoiceNumber        string // numero de factura
	ReceptionStatus      uint8  // 1 = parcialmente entregado |  2 = completado  | 3 = cancelado
	ReceptionNotes       string
	DateExpireInvoice    *time.Time
	CreatedAt            time.Time `gorm:"<-:create"`
	Articles             []Article
}

func (p *Purchase) GetWarehouseId() string { return "warehouse_id" }

// para tener mas detallada la compra
type PurchaseForGet struct {
	Purchase
	Products  []Product
	Warehouse Warehouse
	Supplier  Supplier
	UpdatedAt time.Time `json:"updateAt"`
}
