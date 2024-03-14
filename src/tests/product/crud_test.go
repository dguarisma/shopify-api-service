package product_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed product.json
	pro []byte

	//go:embed products.json
	prods []byte
)

const endpoint = "/product"

func TestInsert(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")

	request := events.APIGatewayProxyRequest{
		HTTPMethod: handleTest.Get,
		MultiValueHeaders: map[string][]string{
			"x-amz-security-token": {
				"DQYAAHNGJOMNHYYBXRGWCYLAURP",
			},
		},
	}

	type Pag struct {
		Limit      int             `json:"limit,omitempty;query:limit"`
		Page       int             `json:"page,string,omitempty;query:page"`
		Sort       string          `json:"sort,omitempty;query:sort"`
		TotalRows  int64           `json:"totalRows"`
		TotalPages int             `json:"totalPages"`
		Rows       []model.Product `json:"Rows"`
	}

	resBody := handleTest.UseHandleRequestTx(t,
		server.Product,
		request,
		http.StatusOK,
	)

	response := Pag{}
	err = json.Unmarshal([]byte(resBody), &response)
	assert.NoError(t, err, "json has a problem")

	handleTest.ShowGetRequest(
		t, endpoint, resBody,
	)
	fmt.Println(resBody)
	//expected := []model.Product{}
	//count := 0

	/*
		t.Run("one empty", func(t *testing.T) {
			tx := handleTest.Begin()
			defer handleTest.Rollback()

			item := map[string]string{
				"Name": "Example",
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

			result := tx.
				Preload("Substance").
				Preload("Substitutes").
				Preload("Warehouses").
				Find(&itemDb, response.ID)

			assert.NoError(t, result.Error, "it shouldn't be an error")
			assert.NotEqual(t, int64(0), result.RowsAffected)
			assert.Equal(t, response.Name, itemDb.Name)

			assert.Equal(t, response.Name, itemDb.Name)
			assert.Equal(t, response.Sku, itemDb.Sku)
			assert.Equal(t, response.Ean, itemDb.Ean)
			assert.Equal(t, response.Taxed, itemDb.Taxed)
			assert.Equal(t, response.Iva, itemDb.Iva)
			assert.Equal(t, response.UrlImage, itemDb.UrlImage)
			assert.Equal(t, response.PackID, itemDb.PackID)
			assert.Equal(t, response.Pack, itemDb.Pack)
			assert.Equal(t, response.PackInfo, itemDb.PackInfo)
			assert.Equal(t, response.Quantity, itemDb.Quantity)
			assert.Equal(t, response.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, response.Weight, itemDb.Weight)
			assert.Equal(t, response.Height, itemDb.Height)
			assert.Equal(t, response.Width, itemDb.Width)
			assert.Equal(t, response.PackUnit, itemDb.PackUnit)
			assert.Equal(t, response.Depth, itemDb.Depth)
			assert.Equal(t, response.Keywords, itemDb.Keywords)
			assert.Equal(t, response.Variation, itemDb.Variation)
			assert.Equal(t, response.Wrapper, itemDb.Wrapper)
			assert.Equal(t, response.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, response.Status, itemDb.Status)
			assert.Equal(t, response.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, response.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, response.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, response.MakerID, itemDb.MakerID)
			assert.Equal(t, response.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, response.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, response.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, response.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, response.Substitutes, itemDb.Substitutes)
			assert.Equal(t, response.Substance, itemDb.Substance)
			assert.Equal(t, response.Warehouses, itemDb.Warehouses)
			assert.Equal(t, response.TypesProductID, itemDb.TypesProductID)
		})

		t.Run("one full", func(t *testing.T) {
			tx := handleTest.Begin()
			defer handleTest.Rollback()

			pack := model.Pack{}
			mark := model.Maker{}

			err := tx.Save(&pack).Error
			assert.NoError(t, err, "error when save pack")

			err = tx.Save(&mark).Error
			assert.NoError(t, err, "error when save mark")

			trakeMark := model.Trademark{MakerID: mark.ID}
			err = tx.Save(&trakeMark).Error
			assert.NoError(t, err, "error when save trademark")

			categoryOne := model.CategoryOne{}
			err = tx.Save(&categoryOne).Error
			assert.NoError(t, err, "error when save categoryOne")

			categoryTwo := model.CategoryTwo{CategoryOneID: categoryOne.ID}
			err = tx.Save(&categoryTwo).Error
			assert.NoError(t, err, "error when save categoryTwo")

			categoryThree := model.CategoryThree{
				CategoryOneID: categoryOne.ID,
				CategoryTwoID: categoryTwo.ID,
			}
			err = tx.Save(&categoryThree).Error
			assert.NoError(t, err, "error when save categoryThree")

			productSubstitutes := model.Product{}
			err = tx.Save(&productSubstitutes).Error
			assert.NoError(t, err, "error when save productSubstitutes")

			productSubstitutes2 := model.Product{}
			err = tx.Save(&productSubstitutes2).Error
			assert.NoError(t, err, "error when save productSubstitutes2")

			typesproduct := model.TypesProduct{}
			err = tx.Save(&typesproduct).Error
			assert.NoError(t, err, "error when save typesproduct")

			warehouses := []model.Warehouse{{Name: "wh-example"}, {Name: "wh-example2"}}
			err = tx.Save(&warehouses).Error
			assert.NoError(t, err, "error when save warehouses")

			subtances := []model.Substance{{}, {}}
			err = tx.Save(&warehouses).Error
			assert.NoError(t, err, "error when save subtances")

			item := map[string]interface{}{}
			err = json.Unmarshal(pro, &item)
			assert.NoError(t, err, "error when unpackage json")

			item["PackID"] = fmt.Sprint(pack.ID)
			item["MakerID"] = fmt.Sprint(mark.ID)
			item["TrademarkID"] = fmt.Sprint(trakeMark.ID)
			item["CategoryOneID"] = fmt.Sprint(categoryOne.ID)
			item["CategoryTwoID"] = fmt.Sprint(categoryTwo.ID)
			item["CategoryThreeID"] = fmt.Sprint(categoryThree.ID)
			item["SubstitutesIDS"] = fmt.Sprintf("%v, %v", productSubstitutes.ID, productSubstitutes2.ID)
			item["Substance"] = fmt.Sprintf("%v, %v", subtances[0].ID, subtances[1].ID)
			item["Warehouses"] = fmt.Sprintf("%v, %v", warehouses[0].ID, warehouses[1].ID)
			item["TypesProductID"] = fmt.Sprint(typesproduct.ID)

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

			result := tx.
				Preload("Substance").
				Preload("Substitutes").
				Preload("Warehouses").
				Find(&itemDb, response.ID)

			assert.NoError(t, result.Error, "it shouldn't be an error")

			assert.NotEqual(t, int64(0), result.RowsAffected)
			assert.Equal(t, response.Name, itemDb.Name)

			assert.Equal(t, response.Name, itemDb.Name)
			assert.Equal(t, response.Sku, itemDb.Sku)
			assert.Equal(t, response.Ean, itemDb.Ean)
			assert.Equal(t, response.Taxed, itemDb.Taxed)
			assert.Equal(t, response.Iva, itemDb.Iva)
			assert.Equal(t, response.UrlImage, itemDb.UrlImage)
			assert.Equal(t, response.PackID, itemDb.PackID)
			assert.Equal(t, response.Pack, itemDb.Pack)
			assert.Equal(t, response.PackInfo, itemDb.PackInfo)
			assert.Equal(t, response.Quantity, itemDb.Quantity)
			assert.Equal(t, response.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, response.Weight, itemDb.Weight)
			assert.Equal(t, response.Height, itemDb.Height)
			assert.Equal(t, response.Width, itemDb.Width)
			assert.Equal(t, response.PackUnit, itemDb.PackUnit)
			assert.Equal(t, response.Depth, itemDb.Depth)
			assert.Equal(t, response.Keywords, itemDb.Keywords)
			assert.Equal(t, response.Variation, itemDb.Variation)
			assert.Equal(t, response.Wrapper, itemDb.Wrapper)
			assert.Equal(t, response.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, response.Status, itemDb.Status)
			assert.Equal(t, response.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, response.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, response.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, response.MakerID, itemDb.MakerID)
			assert.Equal(t, response.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, response.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, response.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, response.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, len(response.Substitutes), len(itemDb.Substitutes), "it should be same count")

			for i := range response.Substitutes {
				assert.Equal(t, response.Substitutes[i].ID, itemDb.Substitutes[i].ID)
			}

			assert.Equal(t, response.Substance, itemDb.Substance)
			assert.Equal(t, response.Warehouses, itemDb.Warehouses)
			assert.Equal(t, response.TypesProductID, itemDb.TypesProductID)
		})

		t.Run("Many", func(t *testing.T) {
			tx := handleTest.Begin()
			defer handleTest.Rollback()

			items := []map[string]interface{}{}

			err := json.Unmarshal(prods, &items)
			assert.NoError(t, err, "error when unpackage json")

			body, err := json.Marshal(items)
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

			response := []model.Product{}
			err = json.Unmarshal([]byte(resBody), &response)
			assert.NoError(t, err, "json has a problem")

			for _, item := range response {
				assert.NotEqual(t, uint(0), item.ID, "it shouldn't be 0 a ID")
			}

			handleTest.ShowRequest(
				t,
				endpoint,
				handleTest.Post,
				body,
				resBody,
			)

			itemsDb := []model.Product{}

			result := tx.
				Preload("Substance").
				Preload("Substitutes").
				Preload("Warehouses").
				Find(&itemsDb)

			assert.NoError(t, result.Error, "it shouldn't be an error")
			assert.NotEqual(t, int64(0), result.RowsAffected)
			for i := range response {
				res := response[i]
				itemDb := itemsDb[i]
				assert.Equal(t, res.ID, itemDb.ID)
				assert.Equal(t, res.Name, itemDb.Name)
				assert.Equal(t, res.Sku, itemDb.Sku)
				assert.Equal(t, res.Ean, itemDb.Ean)
				assert.Equal(t, res.Taxed, itemDb.Taxed)
				assert.Equal(t, res.Iva, itemDb.Iva)
				assert.Equal(t, res.UrlImage, itemDb.UrlImage)
				assert.Equal(t, res.PackID, itemDb.PackID)
				assert.Equal(t, res.Pack, itemDb.Pack)
				assert.Equal(t, res.PackInfo, itemDb.PackInfo)
				assert.Equal(t, res.Quantity, itemDb.Quantity)
				assert.Equal(t, res.MakerUnit, itemDb.MakerUnit)
				assert.Equal(t, res.Weight, itemDb.Weight)
				assert.Equal(t, res.Height, itemDb.Height)
				assert.Equal(t, res.Width, itemDb.Width)
				assert.Equal(t, res.PackUnit, itemDb.PackUnit)
				assert.Equal(t, res.Depth, itemDb.Depth)
				assert.Equal(t, res.Keywords, itemDb.Keywords)
				assert.Equal(t, res.Variation, itemDb.Variation)
				assert.Equal(t, res.Wrapper, itemDb.Wrapper)
				assert.Equal(t, res.WrapperUnit, itemDb.WrapperUnit)
				assert.Equal(t, res.Status, itemDb.Status)
				assert.Equal(t, res.HandlesBog, itemDb.HandlesBog)
				assert.Equal(t, res.HandlesBaq, itemDb.HandlesBaq)
				assert.Equal(t, res.IDFloorProduct, itemDb.IDFloorProduct)
				assert.Equal(t, res.MakerID, itemDb.MakerID)
				assert.Equal(t, res.TrademarkID, itemDb.TrademarkID)
				assert.Equal(t, res.CategoryOneID, itemDb.CategoryOneID)
				assert.Equal(t, res.CategoryTwoID, itemDb.CategoryTwoID)
				assert.Equal(t, res.CategoryThreeID, itemDb.CategoryThreeID)

				assert.Equal(t, len(res.Substitutes), len(itemDb.Substitutes), "it should be same count")

				for i := range res.Substitutes {
					assert.Equal(t, res.Substitutes[i].ID, itemDb.Substitutes[i].ID)
				}

				assert.Equal(t, res.Substance, itemDb.Substance)
				assert.Equal(t, res.Warehouses, itemDb.Warehouses)
				assert.Equal(t, res.TypesProductID, itemDb.TypesProductID)
			}
		})

	*/
}

