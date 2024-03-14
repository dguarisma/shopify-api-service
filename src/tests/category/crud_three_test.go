package category_test

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

func TestCRUDThree(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}

	type DeleteById struct{ ID uint }

	tx := handleTest.Begin()
	defer handleTest.Rollback()
	handleTest.Show()
	base := "/categorythree"

	t.Run("Category three", func(t *testing.T) {
		path := ""
		endpoint := base + path

		categoryOne := model.CategoryOne{}
		categoryOneChange := model.CategoryOne{}
		categoryTwo := model.CategoryTwo{}
		categoryTwoChange := model.CategoryTwo{}

		tx.Save(&categoryOne)
		tx.Save(&categoryOneChange)
		tx.Save(&categoryTwo)
		tx.Save(&categoryTwoChange)

		example := model.CategoryThree{
			Name:          "example",
			Status:        true,
			CategoryOneID: categoryOne.ID,
			CategoryTwoID: categoryTwo.ID,
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.Marshal(example)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}

			method := http.MethodPost

			request := events.APIGatewayProxyRequest{
				Path:       path,
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			handleTest.ShowRequest(
				t,
				endpoint,
				method,
				body,
				resBody,
			)

			res := model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf("Error to parse response: %s", err.Error())
			}
			example.CustomModel = res.CustomModel

			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("Update", func(t *testing.T) {
			example = model.CategoryThree{
				CustomModel:   example.CustomModel,
				Name:          "example-change",
				Status:        false,
				CategoryOneID: categoryOneChange.ID,
				CategoryTwoID: categoryTwoChange.ID,
			}
			body, err := json.Marshal(example)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}

			method := http.MethodPut

			request := events.APIGatewayProxyRequest{
				Path:       path,
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			handleTest.ShowRequest(
				t,
				endpoint,
				method,
				body,
				resBody,
			)

			res := model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf("Error to parse response: %s", err.Error())
			}

			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("GetAll", func(t *testing.T) {

			method := http.MethodGet
			request := events.APIGatewayProxyRequest{
				Path:       path,
				HTTPMethod: method,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			res := []model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}

			handleTest.ShowGetRequest(
				t,
				endpoint,
				resBody,
			)

			if len(res) != 1 {
				t.Errorf("Error the categories should be 1")
			}

			handleTest.DifferentMap(
				t,
				example,
				res[0],
			)
		})

		t.Run("GetById", func(t *testing.T) {
			queryParams := map[string]string{
				"ID": fmt.Sprint(example.ID),
			}
			request := events.APIGatewayProxyRequest{
				Path:                  path,
				HTTPMethod:            http.MethodGet,
				QueryStringParameters: queryParams,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			handleTest.ShowGetByRequest(
				t,
				endpoint,
				resBody,
				queryParams,
			)

			res := model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}
			handleTest.DifferentMap(
				t,
				example,
				res,
			)
		})

		t.Run("GetByName", func(t *testing.T) {

			queryParams := map[string]string{
				"Name": fmt.Sprint(example.Name),
			}

			request := events.APIGatewayProxyRequest{
				Path:                  path,
				HTTPMethod:            http.MethodGet,
				QueryStringParameters: queryParams,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			handleTest.ShowGetByRequest(
				t,
				endpoint,
				resBody,
				queryParams,
			)

			res := []model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}
			if len(res) != 1 {
				t.Errorf("Error the categories should be 1")
			}

			handleTest.DifferentMap(
				t,
				example,
				res[0],
			)
		})

		t.Run("Delete", func(t *testing.T) {
			deleteById := DeleteById{example.ID}
			body, err := json.Marshal(deleteById)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}
			method := http.MethodDelete
			request := events.APIGatewayProxyRequest{
				Path:       path,
				HTTPMethod: method,
				Body:       string(body),
			}

			deleteCategory := model.CategoryThree{}

			resBody := handleTest.UseHandleRequest(t,
				server.CategoryThree,
				request,
				http.StatusOK,
			)

			handleTest.ShowRequest(
				t,
				endpoint,
				method,
				body,
				resBody,
			)

			res := model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}

			handleTest.DifferentMap(
				t,
				deleteCategory,
				res,
			)
		})
	})
}
