package receptions

// go test ./* -v

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/shopifyserv"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"gorm.io/gorm"
)

func SetDatabase(tx *gorm.DB) (*gorm.DB, error) {
	prod := model.Product{
		HandlesBog: "Bogo",
		HandlesBaq: "Barran",
	}
	suppliers := model.Supplier{}
	warehouse := model.Warehouse{
		City: "Barranquilla",
	}
	prod.SetId(1)
	suppliers.SetId(2)
	warehouse.SetId(3)

	if err := tx.Save(&warehouse).Error; err != nil {
		return tx, err
	}
	if err := tx.Save(&suppliers).Error; err != nil {
		return tx, err
	}
	if err := tx.Save(&prod).Error; err != nil {
		return tx, err
	}

	pur := model.Purchase{
		WarehouseID: warehouse.ID,
		SupplierID:  suppliers.ID,
	}
	pur.SetId(4)

	if err := tx.Save(&pur).Error; err != nil {
		return tx, err
	}

	article := model.Article{
		ProductID:  prod.ID,
		PurchaseID: pur.ID,
	}
	article.SetId(5)

	if err := tx.Save(&article).Error; err != nil {
		return tx, err
	}

	reception1 := model.ReceptionArt{
		ArticleID: article.ID,
	}

	if err := tx.Save(&reception1).Error; err != nil {
		return tx, err
	}
	return tx, nil
}

