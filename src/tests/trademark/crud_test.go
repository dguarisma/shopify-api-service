package trademark_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const endpoint = "/trademark"

var (
	handleTest *utils.HandleTest
)

func init() {
	handleTest, _ = utils.NewHandleTest()
	handleTest.Show()
}

func TestGet(t *testing.T) {
	tx := handleTest.Begin()
	defer handleTest.Rollback()

	markTemp := model.Maker{}

	err := tx.Save(&markTemp).Error
	assert.NoError(t, err, "it shouldn't return error when insert in database")

	items := make([]model.Trademark, 100)
	for i, item := range items {
		item.MakerID = markTemp.ID
		if i%2 == 0 {
			item.Status = true
			item.Name = fmt.Sprintf("par-%v", i)
		} else {
			item.Name = fmt.Sprintf("inpar-%v", i)
		}
		if i > 50 && i < 60 {
			item.Name = fmt.Sprint("example")
		}
		items[i] = item
	}

	err = tx.Save(&items).Error
	assert.NoError(t, err, "it shouldn't return error when insert in database")

	t.Run("All", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Get,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)
		response := []model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetRequest(
			t, endpoint, resBody,
		)

		for i, item := range response {
			assert.Equal(t, item.ID, items[i].ID)
			assert.Equal(t, item.MakerID, items[i].MakerID)
			assert.Equal(t, item.Name, items[i].Name)
			assert.Equal(t, item.Status, items[i].Status)
		}
	})
	t.Run("ById", func(t *testing.T) {
		queryParams := map[string]string{"ID": fmt.Sprint(items[10].ID)}

		request := events.APIGatewayProxyRequest{
			HTTPMethod:            handleTest.Get,
			QueryStringParameters: queryParams,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)

		response := model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetByRequest(
			t, endpoint, resBody, queryParams,
		)

		assert.Equal(t, response.ID, items[10].ID)
		assert.Equal(t, response.Name, items[10].Name)
		assert.Equal(t, response.Status, items[10].Status)
	})

	t.Run("Name", func(t *testing.T) {
		queryParams := map[string]string{"Name": fmt.Sprint(items[20].Name)}

		names := []model.Trademark{}
		for _, item := range items {
			if item.Name == queryParams["Name"] {
				names = append(names, item)
			}
		}

		request := events.APIGatewayProxyRequest{
			HTTPMethod:            handleTest.Get,
			QueryStringParameters: queryParams,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)
		response := []model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetByRequest(
			t, endpoint, resBody, queryParams,
		)

		for i, item := range response {
			assert.Equal(t, item.ID, names[i].ID)
			assert.Equal(t, item.Name, names[i].Name)
			assert.Equal(t, item.Status, names[i].Status)
		}
	})
}

