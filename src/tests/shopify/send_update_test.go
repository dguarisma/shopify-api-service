package shopify_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {

	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")
	tx := handleTest.Begin()
	defer tx.Rollback()

	idProductExampleBarranquilla := "6776151638039"

	ware := model.Warehouse{City: "Barranquilla"}
	supplier := model.Supplier{}

	err = tx.Save(&supplier).Error
	assert.NoError(t, err)
	err = tx.Save(&ware).Error
	assert.NoError(t, err)

	products := []model.Product{
		{
			Name:       "example1",
			Sku:        time.Now().String(),
			Ean:        time.Now().String(),
			HandlesBaq: idProductExampleBarranquilla,
			Warehouses: []*model.Warehouse{
				{CustomModel: model.CustomModel{ID: ware.ID}},
			},
		},
	}
	err = tx.Save(&products).Error
	assert.NoError(t, err)

	purchase := model.Purchase{
		SupplierID:  supplier.ID,
		WarehouseID: ware.ID,
	}
	err = tx.Save(&purchase).Error
	assert.NoError(t, err)

	article := model.Article{
		PurchaseID: purchase.ID,
		ProductID:  products[0].ID,
		Count:      4,
	}

	err = tx.Save(&article).Error
	assert.NoError(t, err)

	example := model.ReceptionArt{
		ArticleID: article.ID,
		Count:     article.Count,
	}

	t.Run("UpdateShopify", func(t *testing.T) {

		body, err := json.Marshal(example)
		if err != nil {
			t.Errorf("Mal formato de json: %s", err.Error())
			return
		}

		request := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequest(t,
			server.Reception,
			request,
			http.StatusOK,
		)

		res := model.ReceptionArt{}
		if err := json.Unmarshal([]byte(resBody), &res); err != nil {
			t.Error(err.Error())
		}

		example.CustomModel = res.CustomModel
		handleTest.DifferentMap(t,
			example,
			res,
		)

		purchaseUpdate := model.Purchase{}
		tx.Find(&purchaseUpdate, purchase.ID)

		// deberia ser 2 porque la compra esta completa
		if purchaseUpdate.ReceptionStatus != 2 {
			t.Errorf("deberia ser status 2 porque la compra estaria completa")
		}

	})
}
