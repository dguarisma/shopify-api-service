package server

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/controller/categories"
	"desarrollosmoyan/lambda/src/controller/inventory"
	"desarrollosmoyan/lambda/src/controller/makers"
	"desarrollosmoyan/lambda/src/controller/packs"
	"desarrollosmoyan/lambda/src/controller/products"
	"desarrollosmoyan/lambda/src/controller/purchases"
	"desarrollosmoyan/lambda/src/controller/receptions"
	"desarrollosmoyan/lambda/src/controller/substances"
	"desarrollosmoyan/lambda/src/controller/suppliers"
	"desarrollosmoyan/lambda/src/controller/trademarks"
	"desarrollosmoyan/lambda/src/controller/typesproducts"
	"desarrollosmoyan/lambda/src/controller/warehouses"
	"desarrollosmoyan/lambda/src/mailserv"
	"desarrollosmoyan/lambda/src/shopifyserv"
	"os"
	"time"

	"gorm.io/gorm"
)

const (
	// methods
	Maker uint8 = iota
	Trademark
	Pack
	Typesproduct
	Substance
	Supplier
	Warehouse
	Purchase
	Product
	Reception
	Inventory
	CategoryOne
	CategoryTwo
	CategoryThree
)

func handleRepository(db *gorm.DB, method uint8, path string) (repo controller.CRUD) {
	switch method {
	case Purchase:

		config := &mailserv.MailerConfig{
			Host:     os.Getenv("HOST"),
			Username: os.Getenv("AWS_SES_ACCESS_KEY"),
			Password: os.Getenv("AWS_SES_SECRET_KEY"),
			Sender:   os.Getenv("EMAIL"),
			Port:     587,
			Timeout:  5 * time.Second,
		}
		mailservice := mailserv.New(db, config)
		repo = purchases.NewPurchaseRepository(db, mailservice)

	case Reception:
		shopiServ := shopifyserv.New("2022-10", shopifyserv.GetCredentials(
			os.Getenv("BARRANQULLA"), os.Getenv("BOGOTA"),
		))
		repo = receptions.NewReceptionServ(db, shopiServ)
	case Maker:
		repo = makers.NewMakerRepository(db)
	case Inventory:
		repo = inventory.NewInventoryRepository(db)
	case Trademark:
		repo = trademarks.NewTradeMarkRepository(db)
	case Pack:
		repo = packs.NewPackRepository(db)
	case Typesproduct:
		repo = typesproducts.NewTypesProductRepository(db)
	case Substance:
		repo = substances.NewSubstanceRepository(db)
	case Supplier:
		repo = suppliers.NewSupplierRepository(db)
	case Warehouse:
		repo = warehouses.NewWarehouseRepository(db)
	case Product:
		repo = products.NewProductRepository(db)
	case CategoryOne:
		repo = categories.NewCategoryOneRepository(db)
	case CategoryTwo:
		repo = categories.NewCategoryTwoRepository(db)
	case CategoryThree:
		repo = categories.NewCategoryThreeRepository(db)
	default:
		panic("method don't set")
	}
	return repo
}