func TestInsert(t *testing.T) {
	t.Run("one minimum", func(t *testing.T) {
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		mark := model.Maker{}
		err := tx.Save(&mark).Error
		assert.NoError(t, err)
		item := model.Trademark{MakerID: mark.ID}

		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)

		response := model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.NotEqual(t,
			uint(0),
			response.ID,
			"it shouldn't be 0 a ID",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		itemDb := model.Trademark{}

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, response.ID, itemDb.ID)
		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.MakerID, itemDb.MakerID)
	})
	t.Run("one full", func(t *testing.T) {
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		mark := model.Maker{}
		err := tx.Save(&mark).Error
		assert.NoError(t, err)

		item := model.Trademark{
			MakerID: mark.ID,
			Name:    "example",
			Status:  true,
		}

		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)

		response := model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.NotEqual(t, uint(0), response.ID)
		assert.Equal(t, item.Name, response.Name)
		assert.Equal(t, item.MakerID, response.MakerID)
		assert.Equal(t, item.Status, response.Status)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		itemDb := model.Trademark{}

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, 0, result.RowsAffected)

		assert.Equal(t, response.ID, itemDb.ID)
		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.Status, itemDb.Status)
		assert.Equal(t, response.MakerID, itemDb.MakerID)
	})

	t.Run("three elements", func(t *testing.T) {
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		marks := []model.Maker{{}, {}}
		err := tx.Save(&marks).Error
		assert.NoError(t, err)

		items := []model.Trademark{
			{Name: "example", Status: true, MakerID: marks[0].ID},
			{Name: "example2", Status: false, MakerID: marks[0].ID},
			{Name: "example3", Status: true, MakerID: marks[1].ID},
		}

		body, err := json.Marshal(items)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)

		response := []model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")
		assert.Equal(t,
			len(items), len(response),
			"it should be same count of items and response",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		for i, itemRs := range response {
			preItem := items[i]
			assert.NotEqual(t,
				uint(0),
				itemRs.ID,
				"it shouldn't be 0 a ID",
			)
			assert.Equal(t, preItem.Name, itemRs.Name)
			assert.Equal(t, preItem.Status, itemRs.Status)
			assert.Equal(t, preItem.MakerID, itemRs.MakerID)
		}

		itemsDb := []model.Trademark{}
		result := tx.Find(&itemsDb)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, int64(0), result.RowsAffected)

		for i, itemRs := range response {
			itemDb := itemsDb[i]
			assert.NotEqual(t,
				uint(0),
				itemRs.ID,
				"it shouldn't be 0 a ID",
			)
			assert.Equal(t, itemDb.Name, itemRs.Name)
			assert.Equal(t, itemDb.Status, itemRs.Status)
			assert.Equal(t, itemDb.MakerID, itemRs.MakerID)
		}
	})

	t.Run("Inserts news and updates without delete items", func(t *testing.T) {
		// este test es para poder insertar y actualizar varios items sin borrar
		// los que ya estaban en la base de datos
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		marks := model.Maker{}
		err := tx.Save(&marks).Error

		items := []model.Trademark{
			{Name: "example1", Status: true, MakerID: marks.ID},
			{Name: "example2", Status: false, MakerID: marks.ID},
			{Name: "example3", Status: true, MakerID: marks.ID},
		}

		err = tx.Save(&items).Error
		assert.NoError(t, err, "it shouldn't return error when insert in database")

		items2 := append(items[1:], []model.Trademark{
			{Name: "example-2", Status: true, MakerID: marks.ID},
			{Name: "example-2", Status: false, MakerID: marks.ID},
		}...)

		body, err := json.Marshal(items2)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Post,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)

		response := []model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")
		assert.Equal(t,
			len(items2), len(response),
			"it should be same count of items and response",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Post,
			body,
			resBody,
		)

		for _, item := range response {
			assert.NotEqual(t,
				uint(0),
				item.ID,
				"it shouldn't be 0 a ID",
			)
		}

		itemsDb := []model.Trademark{}
		result := tx.Find(&itemsDb)

		assert.NoError(t, result.Error)
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, len(items2)+1, len(itemsDb))
	})
}

func TestUpdate(t *testing.T) {
	tx := handleTest.Begin()
	defer handleTest.Rollback()

	mark := model.Maker{}
	err := tx.Save(&mark).Error
	assert.NoError(t, err, "save Error")
	itemDb := model.Trademark{MakerID: mark.ID}
	err = tx.Save(&itemDb).Error
	assert.NoError(t, err, "save Error")

	t.Run("one", func(t *testing.T) {
		item := model.Trademark{Name: "example-update", Status: true, MakerID: mark.ID}
		item.ID = itemDb.ID
		body, err := json.Marshal(item)
		assert.NoError(t, err, "json has a problem")

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Put,
			Body:       string(body),
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Trademark,
			request,
			http.StatusOK,
		)
		response := model.Trademark{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		assert.NotEqual(t,
			uint(0),
			response.ID,
			"it shouldn't be 0 a ID",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Put,
			body,
			resBody,
		)

		result := tx.Find(&itemDb, response.ID)
		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, response.ID, itemDb.ID)
	})
}
