package inventory_test

import (
	"desarrollosmoyan/lambda/src/controller/inventory"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetInventory(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	tx := handleTest.Begin()
	defer handleTest.Rollback()

	SetDatabase(tx)
	article := model.Article{}
	if err := tx.
		Preload("ReceptionInfo").
		First(&article).Error; err != nil {
		t.Fatalf("err: %v", err)
	}

	purchase := model.Purchase{CustomModel: model.CustomModel{ID: article.PurchaseID}}
	if err := tx.Find(&purchase, model.Purchase{
		CustomModel: model.CustomModel{ID: article.PurchaseID}},
	).Error; err != nil {
		t.Fatalf("err: %v", err)
	}

	t.Run("GetInventory", func(t *testing.T) {
		endpoint := "/inventory"
		queryParams := map[string]string{
			"WarehouseID": fmt.Sprint(purchase.WarehouseID),
			"From":        fmt.Sprint(time.Now().AddDate(0, 0, -5).Format("2006/01/02")),
			"To":          fmt.Sprint(time.Now().AddDate(0, 0, -1).Format("2006/01/02")),
			"Sku":         "CJJCLRYX",
		}
		fmt.Println(queryParams)
		request := events.APIGatewayProxyRequest{
			Path:                  endpoint,
			HTTPMethod:            http.MethodGet,
			QueryStringParameters: queryParams,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Inventory,
			request,
			http.StatusOK,
		)

		fmt.Println(resBody)
		res := []inventory.InventorySell{}
		err := json.Unmarshal([]byte(resBody), &res)
		assert.NoError(t, err)
		showMeTable(res)

		handleTest.ShowGetByRequest(
			t,
			endpoint,
			resBody,
			queryParams,
		)
	})
	/*
		example := inventory.HandlerTables{Db: tx}
		t.Run("GetPurchaseInventory", func(t *testing.T) {
			inventory := inventory.InventoryInfo{
				ProductId:   article.ProductID,
				WarehouseID: purchase.WarehouseID,
			}

			purchaseInventory, err := example.GetPurchaseInventory(inventory)
			if err != nil {
				assert.NoError(t, err, "the return shouldn't be a error ")
			}
			assert.Equal(t,
				len(*purchaseInventory),
				len(article.ReceptionInfo),
				"it should has same count of items",
			)

			for i, reception := range article.ReceptionInfo {
				purcha := (*purchaseInventory)[i]
				assert.Equal(t,
					int64(reception.Count),
					purcha.Purchase.Units,
					"it should be same count",
				)
			}
		})

		t.Run("SortTables", func(t *testing.T) {
			t.Run("error for empty tables", func(t *testing.T) {
				purchaseList := []inventory.InventorySell{}
				sellList := []inventory.InventorySell{}
				fullList, err := example.SortTables(sellList, purchaseList)

				assert.NotNil(t,
					err,
					"it should return an error that say:`No hay compras o ventas registradas` not nil",
				)
				if fullList != nil {
					t.Errorf("it should be a nil because doesn't have table to sort")
				}
			})

			t.Run("one row for only one item", func(t *testing.T) {
				purchaseList := []inventory.InventorySell{
					{
						Date:        example.SetDate(time.Now()),
						ProductId:   "example",
						ProductName: "example name",
						Sku:         "example sku",
					},
				}
				sellList := []inventory.InventorySell{}
				fullList, err := example.SortTables(sellList, purchaseList)

				if err != nil {
					t.Errorf("it should not return an error: %v", err)
				}
				if len(*fullList) != 1 {
					t.Errorf("it should return an one element")
				}
			})

			t.Run("one row for purchase and sell itself day", func(t *testing.T) {
				purchaseList := []inventory.InventorySell{
					{
						Date:        example.SetDate(time.Now()),
						ProductId:   "example",
						ProductName: "example name",
						Sku:         "example sku",
						Purchase: inventory.Store{
							Units: 10,
							Cost:  1000,
						},
					},
				}
				sellList := []inventory.InventorySell{
					{
						Date:        example.SetDate(time.Now()),
						ProductId:   "example",
						ProductName: "example name",
						Sku:         "example sku",
						Selled: inventory.Store{
							Units: 9,
							Cost:  1100,
						},
					},
				}
				fullList := inventory.SortTables(purchaseList, sellList)

				if err != nil {
					t.Errorf("it should not return an error: %v", err)
				}
				if len(fullList) != 1 {
					t.Errorf("it should return an one element")
				}

				fullRow := fullList[0]
				sellRow := sellList[0]
				purchaselRow := purchaseList[0]

				assert.Equal(t,
					sellRow.Date, fullRow.Date,
					"the date should be a same",
				)

				assert.Equal(t,
					sellRow.ProductId, fullRow.ProductId,
					"the product id should be a same",
				)

				assert.Equal(t,
					sellRow.ProductName, fullRow.ProductName,
					"the name should be a same",
				)

				assert.Equal(t,
					sellRow.Sku, fullRow.Sku,
					"the sku should be a same",
				)

				assert.Equal(t,
					sellRow.Selled, fullRow.Selled,
					"the selled should be a same",
				)

				assert.Equal(t,
					purchaselRow.Purchase, fullRow.Purchase,
					"the purchaseshould be a same",
				)
			})
			t.Run("aaa", func(t *testing.T) {

				purchaseList := []inventory.InventorySell{
					{
						Date:        example.SetDate(time.Date(2022, 02, 28, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Purchase:    inventory.Store{Units: 100, Cost: 2000},
					}, {
						Date:        example.SetDate(time.Date(2022, 03, 01, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Purchase:    inventory.Store{Units: 200, Cost: 2100},
					}, {
						Date:        example.SetDate(time.Date(2022, 03, 02, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Purchase:    inventory.Store{Units: 300, Cost: 2000},
					}, {
						Date:        example.SetDate(time.Date(2022, 03, 03, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Purchase:    inventory.Store{Units: 400, Cost: 2050},
					},
				}

				if j, _ := json.Marshal(purchaseList); true {
					fmt.Printf("\n\n%v\n\n", string(j))
				}
				sellList := []inventory.InventorySell{
					{
						Date:        example.SetDate(time.Date(2022, 03, 01, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Selled:      inventory.Store{Units: 9, Cost: 1100},
					}, {
						Date:        example.SetDate(time.Date(2022, 03, 01, 0, 0, 0, 0, time.UTC)),
						ProductId:   "6639852191767",
						ProductName: "ACETAMINOFEN",
						Sku:         "CJHSBEPL",
						Selled:      inventory.Store{Units: 9, Cost: 1100},
					},
				}
				fullList := inventory.SortTables(purchaseList, sellList)

				if err != nil {
					t.Errorf("it should not return an error: %v", err)
				}
				if len(fullList) != 1 {
					t.Errorf("it should return an one element")
				}

				fullRow := fullList[0]
				sellRow := sellList[0]
				purchaselRow := purchaseList[0]

				assert.Equal(t,
					sellRow.Date, fullRow.Date,
					"the date should be a same",
				)

				assert.Equal(t,
					sellRow.ProductId, fullRow.ProductId,
					"the product id should be a same",
				)

				assert.Equal(t,
					sellRow.ProductName, fullRow.ProductName,
					"the name should be a same",
				)

				assert.Equal(t,
					sellRow.Sku, fullRow.Sku,
					"the sku should be a same",
				)

				assert.Equal(t,
					sellRow.Selled, fullRow.Selled,
					"the selled should be a same",
				)

				assert.Equal(t,
					purchaselRow.Purchase, fullRow.Purchase,
					"the purchaseshould be a same",
				)
			})
		})

	*/
}

func showMeTable(superLista []inventory.InventorySell) {
	tabla := [][]string{
		{
			"fecha",
			"productId",
			"name",
			"sku",
			"init unit",
			"init cost",
			"purcha unit",
			"purcha cost",
			"available unit",
			"available cost",
			"selled unit",
			"selled cost",
			"end unit",
			"end cost",
		},
	}

	for _, data := range superLista {
		tabla = append(tabla, []string{
			fmt.Sprint(data.Date),
			fmt.Sprint(data.ProductId),
			fmt.Sprint(data.ProductName),
			fmt.Sprint(data.Sku),
			fmt.Sprint(data.InitInventory.Units),
			fmt.Sprintf("%.2f", data.InitInventory.Cost),
			fmt.Sprint(data.Purchase.Units),
			fmt.Sprintf("%.2f", data.Purchase.Cost),
			fmt.Sprint(data.Available.Units),
			fmt.Sprintf("%.2f", data.Available.Cost),
			fmt.Sprint(data.Selled.Units),
			fmt.Sprintf("%.2f", data.Selled.Cost),
			fmt.Sprint(data.EndInventory.Units),
			fmt.Sprintf("%.2f", data.EndInventory.Cost),
		})
	}

	log := make([]uint, len(tabla[0]))

	for _, row := range tabla {
		for i, item := range row {
			if log[i] < uint(len(item)) {
				log[i] = uint(len(item))
			}
		}
	}
	for _, row := range tabla {
		fmt.Print("|")

		for i, v := range row {
			space := log[i] - uint(len(v)) + 1
			fmt.Printf(" %v", v)

			for j := 0; j < int(space); j++ {
				fmt.Print(" ")
			}
			fmt.Print("|")
		}
		fmt.Print("\n")
	}
}

func SetDatabase(tx *gorm.DB) (*gorm.DB, error) {
	prod := model.Product{
		Name:       "example",
		Sku:        "CJJCLRYX",
		HandlesBog: "idForBogota",
		HandlesBaq: "idForBarranquilla",
		Warehouses: []*model.Warehouse{{City: "Bogotá"}},
	}

	if err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&prod).Error; err != nil {
		return tx, err
	}

	suppliers := model.Supplier{}
	if err := tx.Save(&suppliers).Error; err != nil {
		return tx, err
	}

	pur := model.Purchase{
		WarehouseID: prod.Warehouses[0].ID,
		SupplierID:  suppliers.ID,
		Articles: []model.Article{
			{
				Count:     100,
				BasePrice: 100000,
				ProductID: prod.ID,
				ReceptionInfo: []model.ReceptionArt{
					{Count: 45, Date: time.Now().AddDate(0, 0, -5)},
					{Count: 45, Date: time.Now().AddDate(0, 0, -105)},
					{Count: 55, Date: time.Now().AddDate(0, 0, -2)},
				},
			},
		},
	}

	if err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&pur).Error; err != nil {
		return tx, err
	}

	return tx, nil
}
