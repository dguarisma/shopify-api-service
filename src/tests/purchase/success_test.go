package purchases

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"testing"

	"github.com/aws/aws-lambda-go/events"
	"gorm.io/gorm"
)

func TestCrud(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	deletes := []interface{}{
		&model.Product{},
		&model.Article{},
		&model.Receiving{},
		&model.Purchase{},
		&model.Supplier{},
		&model.Warehouse{},
	}
	defer handleTest.DeleteInfo(t, deletes)
	/*
	 */
	// para ver como son las request
	handleTest.Show()
	//handleTest.Show()

	t.Run("Simple CRUD (one)", func(t *testing.T) {
		supplier := model.Supplier{}
		warehouse := model.Warehouse{}
		handleTest.Db.Save(&supplier)
		handleTest.Db.Save(&warehouse)

		example := model.Purchase{
			Articles:             []model.Article{},
			SupplierID:           supplier.ID,
			WarehouseID:          warehouse.ID,
			Notes:                "example",
			Discount:             1,
			DiscountEarliyPay:    2,
			SubtotalWithDiscount: 3,
			SubTotal:             4,
			Total:                5,
			Tax:                  1,
			Status:               0,
			DiscountGlobal:       40,
			ReceptionStatus:      1,
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.MarshalIndent(example, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Purchase,
				request,
				http.StatusOK,
			)

			res := model.Purchase{}
			if err := json.
				Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}

			example.ID = res.ID
			example.CreatedAt = res.CreatedAt
			example.DateExpireInvoice = res.DateExpireInvoice

			handleTest.ShowRequest(t,
				http.MethodPost,
				string(body),
				resBody,
			)

			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("Update", func(t *testing.T) {
			ti := time.Now()
			example.DateExpireInvoice = &ti
			body, err := json.MarshalIndent(example, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPut,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Purchase,
				request,
				http.StatusOK,
			)

			res := model.Purchase{}
			if err := json.
				Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}

			handleTest.ShowRequest(t,
				http.MethodPut,
				string(body),
				resBody,
			)

			res.CreatedAt = example.CreatedAt
			if res.DateExpireInvoice.Equal(*example.DateExpireInvoice) {
				example.DateExpireInvoice = res.DateExpireInvoice
			}
			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("GetById", func(t *testing.T) {
			path := map[string]string{
				"ID": fmt.Sprint(example.ID),
			}

			expect := model.PurchaseForGet{
				Purchase: example,
			}

			err := handleTest.Db.Transaction(func(tx *gorm.DB) error {
				warehouse := model.Warehouse{}
				supplier := model.Supplier{}

				if err := tx.
					Find(&warehouse, example.WarehouseID).
					Error; err != nil {
					return err
				}

				if err := tx.
					Find(&supplier, example.SupplierID).
					Error; err != nil {
					return err
				}
				supplier.CustomModel = model.CustomModel{ID: supplier.ID}
				warehouse.CustomModel = model.CustomModel{ID: warehouse.ID}
				expect.Supplier = supplier
				expect.Warehouse = warehouse

				return nil
			})
			if err != nil {
				t.Fatal(err.Error())
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				QueryStringParameters: path,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Purchase,
				request,
				http.StatusOK,
			)

			res := model.PurchaseForGet{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}

			*res.Purchase.DateExpireInvoice = res.Purchase.DateExpireInvoice.UTC().Round(time.Second)
			*expect.Purchase.DateExpireInvoice = expect.Purchase.DateExpireInvoice.UTC().Round(time.Second)
			if res.Purchase.CreatedAt.Equal(expect.CreatedAt) {
				expect.CreatedAt = res.CreatedAt
			}

			expect.Purchase.CustomModel = res.Purchase.CustomModel

			handleTest.DifferentMap(
				t,
				expect,
				res,
			)
		})

		t.Run("DeleteById", func(t *testing.T) {
			path := map[string]uint{
				"ID": example.ID,
			}

			body, err := json.MarshalIndent(path, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Purchase,
				request,
				http.StatusOK,
			)

			res := model.Purchase{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}
			handleTest.ShowRequest(t, http.MethodDelete, string(body), resBody)

			handleTest.DifferentMap(
				t,
				model.Purchase{},
				res,
			)
		})

	})

	t.Run("Complex CRUD (one)", func(t *testing.T) {

		supplier := model.Supplier{}
		warehouse := model.Warehouse{}

		products := []model.Product{
			{
				Name: "example1",
				Sku:  "example1",
				Ean:  "example1",
			},
		}

		handleTest.Db.Save(&supplier)
		handleTest.Db.Save(&warehouse)
		handleTest.Db.Save(&products)

		example := model.Purchase{
			//Articles:    make([]model.Article, len(products)),
			Articles:    make([]model.Article, len(products)),
			SupplierID:  supplier.ID,
			WarehouseID: warehouse.ID,
		}
		for i, product := range products {
			example.Articles[i] = model.Article{
				ProductID:          product.ID,
				BasePrice:          1,
				Count:              2,
				Tax:                3,
				Discount:           4,
				DiscountAdditional: 5,
				Bonus:              6,
				SubTotal:           7,
				Total:              8,
			}
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.MarshalIndent(example, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Purchase,
				request,
				http.StatusOK,
			)

			res := model.Purchase{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}

			example.ID = res.ID
			example.CreatedAt = res.CreatedAt
			//example.Reception = res.Reception

			example.ID = res.ID
			example.CreatedAt = res.CreatedAt
			example.DateExpireInvoice = res.DateExpireInvoice

			if j, _ := json.MarshalIndent(example, "", "\t"); true {
				fmt.Printf("\n\n%v\n\n", string(j))
			}

			if j, _ := json.MarshalIndent(res, "", "\t"); true {
				fmt.Printf("\n\n%v\n\n", string(j))
			}

			/*
				for i, resArticle := range res.Articles {
					example.Articles[i].ID = resArticle.ID
					example.Articles[i].PurchaseID = resArticle.PurchaseID
				}
			*/

			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		/*
			t.Run("Update", func(t *testing.T) {
				t.Run("supplier and warehouse", func(t *testing.T) {
					supplier := model.Supplier{}
					warehouse := model.Warehouse{}
					handleTest.Db.Save(&supplier)
					handleTest.Db.Save(&warehouse)

					example.SupplierID = supplier.ID
					example.WarehouseID = warehouse.ID

					body, err := json.Marshal(example)
					if err != nil {
						t.Errorf("Mal formato de json: %s", err.Error())
						return
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodPut,
						Body:       string(body),
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res := model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					example.ID = res.ID
					example.CreatedAt = res.CreatedAt
					//	example.Reception = res.Reception

					for i, resArticle := range res.Articles {
						example.Articles[i].ID = resArticle.ID
						example.Articles[i].PurchaseID = resArticle.PurchaseID
					}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)

				})

				t.Run("_article", func(t *testing.T) {

					product := model.Product{
						Name: "example10",
						Sku:  "example10",
						Ean:  "example10",
					}
					handleTest.Db.Save(&product)
					artTemp := example.Articles[0]

					example.Articles[0] = model.Article{
						CustomModel:        model.CustomModel{ID: artTemp.ID},
						PurchaseID:         artTemp.PurchaseID,
						ProductID:          product.ID,
						Count:              artTemp.Count + 1,
						BasePrice:          artTemp.BasePrice + 1,
						Tax:                artTemp.Tax + 1,
						Discount:           artTemp.Discount + 1,
						DiscountAdditional: artTemp.DiscountAdditional + 1,
						Bonus:              artTemp.Bonus + 1,
						SubTotal:           artTemp.SubTotal + 1,
						Total:              artTemp.Total + 1,
						//ReceptionInfo:      artTemp.ReceptionInfo,
					}

					body, err := json.Marshal(example)
					if err != nil {
						t.Errorf("Mal formato de json: %s", err.Error())
						return
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodPut,
						Body:       string(body),
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res := model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)

				})
				t.Run("add_article", func(t *testing.T) {
					product := model.Product{
						Name: "example11",
						Sku:  "example11",
						Ean:  "example11",
					}

					handleTest.Db.Save(&product)

					artTemp := example.Articles[0]

					example.Articles = append(example.Articles, model.Article{
						PurchaseID:         artTemp.PurchaseID,
						ProductID:          product.ID,
						Count:              artTemp.Count + 1,
						BasePrice:          artTemp.BasePrice + 1,
						Tax:                artTemp.Tax + 1,
						Discount:           artTemp.Discount + 1,
						DiscountAdditional: artTemp.DiscountAdditional + 1,
						Bonus:              artTemp.Bonus + 1,
						SubTotal:           artTemp.SubTotal + 1,
						Total:              artTemp.Total + 1,
						//ReceptionInfo:      artTemp.ReceptionInfo,
					})

					body, err := json.Marshal(example)
					if err != nil {
						t.Errorf("Mal formato de json: %s", err.Error())
						return
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodPut,
						Body:       string(body),
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res := model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					for i, resArticle := range res.Articles {
						example.Articles[i].ID = resArticle.ID
						example.Articles[i].PurchaseID = resArticle.PurchaseID
					}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)
					// check update
					path := map[string]string{
						"ID": fmt.Sprint(example.ID),
					}

					request = events.APIGatewayProxyRequest{
						HTTPMethod:            http.MethodGet,
						QueryStringParameters: path,
					}

					resBody = handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res = model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					example.ID = res.ID
					example.CreatedAt = res.CreatedAt
					//	example.Reception = res.Reception
						if example.CreatedAt.Equal(*res.CreatedAt) {
							// da la misma fecha solo que cambia el formato
							example.CreatedAt = res.CreatedAt
						}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)

				})

				t.Run("delete_article", func(t *testing.T) {
					example.Articles = example.Articles[1:]
					body, err := json.Marshal(example)
					if err != nil {
						t.Errorf("Mal formato de json: %s", err.Error())
						return
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodPut,
						Body:       string(body),
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res := model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					for i, resArticle := range res.Articles {
						example.Articles[i].ID = resArticle.ID
						example.Articles[i].PurchaseID = resArticle.PurchaseID
					}

					example.CreatedAt = res.CreatedAt
					//example.Reception = res.Reception

					handleTest.DifferentMap(
						t,
						example,
						res,
					)
					// check update
					path := map[string]string{
						"ID": fmt.Sprint(example.ID),
					}

					request = events.APIGatewayProxyRequest{
						HTTPMethod:            http.MethodGet,
						QueryStringParameters: path,
					}

					resBody = handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res = model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					example.ID = res.ID
					example.CreatedAt = res.CreatedAt
					//example.Reception = res.Reception
						if example.CreatedAt.Equal(*res.CreatedAt) {
							example.CreatedAt = res.CreatedAt
						}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)

				})

				t.Run("add_update_article", func(t *testing.T) {
					product := model.Product{
						Name: "example18",
						Sku:  "example18",
						Ean:  "example18",
					}
					handleTest.Db.Save(&product)
					artTemp := example.Articles[0]

					example.Articles[0] = model.Article{
						CustomModel:        model.CustomModel{ID: artTemp.ID},
						PurchaseID:         artTemp.PurchaseID,
						ProductID:          product.ID,
						Count:              artTemp.Count + 1,
						BasePrice:          artTemp.BasePrice + 1,
						Tax:                artTemp.Tax + 1,
						Discount:           artTemp.Discount + 1,
						DiscountAdditional: artTemp.DiscountAdditional + 1,
						Bonus:              artTemp.Bonus + 1,
						SubTotal:           artTemp.SubTotal + 1,
						Total:              artTemp.Total + 1,
						ReceptionInfo:      artTemp.ReceptionInfo,
					}
					example.Articles = append(example.Articles, model.Article{
						ProductID: product.ID,
					})

					body, err := json.Marshal(example)
					if err != nil {
						t.Errorf("Mal formato de json: %s", err.Error())
						return
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodPut,
						Body:       string(body),
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res := model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					for i, resArticle := range res.Articles {
						example.Articles[i].ID = resArticle.ID
						example.Articles[i].PurchaseID = resArticle.PurchaseID
					}

					example.CreatedAt = res.CreatedAt
					//example.Reception = res.Reception

					handleTest.DifferentMap(
						t,
						example,
						res,
					)
					// check update
					path := map[string]string{
						"ID": fmt.Sprint(example.ID),
					}

					request = events.APIGatewayProxyRequest{
						HTTPMethod:            http.MethodGet,
						QueryStringParameters: path,
					}

					resBody = handleTest.UseHandleRequest(t,
						server.Purchase,
						request,
						http.StatusOK,
					)

					res = model.Purchase{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					example.ID = res.ID
					example.CreatedAt = res.CreatedAt
					// example.Reception = res.Reception
						if example.CreatedAt.Equal(*res.CreatedAt) {
							example.CreatedAt = res.CreatedAt
						}

					handleTest.DifferentMap(
						t,
						example,
						res,
					)

				})


			})

		*/
	})
}
