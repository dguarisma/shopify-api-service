package reception_test

/*
func TestMsg(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	active := false
	if active {

		deletes := []interface{}{
			&model.Product{},
			&model.Article{},
			&model.Purchase{},
			&model.Supplier{},
			&model.Warehouse{},
		}
		defer handleTest.DeleteInfo(t, deletes)

		suppl := model.Supplier{}
		ware := model.Warehouse{Name: "example1"}
		ware2 := model.Warehouse{Name: "example2"}
		handleTest.Db.Save(&suppl)
		handleTest.Db.Save(&ware)
		handleTest.Db.Save(&ware2)

		products := []model.Product{
			{
				Name: "example1",
				Sku:  time.Now().String(),
				Ean:  time.Now().String(),
				Warehouses: []*model.Warehouse{
					{CustomModel: model.CustomModel{ID: ware.ID}},
					{CustomModel: model.CustomModel{ID: ware2.ID}},
				},
			},
			{
				Name: "example2",
				Sku:  time.Now().String(),
				Ean:  time.Now().String(),
				Warehouses: []*model.Warehouse{
					{CustomModel: model.CustomModel{ID: ware.ID}},
				},
			},
			{
				Name: "example3",
				Sku:  time.Now().String(),
				Ean:  time.Now().String(),
				Warehouses: []*model.Warehouse{
					{CustomModel: model.CustomModel{ID: ware2.ID}},
				},
			},
		}

		handleTest.Db.Create(&products)

		purchas := model.Purchase{
			SupplierID:  suppl.ID,
			WarehouseID: ware.ID,
			Articles:    make([]model.Article, len(products)),
			Reception:   &model.Receiving{},
		}

		for i, pro := range products {
			purchas.Articles[i] = model.Article{ProductID: pro.ID}
		}

		handleTest.Db.Save(&purchas)

		data := &receptions.Adapt{
			PurchaseID: purchas.ID,
			Articles:   []receptions.ArticlesAdapt{},
		}

		countExample := []uint{1, 5, 10}
		for i, art := range purchas.Articles {
			data.Articles = append(data.Articles, receptions.ArticlesAdapt{
				ArticleID: art.ID,
				Batch:     "aa",
				Count:     countExample[i],
			})
		}

	}
}
*/