/*
func TestUpdate(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")

	t.Run("one", func(t *testing.T) {
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		item := model.Product{
			Name: "original name",
			Sku:  "examplesku",
		}
		err := tx.Save(&item).Error
		assert.NoError(t, err, "it shouldn't be an error")

		update := map[string]interface{}{
			"ID":     fmt.Sprint(item.ID),
			"Name":   "example2",
			"Status": true,
			"Sku":    item.Sku,
		}

		body, err := json.Marshal(update)
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

		assert.NotEqual(t,
			0,
			response,
			"it shouldn't be 0 a ID",
		)

		handleTest.ShowRequest(
			t,
			endpoint,
			handleTest.Put,
			body,
			resBody,
		)

		itemDb := model.Product{}

		result := tx.
			Preload("Substance").
			Preload("Substitutes").
			Preload("Warehouses").
			Find(&itemDb, response.ID)

		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, response.Name, itemDb.Name)

		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.Sku, itemDb.Sku)
		assert.Equal(t, response.Ean, itemDb.Ean)
		assert.Equal(t, response.Taxed, itemDb.Taxed)
		assert.Equal(t, response.Iva, itemDb.Iva)
		assert.Equal(t, response.UrlImage, itemDb.UrlImage)
		assert.Equal(t, response.PackID, itemDb.PackID)
		assert.Equal(t, response.Pack, itemDb.Pack)
		assert.Equal(t, response.PackInfo, itemDb.PackInfo)
		assert.Equal(t, response.Quantity, itemDb.Quantity)
		assert.Equal(t, response.MakerUnit, itemDb.MakerUnit)
		assert.Equal(t, response.Weight, itemDb.Weight)
		assert.Equal(t, response.Height, itemDb.Height)
		assert.Equal(t, response.Width, itemDb.Width)
		assert.Equal(t, response.PackUnit, itemDb.PackUnit)
		assert.Equal(t, response.Depth, itemDb.Depth)
		assert.Equal(t, response.Keywords, itemDb.Keywords)
		assert.Equal(t, response.Variation, itemDb.Variation)
		assert.Equal(t, response.Wrapper, itemDb.Wrapper)
		assert.Equal(t, response.WrapperUnit, itemDb.WrapperUnit)
		assert.Equal(t, response.Status, itemDb.Status)
		assert.Equal(t, response.HandlesBog, itemDb.HandlesBog)
		assert.Equal(t, response.HandlesBaq, itemDb.HandlesBaq)
		assert.Equal(t, response.IDFloorProduct, itemDb.IDFloorProduct)
		assert.Equal(t, response.MakerID, itemDb.MakerID)
		assert.Equal(t, response.TrademarkID, itemDb.TrademarkID)
		assert.Equal(t, response.CategoryOneID, itemDb.CategoryOneID)
		assert.Equal(t, response.CategoryTwoID, itemDb.CategoryTwoID)
		assert.Equal(t, response.CategoryThreeID, itemDb.CategoryThreeID)
		assert.Equal(t, response.Substitutes, itemDb.Substitutes)
		assert.Equal(t, response.Substance, itemDb.Substance)
		assert.Equal(t, response.Warehouses, itemDb.Warehouses)
		assert.Equal(t, response.TypesProductID, itemDb.TypesProductID)
	})

	t.Run("many", func(t *testing.T) {
		tx := handleTest.Begin()
		defer handleTest.Rollback()

		items := []model.Product{
			{
				Name: "example1",
				Sku:  "examplesku1",
			},
			{
				Name: "example2",
				Sku:  "examplesku2",
			},
			{
				Name: "example3",
				Sku:  "examplesku3",
			},
		}
		err := tx.Save(&items).Error
		assert.NoError(t, err, "it shouldn't be an error")

		update := []map[string]interface{}{}
		expected := append([]model.Product{}, items[0])

		for _, item := range items[1:] {
			update = append(update, map[string]interface{}{
				"ID":   fmt.Sprint(item.ID),
				"Name": item.Name + "update",
				"Sku":  item.Sku,
			})
			expected = append(expected, model.Product{
				CustomModel: item.CustomModel,
				Name:        item.Name + "update",
				Sku:         item.Sku,
			})
		}

		body, err := json.Marshal(update)
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

		response := []model.Product{}
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

		itemsDb := []model.Product{}

		result := tx.
			Preload("Substance").
			Preload("Substitutes").
			Preload("Warehouses").
			Find(&itemsDb)

		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)

		for i := range response {
			res := itemsDb[i]
			itemDb := expected[i]
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Sku, itemDb.Sku)
			assert.Equal(t, res.Ean, itemDb.Ean)
			assert.Equal(t, res.Taxed, itemDb.Taxed)
			assert.Equal(t, res.Iva, itemDb.Iva)
			assert.Equal(t, res.UrlImage, itemDb.UrlImage)
			assert.Equal(t, res.PackID, itemDb.PackID)
			assert.Equal(t, res.Pack, itemDb.Pack)
			assert.Equal(t, res.PackInfo, itemDb.PackInfo)
			assert.Equal(t, res.Quantity, itemDb.Quantity)
			assert.Equal(t, res.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, res.Weight, itemDb.Weight)
			assert.Equal(t, res.Height, itemDb.Height)
			assert.Equal(t, res.Width, itemDb.Width)
			assert.Equal(t, res.PackUnit, itemDb.PackUnit)
			assert.Equal(t, res.Depth, itemDb.Depth)
			assert.Equal(t, res.Keywords, itemDb.Keywords)
			assert.Equal(t, res.Variation, itemDb.Variation)
			assert.Equal(t, res.Wrapper, itemDb.Wrapper)
			assert.Equal(t, res.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, res.Status, itemDb.Status)
			assert.Equal(t, res.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, res.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, res.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, res.MakerID, itemDb.MakerID)
			assert.Equal(t, res.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, res.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, res.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, res.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, res.TypesProductID, itemDb.TypesProductID)
		}
	})
}

func TestGet(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")

	endpoint := "/product"
	products := []model.Product{
		{Name: "example", Sku: "13005", Ean: "a16asd0"},
		{Name: "example1", Sku: "1234", Ean: "a12"},
		{Name: "example2", Sku: "1235", Ean: "a13"},
		{Name: "example3", Sku: "1800", Ean: "a14"},
		{Name: "example4", Sku: "2122", Ean: "a12"},
		{Name: "gel", Sku: "3000", Ean: "a1400"},
		{Name: "gel2", Sku: "13000", Ean: "a160"},
		{Name: "a", Sku: "13001", Ean: "a"},
		{Name: "b", Sku: "13002", Ean: "a10"},
		{Name: "c", Sku: "13003", Ean: "a1asd60"},
		{Name: "d", Sku: "13004", Ean: "a16asd0"},
		{Name: "d", Sku: "13004", Ean: "a16asd0"},
		{Name: "d", Sku: "13004", Ean: "a16asd0"},
		{Name: "d", Sku: "13004", Ean: "a16asd0"},
		{Name: "d", Sku: "13004", Ean: "a16asd0"},
	}

	tx := handleTest.Begin()
	defer handleTest.Rollback()
	err = tx.Save(&products).Error
	assert.NoError(t, err, "error when save products")

	type Pag struct {
		Limit      int             `json:"limit,omitempty;query:limit"`
		Page       int             `json:"page,string,omitempty;query:page"`
		Sort       string          `json:"sort,omitempty;query:sort"`
		TotalRows  int64           `json:"totalRows"`
		TotalPages int             `json:"totalPages"`
		Rows       []model.Product `json:"Rows"`
	}

	t.Run("AllFirstPage", func(t *testing.T) {

		request := events.APIGatewayProxyRequest{
			HTTPMethod: handleTest.Get,
		}

		resBody := handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response := Pag{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetRequest(
			t, endpoint, resBody,
		)
		expected := []model.Product{}
		//count := 0
		for i := len(products) - 1; i >= 0; i-- {
			expected = append(expected, products[i])
			if len(expected) == len(response.Rows) {
				break
			}
		}

		fmt.Println()
		//itemDb := model.Product{}
		for i := range response.Rows {
			res := response.Rows[i]
			itemDb := expected[i]

			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Sku, itemDb.Sku)
			assert.Equal(t, res.Ean, itemDb.Ean)
			assert.Equal(t, res.Taxed, itemDb.Taxed)
			assert.Equal(t, res.Iva, itemDb.Iva)
			assert.Equal(t, res.UrlImage, itemDb.UrlImage)
			assert.Equal(t, res.PackID, itemDb.PackID)
			assert.Equal(t, res.Pack, itemDb.Pack)
			assert.Equal(t, res.PackInfo, itemDb.PackInfo)
			assert.Equal(t, res.Quantity, itemDb.Quantity)
			assert.Equal(t, res.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, res.Weight, itemDb.Weight)
			assert.Equal(t, res.Height, itemDb.Height)
			assert.Equal(t, res.Width, itemDb.Width)
			assert.Equal(t, res.PackUnit, itemDb.PackUnit)
			assert.Equal(t, res.Depth, itemDb.Depth)
			assert.Equal(t, res.Keywords, itemDb.Keywords)
			assert.Equal(t, res.Variation, itemDb.Variation)
			assert.Equal(t, res.Wrapper, itemDb.Wrapper)
			assert.Equal(t, res.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, res.Status, itemDb.Status)
			assert.Equal(t, res.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, res.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, res.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, res.MakerID, itemDb.MakerID)
			assert.Equal(t, res.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, res.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, res.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, res.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, res.TypesProductID, itemDb.TypesProductID)
		}
	})

	t.Run("AllFirstAndSecondPage", func(t *testing.T) {
		queryParams := map[string]string{
			"Limit": "2",
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

		response := Pag{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetRequest(
			t, endpoint, resBody,
		)

		expected := []model.Product{}
		//count := 0
		for i := len(products) - 1; i >= 0; i-- {
			expected = append(expected, products[i])
			if len(expected) == len(response.Rows) {
				break
			}
		}

		for i := range response.Rows {
			res := response.Rows[i]
			itemDb := expected[i]

			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Sku, itemDb.Sku)
			assert.Equal(t, res.Ean, itemDb.Ean)
			assert.Equal(t, res.Taxed, itemDb.Taxed)
			assert.Equal(t, res.Iva, itemDb.Iva)
			assert.Equal(t, res.UrlImage, itemDb.UrlImage)
			assert.Equal(t, res.PackID, itemDb.PackID)
			assert.Equal(t, res.Pack, itemDb.Pack)
			assert.Equal(t, res.PackInfo, itemDb.PackInfo)
			assert.Equal(t, res.Quantity, itemDb.Quantity)
			assert.Equal(t, res.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, res.Weight, itemDb.Weight)
			assert.Equal(t, res.Height, itemDb.Height)
			assert.Equal(t, res.Width, itemDb.Width)
			assert.Equal(t, res.PackUnit, itemDb.PackUnit)
			assert.Equal(t, res.Depth, itemDb.Depth)
			assert.Equal(t, res.Keywords, itemDb.Keywords)
			assert.Equal(t, res.Variation, itemDb.Variation)
			assert.Equal(t, res.Wrapper, itemDb.Wrapper)
			assert.Equal(t, res.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, res.Status, itemDb.Status)
			assert.Equal(t, res.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, res.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, res.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, res.MakerID, itemDb.MakerID)
			assert.Equal(t, res.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, res.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, res.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, res.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, res.TypesProductID, itemDb.TypesProductID)
		}

		queryParams = map[string]string{
			"Limit": "2",
			"Page":  "2",
		}

		request = events.APIGatewayProxyRequest{
			HTTPMethod:            handleTest.Get,
			QueryStringParameters: queryParams,
		}

		resBody = handleTest.UseHandleRequestTx(t,
			server.Product,
			request,
			http.StatusOK,
		)

		response = Pag{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetRequest(
			t, endpoint, resBody,
		)

		expected = []model.Product{}

		for i := len(products) - 1 - 2; i >= 0; i-- {
			expected = append(expected, products[i])
			if len(expected) == len(response.Rows) {
				break
			}
		}

		for i := range response.Rows {
			res := response.Rows[i]
			itemDb := expected[i]
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Name, itemDb.Name)
			assert.Equal(t, res.Sku, itemDb.Sku)
			assert.Equal(t, res.Ean, itemDb.Ean)
			assert.Equal(t, res.Taxed, itemDb.Taxed)
			assert.Equal(t, res.Iva, itemDb.Iva)
			assert.Equal(t, res.UrlImage, itemDb.UrlImage)
			assert.Equal(t, res.PackID, itemDb.PackID)
			assert.Equal(t, res.Pack, itemDb.Pack)
			assert.Equal(t, res.PackInfo, itemDb.PackInfo)
			assert.Equal(t, res.Quantity, itemDb.Quantity)
			assert.Equal(t, res.MakerUnit, itemDb.MakerUnit)
			assert.Equal(t, res.Weight, itemDb.Weight)
			assert.Equal(t, res.Height, itemDb.Height)
			assert.Equal(t, res.Width, itemDb.Width)
			assert.Equal(t, res.PackUnit, itemDb.PackUnit)
			assert.Equal(t, res.Depth, itemDb.Depth)
			assert.Equal(t, res.Keywords, itemDb.Keywords)
			assert.Equal(t, res.Variation, itemDb.Variation)
			assert.Equal(t, res.Wrapper, itemDb.Wrapper)
			assert.Equal(t, res.WrapperUnit, itemDb.WrapperUnit)
			assert.Equal(t, res.Status, itemDb.Status)
			assert.Equal(t, res.HandlesBog, itemDb.HandlesBog)
			assert.Equal(t, res.HandlesBaq, itemDb.HandlesBaq)
			assert.Equal(t, res.IDFloorProduct, itemDb.IDFloorProduct)
			assert.Equal(t, res.MakerID, itemDb.MakerID)
			assert.Equal(t, res.TrademarkID, itemDb.TrademarkID)
			assert.Equal(t, res.CategoryOneID, itemDb.CategoryOneID)
			assert.Equal(t, res.CategoryTwoID, itemDb.CategoryTwoID)
			assert.Equal(t, res.CategoryThreeID, itemDb.CategoryThreeID)
			assert.Equal(t, res.TypesProductID, itemDb.TypesProductID)
		}
	})

	t.Run("byId", func(t *testing.T) {
		queryParams := map[string]string{
			"ID": fmt.Sprint(products[len(products)-1].ID),
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
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetByRequest(
			t,
			endpoint,
			resBody,
			queryParams,
		)

		itemDb := model.Product{}

		result := tx.Find(&itemDb, response.ID)

		assert.NoError(t, result.Error, "it shouldn't be an error")
		assert.NotEqual(t, int64(0), result.RowsAffected)
		assert.Equal(t, response.Name, itemDb.Name)
		assert.Equal(t, response.Sku, itemDb.Sku)
		assert.Equal(t, response.Ean, itemDb.Ean)
	})

	t.Run("byName", func(t *testing.T) {
		pattern := "exam"
		queryParams := map[string]string{
			"Name": pattern,
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

		response := Pag{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetByRequest(
			t,
			endpoint,
			resBody,
			queryParams,
		)

		produtsFilter := []model.Product{}
		for i := len(products) - 1; i >= 0; i-- {
			prod := products[i]
			isMatch := regexp.MustCompile(pattern).MatchString(prod.Name)
			if isMatch {
				produtsFilter = append(produtsFilter, prod)
			}
		}

		for i, item := range response.Rows {
			assert.Equal(t, item.ID, produtsFilter[i].ID)
			assert.Equal(t, item.Name, produtsFilter[i].Name)
			assert.Equal(t, item.Sku, produtsFilter[i].Sku)
			assert.Equal(t, item.Ean, produtsFilter[i].Ean)
		}
	})
	t.Run("bySkus", func(t *testing.T) {
		pattern := "1234,1235,1800,21222,3000,13000,13001,13002"
		queryParams := map[string]string{
			"Skus": pattern,
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

		//model.Pagination{}

		response := []model.Product{}
		err = json.Unmarshal([]byte(resBody), &response)
		assert.NoError(t, err, "json has a problem")

		handleTest.ShowGetByRequest(
			t,
			endpoint,
			resBody,
			queryParams,
		)
		produtsFilter := []model.Product{}
		for _, prod := range products {
			skus := strings.Split(pattern, ",")
			for _, sku := range skus {
				if prod.Sku == sku {
					produtsFilter = append(produtsFilter, prod)
				}
			}
		}

		for i, item := range response {
			assert.Equal(t, item.ID, produtsFilter[i].ID)
			assert.Equal(t, item.Name, produtsFilter[i].Name)
			assert.Equal(t, item.Sku, produtsFilter[i].Sku)
			assert.Equal(t, item.Ean, produtsFilter[i].Ean)
		}

	})

}

*/
