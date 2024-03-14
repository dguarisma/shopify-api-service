package model

type Substance struct {
	CustomModel
	ProductID *uint // para borrar
	Name      string
	Status    bool
	Products  []*Product `gorm:"many2many:substance_products;"`
}