func TestExample(t *testing.T) {
	handlerTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tx, err := SetDatabase(handlerTest.Db.Begin())
	if err != nil {
		tx.Rollback()
		log.Fatalf("error montando la base de datos:%v", err.Error())
	}
	defer tx.Rollback()

	//shopiServ := shopifyserv.New("2022-10", shopifyserv.GetCredentials())
	handlerSend := &HandlerPurchaseStatus{
		db:          tx,
		shopifyserv: &ShopifyMockService{},
	}

	article := model.Article{}
	if err := tx.First(&article).Error; err != nil {
		log.Fatalf("error montando la base de datos:%v", err.Error())
	}

	article.Count = 100

	if err := tx.Save(&article).Error; err != nil {
		log.Fatalf(err.Error())
	}

	preReception := model.ReceptionArt{
		Count:     10,
		ArticleID: article.ID,
	}

	if err := tx.Save(&preReception).Error; err != nil {
		log.Fatalf(err.Error())
	}

	t.Run("Insert One", func(t *testing.T) {
		reception := model.ReceptionArt{
			Count:     article.Count - (preReception.Count + 60), // count 30
			ArticleID: article.ID,
		}

		reception2, err := handlerSend.InsertOne(tx, &reception)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if reception2.ID == 0 {
			t.Error("Deberia retornar una recepcion con id")
		}

		receptionDb := model.ReceptionArt{}
		if err := handlerTest.Db.Find(&receptionDb, reception2.ID).Error; err != nil {
			t.Error("Deberia existir la reception en la base de datos")
		}
	})

	t.Run("IsArticleOverflow", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := handlerSend.IsArticleOverflow(tx, article.ID)
			if err != nil {
				t.Error("Deberia retornar una nil, porque no se supero la cantidad de articulo de la compra")
			}
		})

		t.Run("Error for Overflow", func(t *testing.T) {
			reception := model.ReceptionArt{
				Count:     article.Count, // 100
				ArticleID: article.ID,
			}

			_, err := handlerSend.InsertOne(tx, &reception)
			if err != nil {
				log.Fatalf(err.Error())
			}

			err = handlerSend.IsArticleOverflow(tx, article.ID)
			if err == nil {
				t.Error("Deberia retornar un error, porque se supero la cantidad de articulo de la compra")
			}
		})
	})

	t.Run("getCityAndProductShopifyId", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			shopifyUpdate, err := handlerSend.getCityAndProductShopifyId(tx, article.ID)
			if err == nil {
				t.Error("Deberia retornar un error, porque se supero la cantidad de articulo de la compra")
			}
			if j, _ := json.MarshalIndent(shopifyUpdate, "", "\t"); true {
				fmt.Printf("\n\n%v\n\n", string(j))
			}
		})

	})

	/*
		t.Run("Update Element", func(t *testing.T) {
			reception := model.ReceptionArt{}
			if err := tx.
				First(&reception, model.ReceptionArt{ArticleID: article.ID}).
				Error; err != nil {
				t.Errorf(err.Error())
			}

			reception.Count++

			status, err := handlerSend.thereCheckChangeInCount(tx, &reception)
			if err != nil {
				t.Error(err.Error())
			}

			if status != true {
				t.Error("el estado deberia ser true, se cambio la cantidad de un elemento")
			}
		})

		t.Run("Update Element no Count", func(t *testing.T) {
			preReception := model.ReceptionArt{}
			if err := tx.First(&preReception).Error; err != nil {
				t.Errorf(err.Error())
			}

			preReception.Batch = "example"

			status, err := handlerSend.thereCheckChangeInCount(tx, &preReception)
			if err != nil {
				t.Error(err.Error())
			}

			if status != false {
				t.Error("el estado deberia ser true, se cambio la cantidad de un elemento")
			}
		})

	*/
	/*
		t.Run("isArticleComplete", func(t *testing.T) {

			t.Run("With some items", func(t *testing.T) {
				status, err := handlerSend.isArticleComplete(tx, article.ID, article.Count)
				if err != nil {
					t.Error(err.Error())
				}
				if status == true {
					t.Error("el estado deberia ser falso porque el articulo no esta completado")
				}
			})

			t.Run("Should be full", func(t *testing.T) {
				rececion := model.ReceptionArt{
					ArticleID: article.ID,
					Count:     article.Count - preReception.Count,
				}
				if err := tx.Save(&rececion).Error; err != nil {
					t.Error(err.Error())
				}

				status, err := handlerSend.isArticleComplete(tx, article.ID, article.Count)
				if err != nil {
					t.Error(err.Error())
				}
				if status == false {
					t.Error("el estado deberia ser verdadero porque el articulo esta completado")
				}
			})
		})

		t.Run("getByid", func(t *testing.T) {

			warehouse := &model.Warehouse{}
			if err := tx.First(&warehouse).Error; err != nil {
				t.Errorf(err.Error())
			}

			example, err := getById[model.Warehouse](tx, warehouse.ID, "warehouse")
			if err != nil {
				t.Errorf(err.Error())
			}
			handlerTest.DifferentMap(t, warehouse, example)

			if j, _ := json.MarshalIndent(example, "", "\t"); true {
				fmt.Printf("\n\n%v\n\n", string(j))
			}

		})

		t.Run("CheckAndUpdateShopify", func(t *testing.T) {
			mock := &ShopifyMockService{}
			handlerSend := &HandlerPurchaseStatus{
				db:          tx,
				shopifyserv: mock,
			}

			purchase := &model.Purchase{}
			if err := tx.First(&purchase).Error; err != nil {
				t.Errorf(err.Error())
			}

			prod := &model.Product{}
			if err := tx.First(&prod).Error; err != nil {
				t.Errorf(err.Error())
			}

			article := model.Article{
				PurchaseID: purchase.ID,
				ProductID:  prod.ID,
				Count:      100,
			}
			if err := tx.Save(&article).Error; err != nil {
				t.Errorf(err.Error())
			}
			reception := model.ReceptionArt{
				ArticleID: article.ID,
				Count:     100,
			}

			if err := tx.Save(&reception).Error; err != nil {
				t.Errorf(err.Error())
			}
			if _, err := handlerSend.
				CheckAndUpdateShopify(tx, &reception); err != nil {
				t.Errorf(err.Error())
			}
			data := mock.getItem()
			if data.Items[0].Available != uint64(article.Count) {
				t.Error("La cantidad de items debe ser igual a la de articulos")
			}

		})

	*/
}

type ShopifyMockService struct {
	item shopifyserv.ProductsUpdates
}

func (sms *ShopifyMockService) SendUpdateInventory(item shopifyserv.ProductsUpdates) error {
	sms.item = item
	fmt.Println(item)
	return nil
}

func (sms *ShopifyMockService) getItem() shopifyserv.ProductsUpdates {
	return sms.item
}
