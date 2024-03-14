package main

import (
	"desarrollosmoyan/lambda/src/controller/substances"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestAll(t *testing.T) {
	defer utils.DeleteInfo(&model.Pack{})
	first := substances.Adapt{
		Name:   "Example1",
		Status: true,
	}
	/*
		manyWithoutUpdate := []substances.Adapt{
			{
				Name:   "Example",
				Status: false,
			},
			{
				Name:   "Example",
				Status: false,
			},
		}

		manyWithUpdate := []substances.Adapt{}

		newElement := substances.Adapt{
			Name:   "New-Item",
			Status: true,
		}
	*/

	t.Run("Insert one", func(t *testing.T) {
		method := http.MethodPost
		body, err := json.MarshalIndent(first, "", "  ")
		if err != nil {
			t.Errorf("Mal formato de json: %s", err.Error())
			return
		}

		request := events.APIGatewayProxyRequest{
			HTTPMethod: method,
			Body:       string(body),
		}

		resBody := utils.UseHandleRequest(t, request, http.StatusOK)

		response := substances.Adapt{}
		if err := json.Unmarshal([]byte(resBody), &response); err != nil {
			t.Errorf("Mal retorno de body: %s", err.Error())
		}

		if response.ID == 0 {
			t.Errorf("El retorno del ID no puede ser cero")
		}

		first.ID = response.ID

		utils.ShowRequest(t, method, string(body), resBody)

		utils.DifferentMap(t,
			first,
			response,
		)
	})

	/*
		t.Run("Insert many without update", func(t *testing.T) {
			method := http.MethodPost
			body, err := json.MarshalIndent(manyWithoutUpdate, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := []substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
			}
			for i, element := range response {
				if element.ID == 0 {
					t.Errorf("El retorno del ID no puede ser cero")
					return
				}
				manyWithoutUpdate[i].ID = element.ID
			}

			utils.ShowRequest(t, method, string(body), resBody)
			utils.DifferentMap(t,
				manyWithoutUpdate,
				response,
			)
		})

		t.Run("Insert many with update", func(t *testing.T) {
			method := http.MethodPost

			for i, element := range manyWithoutUpdate {
				element.Name += "change"
				element.Status = !element.Status
				manyWithoutUpdate[i] = element
			}

			manyWithUpdate = append(manyWithoutUpdate, newElement)

			body, err := json.MarshalIndent(manyWithUpdate, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := []substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}

			lastIndex := len(manyWithUpdate) - 1

			for i := 0; i < lastIndex-1; i++ {
				element := manyWithUpdate[i]
				elementRes := response[i]
				if element.ID != elementRes.ID {
					t.Errorf("Id expected: %d | got: %d", element.ID, elementRes.ID)
					return
				}

			}
			manyWithUpdate[lastIndex].ID = response[lastIndex].ID

			utils.ShowRequest(t, method, string(body), resBody)
			utils.DifferentMap(t,
				manyWithUpdate,
				response,
			)
		})

		t.Run("Update first Item", func(t *testing.T) {
			method := http.MethodPut
			first.Name += "second-change"
			first.Status = !first.Status

			body, err := json.MarshalIndent(first, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}

			utils.ShowRequest(t,
				method,
				string(body),
				resBody,
			)

			utils.DifferentMap(t,
				first,
				response,
			)
		})

		t.Run("Get all", func(t *testing.T) {
			method := http.MethodGet
			queryParams := map[string]string{}

			request := events.APIGatewayProxyRequest{
				HTTPMethod:            method,
				QueryStringParameters: queryParams,
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := []substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}
			utils.ShowRequestEmpty(t, method, resBody)

			all := []substances.Adapt{first}
			all = append(all, manyWithUpdate...)
			utils.DifferentMap(t, all, response)
		})

		t.Run("Get by name", func(t *testing.T) {
			method := http.MethodGet
			queryParams := map[string]string{
				"Name": fmt.Sprint(manyWithoutUpdate[0].Name),
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod:            method,
				QueryStringParameters: queryParams,
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := []substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}
			utils.ShowRequestBy(t, queryParams, "Name", method, resBody)

			manyWithoutLast := manyWithUpdate[:len(manyWithUpdate)-1]
			utils.DifferentMap(t, manyWithoutLast, response)
		})

		t.Run("Get by ID", func(t *testing.T) {
			method := http.MethodGet
			queryParams := map[string]string{
				"ID": fmt.Sprint(first.ID),
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod:            method,
				QueryStringParameters: queryParams,
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			response := substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &response); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}
			utils.ShowRequestBy(t, queryParams, "ID", method, resBody)
			utils.DifferentMap(t, first, response)
		})

		t.Run("Delete first Item", func(t *testing.T) {
			method := http.MethodDelete
			element := substances.Adapt{
				ID: first.ID,
			}

			body, err := json.MarshalIndent(element, "", "  ")
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: method,
				Body:       string(body),
			}

			resBody := utils.UseHandleRequest(t, request, http.StatusOK)

			responseAdap := substances.Adapt{}
			if err := json.Unmarshal([]byte(resBody), &responseAdap); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}

			utils.ShowRequest(t,
				method,
				string(body),
				resBody,
			)

			utils.DifferentMap(t,
				substances.Adapt{},
				responseAdap,
			)

			request = events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				QueryStringParameters: map[string]string{
					"ID": fmt.Sprint(element.ID),
				},
			}

			resBody = utils.UseHandleRequest(t, request, http.StatusNotFound)

			errMsg := response.ErrMsg{}
			if err := json.Unmarshal([]byte(resBody), &errMsg); err != nil {
				t.Errorf("Mal retorno de body: %s", err.Error())
				return
			}

			if errMsg.Error != "Elemento no encontrado" {
				t.Errorf("El elemento no deberia existir")
				return
			}
		})
	*/
}
