package purchases_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGets(t *testing.T) {
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
	t.Run("Get", func(t *testing.T) {

		supplier := model.Supplier{
			BusinessName: "BusinessName",
			DaysPayment:  "DaysPayment",
			Nit:          "Nit",
			PaymenTerm:   "PaymenTerm",
			Cupo:         1,
			Discount:     1,
			LeadTimeBaq:  1,
			LeadTimeBog:  1,
			NameContact:  "NameContact",
			EmailContact: "EmailContact",
			PhoneContact: "PhoneContact",
			Status:       true,
			Location:     "Location",
		}
		warehouse := model.Warehouse{
			Name:       "Name",
			Department: "Department",
			City:       "City",
			Location:   "Location",
			Status:     true,
		}
		handleTest.Db.Save(&supplier)
		handleTest.Db.Save(&warehouse)

		example := model.Purchase{
			Articles:    []model.Article{},
			SupplierID:  supplier.ID,
			WarehouseID: warehouse.ID,
		}

		t.Run("Set", func(t *testing.T) {
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
			example.Reception = res.Reception

			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("ById", func(t *testing.T) {
			expect := model.PurchaseForGet{
				Purchase: example,
			}
			expect.Warehouse = warehouse
			expect.Supplier = supplier

			path := map[string]string{
				"ID": fmt.Sprint(expect.ID),
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

			if expect.CreatedAt.Equal(*res.CreatedAt) {
				// da la misma fecha solo que cambia el formato
				expect.CreatedAt = res.CreatedAt
			}

			res.Warehouse.CreatedAt = expect.Warehouse.CreatedAt
			res.Warehouse.UpdatedAt = expect.Warehouse.UpdatedAt
			res.Supplier.CreatedAt = expect.Supplier.CreatedAt
			res.Supplier.UpdatedAt = expect.Supplier.UpdatedAt

			handleTest.DifferentMap(
				t,
				expect,
				res,
			)
		})

	})
}
