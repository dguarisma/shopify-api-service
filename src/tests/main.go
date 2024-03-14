package main

import (
	"desarrollosmoyan/lambda/src/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No se cargaron la variables de entorno")
	}
}

func main() {
	dbUri := os.Getenv("DB_URI") // test

	db, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	CreateDatabase(db)
	// ResetDatabase(db)
}

func ResetDatabase(db *gorm.DB) {

	type SubstanceProducts struct {
		ProductId   uint
		SubstanceId uint
	}

	type WarehouseProducts struct {
		ProductId   uint
		WarehouseId uint
	}
	if err := db.Model(&model.Product{}).
		Where("1=1").
		Update("MakerID", nil).Error; err != nil {
		log.Println(err.Error())
		return
	}

	if err := db.Model(&model.Product{}).
		Where("1=1").
		Update("TrademarkID", nil).Error; err != nil {
		log.Println(err.Error())
		return
	}

	var tables []interface{} = []interface{}{
		&SubstanceProducts{},
		&WarehouseProducts{},
		&model.Product{},
		&model.Trademark{},
		&model.Maker{},
		&model.TypesProduct{},
		&model.ReceptionArt{},
		&model.Article{},
		&model.Purchase{},
		&model.Substance{},
		&model.Warehouse{},
		&model.Supplier{},
		&model.Pack{},
		&model.CategoryThree{},
		&model.CategoryTwo{},
		&model.CategoryOne{},
	}

	for _, table := range tables {
		if err := db.Unscoped().
			Where("1=1").
			Delete(table).Error; err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func CreateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(
		/* &model.Pack{},
		&model.Supplier{},
		&model.Warehouse{},
		&model.Purchase{},
		&model.Maker{},
		&model.Trademark{},
		&model.CategoryOne{},
		&model.CategoryTwo{},
		&model.CategoryThree{},
		&model.TypesProduct{},
		&model.Article{},
		&model.ReceptionArt{},
		&model.Substance{}, */
		&model.Product{},
	); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Base de datos cargada con exito !!!")
}
