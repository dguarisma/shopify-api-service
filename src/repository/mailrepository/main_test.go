package mailrepository_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository/mailrepository"
	"desarrollosmoyan/lambda/src/tests/utils"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMsg2(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	email := "example"

	tx := handleTest.Begin()
	defer handleTest.Rollback()
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
		EmailContact: email,
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

	product := model.Product{
		Name: "exa1",
		Sku:  "exa1",
		Ean:  "exa1",
	}

	err = tx.Save(&supplier).Error
	assert.NoError(t, err)
	err = tx.Save(&warehouse).Error
	assert.NoError(t, err)
	err = tx.Save(&product).Error
	assert.NoError(t, err)

	example := model.Purchase{
		Status:         1,
		Tax:            float32(rand.Int31n(30)),
		DiscountGlobal: float32(rand.Int31n(30)),
		SubTotal:       float32(rand.Int31n(30)),
		Total:          float32(rand.Int31n(100)),
		Articles:       []model.Article{},
		SupplierID:     supplier.ID,
		WarehouseID:    warehouse.ID,
	}

	for i := 0; i < 2; i++ {
		example.Articles = append(example.Articles, model.Article{
			ProductID: product.ID,
			Count:     uint(rand.Intn(100)),
			BasePrice: float32(rand.Int31n(30)),
			Tax:       float32(rand.Int31n(30)),
			Discount:  float32(rand.Int31n(30)),
			Bonus:     uint(rand.Int31n(30)),
			SubTotal:  float32(rand.Int31n(30)),
			Total:     float32(rand.Int31n(100)),
		})
	}
	fmt.Println(tx)
	mailRepo := mailrepository.New(tx)
	_, err = mailRepo.GetMsg2(&example)
	assert.NoError(t, err)

}
