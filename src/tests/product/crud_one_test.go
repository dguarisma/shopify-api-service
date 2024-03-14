package product_test

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
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func SetDatabase(tx *gorm.DB) {

}

func TestSuccess(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")
	product := []model.Product{}

	if err := handleTest.Db.Find(&product).Error; err != nil {
		log.Println(err.Error())
	}

	if j, _ := json.MarshalIndent(product, "", "\t"); true {
		fmt.Printf("\n\n%v\n\n", string(j))
	}

	tx := handleTest.Begin()
	defer handleTest.Rollback()
	endpoint := "/product"

	// handleTest.Show()

	t.Run("Insert one minimum", func(t *testing.T) {
		item := model.Product{}
		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.NotEqual(t,
			0,
			response,
			"it shouldn't be 0 a ID",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		itemDb := model.Product{}

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, response.ID, itemDb.ID)
	})

	t.Run("Insert one full", func(t *testing.T) {
		item := model.Product{
			Name:   "example",
			Status: true,
		}
		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.NotEqual(t, 0, response.ID)
		assert.Equal(t, item.Name, response.Name)
		assert.Equal(t, item.Status, response.Status)

		assert.NotEqual(t,
			0,
			response,
			"it shouldn't be 0 a ID",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		itemDb := model.Product{}

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, 0, result.RowsAffected)

		assert.Equal(t, response.ID, itemDb.ID)
		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.Status, itemDb.Status)
	})

	t.Run("Update one", func(t *testing.T) {
		item := model.Product{}

		err := tx.First(&item).Error
		assert.NoError(t, err)

		item = model.Product{
			CustomModel: model.CustomModel{ID: item.ID},
			Name:        "example",
			Status:      true,
		}

		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Put,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.Equal(t, item.ID, response.ID)
		assert.Equal(t, item.Name, response.Name)
		assert.Equal(t, item.Status, response.Status)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Put,
			body,
			resBody,
		)

		itemDb := model.Product{}

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, int64(0), result.RowsAffected)

		assert.Equal(t, response.ID, itemDb.ID)
		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.Status, itemDb.Status)
	})

	t.Run("GetById", func(t *testing.T) {
		item := model.Product{
			Name:   "example",
			Status: true,
		}

		err := tx.Save(&item).Error
		assert.NoError(t, err)

		queryParams := map[string]string{
			"ID": fmt.Sprint(item.ID),
		}

		request := events.APIGatewayProxyRequest{
			HTTPMethod:            handleTest.Get,
			QueryStringParameters: queryParams,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err)

		assert.Equal(t, item.ID, response.ID)
		assert.Equal(t, item.Name, response.Name)
		assert.Equal(t, item.Status, response.Status)
	})

	t.Run("Delete one", func(t *testing.T) {
		deleteById := utils.DeleteById{}

		err := tx.Model(&model.Product{}).
			Select("id").Limit(1).
			Scan(&deleteById).Error

		assert.NoError(t, err)

		body, err := json.Marshal(deleteById)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Delete,
			Body:       string(body),
		}
		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), response.ID)
		assert.Equal(t, "", response.Name)
		assert.Equal(t, false, response.Status)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Delete,
			body,
			resBody,
		)

		itemDb := model.Product{}

		result := tx.Find(&itemDb, deleteById.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.Equal(t, int64(0), result.RowsAffected)
	})
}
