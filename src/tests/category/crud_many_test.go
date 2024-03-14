package category_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestCrudOne(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	deletes := []interface{}{
		&model.CategoryOne{},
		&model.CategoryTwo{},
		&model.CategoryThree{},
	}

	handleTest.DeleteInfo(t, deletes)
	defer handleTest.DeleteInfo(t, deletes)

	t.Run("Category one", func(t *testing.T) {

		categories := []model.CategoryOne{
			{
				Name:   "example1",
				Status: true,
			},
			{
				Name:   "example2",
				Status: false,
			},
			{
				Name:   "example3",
				Status: true,
			},
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.Marshal(categories)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}

			request := events.APIGatewayProxyRequest{
				Path:       "/one",
				HTTPMethod: http.MethodPost,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
				request,
				http.StatusOK,
			)

			res := []model.CategoryOne{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf("Error to parse response: %s", err.Error())
			}

			for i, curCategory := range res {
				categories[i].CustomModel = curCategory.CustomModel
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})

		t.Run("Update", func(t *testing.T) {
			t.Run("all", func(t *testing.T) {
				for i, curCategory := range categories {
					categories[i] = model.CategoryOne{
						CustomModel: curCategory.CustomModel,
						Name:        curCategory.Name + "-change",
						Status:      !curCategory.Status,
					}
				}

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/one",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}
				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryOne{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("add_one", func(t *testing.T) {
				categories = append(categories, model.CategoryOne{
					Name:   "new-category",
					Status: false,
				})
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/one",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryOne{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
			t.Run("delete_one", func(t *testing.T) {
				categories = categories[1:]
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/one",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryOne{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("delete_update_add", func(t *testing.T) {
				categories = categories[1:]

				cateTemp := categories[0]
				categories[0] = model.CategoryOne{
					CustomModel: cateTemp.CustomModel,
					Name:        cateTemp.Name + "-other-change",
					Status:      !cateTemp.Status,
				}

				categories = append(categories, model.CategoryOne{
					Name:   "new-category",
					Status: false,
				})

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/one",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryOne{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
		})

		t.Run("GetAll", func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Path:       "/one",
				HTTPMethod: http.MethodGet,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
				request,
				http.StatusOK,
			)

			res := []model.CategoryOne{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}

			if len(res) != len(categories) {
				t.Errorf("Error the categories should be %v", len(categories))
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})
	})

	t.Run("Category two", func(t *testing.T) {

		categoryOne := model.CategoryOne{}
		categoryOneSecond := model.CategoryOne{}
		handleTest.Db.Save(&categoryOne)
		handleTest.Db.Save(&categoryOneSecond)

		categories := []model.CategoryTwo{
			{
				Name:          "example1",
				Status:        true,
				CategoryOneID: categoryOne.ID,
			},
			{
				Name:          "example2",
				Status:        false,
				CategoryOneID: categoryOne.ID,
			},
			{
				Name:          "example3",
				Status:        true,
				CategoryOneID: categoryOne.ID,
			},
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.Marshal(categories)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}

			request := events.APIGatewayProxyRequest{
				Path:       "/two",
				HTTPMethod: http.MethodPost,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
				request,
				http.StatusOK,
			)

			res := []model.CategoryTwo{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf("Error to parse response: %s", err.Error())
			}

			for i, curCategory := range res {
				categories[i].CustomModel = curCategory.CustomModel
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})

		t.Run("Update", func(t *testing.T) {
			t.Run("all", func(t *testing.T) {
				for i, curCategory := range categories {
					categories[i] = model.CategoryTwo{
						CustomModel:   curCategory.CustomModel,
						Name:          curCategory.Name + "-change",
						Status:        !curCategory.Status,
						CategoryOneID: categoryOne.ID,
					}
				}
				categories[0].CategoryOneID = categoryOneSecond.ID

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/two",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}
				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryTwo{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("add_one", func(t *testing.T) {
				categories = append(categories, model.CategoryTwo{
					Name:          "new-category",
					Status:        false,
					CategoryOneID: categoryOne.ID,
				})
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/two",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryTwo{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
			t.Run("delete_one", func(t *testing.T) {
				categories = categories[1:]
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/two",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryTwo{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("delete_update_add", func(t *testing.T) {
				categories = categories[1:]

				cateTemp := categories[0]
				categories[0] = model.CategoryTwo{
					CustomModel:   cateTemp.CustomModel,
					Name:          cateTemp.Name + "-other-change",
					Status:        !cateTemp.Status,
					CategoryOneID: cateTemp.CategoryOneID,
				}

				categories = append(categories, model.CategoryTwo{
					Name:          "new-category",
					Status:        false,
					CategoryOneID: cateTemp.CategoryOneID,
				})

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/two",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryTwo{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
		})

		t.Run("GetAll", func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Path:       "/two",
				HTTPMethod: http.MethodGet,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
				request,
				http.StatusOK,
			)

			res := []model.CategoryTwo{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf(
					"Error to parse response: %s",
					err.Error(),
				)
			}

			if len(res) != len(categories) {
				t.Errorf("Error the categories should be %v", len(categories))
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})
	})

	t.Run("Category three", func(t *testing.T) {

		categoryOne := model.CategoryOne{}
		categoryOneSecond := model.CategoryOne{}

		categoryTwo := model.CategoryTwo{}
		categoryTwoSecond := model.CategoryThree{}

		handleTest.Db.Save(&categoryOne)
		handleTest.Db.Save(&categoryOneSecond)

		handleTest.Db.Save(&categoryTwo)
		handleTest.Db.Save(&categoryTwoSecond)

		categories := []model.CategoryThree{
			{
				Name:          "example1",
				Status:        true,
				CategoryOneID: categoryOne.ID,
				CategoryTwoID: categoryTwo.ID,
			},
			{
				Name:          "example2",
				Status:        false,
				CategoryOneID: categoryOne.ID,
				CategoryTwoID: categoryTwo.ID,
			},
			{
				Name:          "example3",
				Status:        true,
				CategoryOneID: categoryOne.ID,
				CategoryTwoID: categoryTwo.ID,
			},
		}

		t.Run("Insert", func(t *testing.T) {
			body, err := json.Marshal(categories)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
			}

			request := events.APIGatewayProxyRequest{
				Path:       "/three",
				HTTPMethod: http.MethodPost,
				Body:       string(body),
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
				request,
				http.StatusOK,
			)

			res := []model.CategoryThree{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Errorf("Error to parse response: %s", err.Error())
			}

			for i, curCategory := range res {
				categories[i].CustomModel = curCategory.CustomModel
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})

		t.Run("Update", func(t *testing.T) {
			t.Run("all", func(t *testing.T) {
				for i, curCategory := range categories {
					categories[i] = model.CategoryThree{
						CustomModel:   curCategory.CustomModel,
						Name:          curCategory.Name + "-change",
						Status:        !curCategory.Status,
						CategoryOneID: curCategory.CategoryOneID,
						CategoryTwoID: curCategory.CategoryTwoID,
					}
				}

				categories[0].CategoryOneID = categoryOneSecond.ID

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/three",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}
				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryThree{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("add_one", func(t *testing.T) {
				categories = append(categories, model.CategoryThree{
					Name:          "new-category",
					Status:        false,
					CategoryOneID: categoryOne.ID,
					CategoryTwoID: categoryTwo.ID,
				})
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/three",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryThree{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
			t.Run("delete_one", func(t *testing.T) {
				categories = categories[1:]
				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/three",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				res := []model.CategoryThree{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})

			t.Run("delete_update_add", func(t *testing.T) {
				categories = categories[1:]

				cateTemp := categories[0]
				categories[0] = model.CategoryThree{
					CustomModel:   cateTemp.CustomModel,
					Name:          cateTemp.Name + "-other-change",
					Status:        !cateTemp.Status,
					CategoryOneID: cateTemp.CategoryOneID,
					CategoryTwoID: cateTemp.CategoryTwoID,
				}

				categories = append(categories, model.CategoryThree{
					Name:          "new-category",
					Status:        false,
					CategoryOneID: cateTemp.CategoryOneID,
				})

				body, err := json.Marshal(categories)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
				}

				request := events.APIGatewayProxyRequest{
					Path:       "/three",
					HTTPMethod: http.MethodPost,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Category,
					request,
					http.StatusOK,
				)

				lastCategory := len(categories) - 1
				categories[lastCategory].ID = categories[lastCategory-1].ID + 1

				res := []model.CategoryThree{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Errorf("Error to parse response: %s", err.Error())
				}

				handleTest.DifferentMap(
					t,
					categories,
					res,
				)
			})
		})

		t.Run("GetAll", func(t *testing.T) {
			request := events.APIGatewayProxyRequest{
				Path:       "/three",
				HTTPMethod: http.MethodGet,
			}

			resBody := handleTest.UseHandleRequest(t,
				server.Category,
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

			if len(res) != len(categories) {
				t.Errorf("Error the categories should be %v", len(categories))
			}

			handleTest.DifferentMap(
				t,
				categories,
				res,
			)
		})
	})
}
