package model

type Product struct {
	CustomModel

	// Datos Basicos
	Name     string
	Sku      string //`gorm:"unique"`
	Ean      string //`gorm:"unique"`
	Taxed    bool
	Iva      float32 `json:"iva"`
	UrlImage string  // Icono del producto

	// Datos medidas
	PackID    *uint
	Pack      Pack `gorm:"foreignKey:PackID;references:ID"`
	PackInfo  string
	Quantity  uint64
	MakerUnit string `gorm:"size:100"`
	Weight    float32
	Height    float32
	Width     float32
	PackUnit  uint64
	Depth     uint64

	// Datos Extra
	Keywords       string
	Variation      string
	Wrapper        string
	WrapperUnit    string
	Status         bool
	HandlesBog     string `gorm:"size:255"`
	HandlesBaq     string `gorm:"size:255"`
	IDFloorProduct uint

	// Datos adicionales
	MakerID         *uint
	TrademarkID     *uint
	CategoryOneID   *uint
	CategoryTwoID   *uint
	CategoryThreeID *uint

	Maker         Maker
	CategoryOne   CategoryOne   `gorm:"foreignKey:CategoryOneID;references:ID"`
	CategoryTwo   CategoryTwo   `gorm:"foreignKey:CategoryTwoID;references:ID"`
	CategoryThree CategoryThree `gorm:"foreignKey:CategoryThreeID;references:ID"`

	Substitutes    []*Product   `gorm:"many2many:product_substitues"`
	Substance      []*Substance `gorm:"many2many:substance_products;"`
	Warehouses     []*Warehouse `gorm:"many2many:warehouse_products;"`
	TypesProductID *uint        // seteo un int en vez de un uint
}
