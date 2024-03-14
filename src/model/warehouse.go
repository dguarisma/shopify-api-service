package model

type Warehouse struct {
	CustomModel
	Name       string
	Department string
	City       string
	Location   string
	Status     bool
	Products   []*Product `gorm:"many2many:warehouse_products;"`
	Purchase   []Purchase
}
