package reception_test

/*
func TestReceptionStatus(t *testing.T) {
	fmt.Println("hoas")
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	deletes := []interface{}{
		&model.ReceptionArt{},
		&model.Article{},
		&model.Product{},
		&model.Purchase{},
		&model.Supplier{},
		&model.Warehouse{},
	}
	handleTest.DeleteInfo(t, deletes)
	defer handleTest.DeleteInfo(t, deletes)

	t.Run("example 1", func(t *testing.T) {

		supplier := model.Supplier{}
		warehouse := model.Warehouse{}
		product1 := model.Product{
			Name: "example1",
			Ean:  "example1",
			Sku:  "example1",
		}

		handleTest.Db.Save(&supplier)
		handleTest.Db.Save(&warehouse)
		handleTest.Db.Save(&product1)

		purchase := model.Purchase{
			WarehouseID: warehouse.ID,
			SupplierID:  supplier.ID,
			Articles: []model.Article{
				{
					ProductID: product1.ID,
					Count:     10,
				},

				{
					ProductID: product1.ID,
					Count:     100,
				},
			},
		}

		t.Run("Insert(purchase)", func(t *testing.T) {
			body, err := json.Marshal(purchase)
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

			purchase.CustomModel = res.CustomModel
			purchase.CreatedAt = res.CreatedAt

			for i := range purchase.Articles {
				purchase.Articles[i].CustomModel = res.Articles[i].CustomModel
				purchase.Articles[i].PurchaseID = res.ID
			}

			handleTest.DifferentMap(t,
				purchase,
				res,
			)
		})

		t.Run("Insert(purchase)", func(t *testing.T) {
			body, err := json.Marshal(purchase)
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

			purchase.CustomModel = res.CustomModel
			purchase.CreatedAt = res.CreatedAt

			for i := range purchase.Articles {
				purchase.Articles[i].CustomModel = res.Articles[i].CustomModel
				purchase.Articles[i].PurchaseID = res.ID
			}

			handleTest.DifferentMap(t,
				purchase,
				res,
			)

			// chequeando estado del purchaDb
			purchaDb := model.Purchase{CustomModel: model.CustomModel{ID: purchase.ID}}

			handleTest.Db.Find(&purchaDb)
			if purchaDb.ReceptionStatus != 0 {
				t.Errorf("el estado de la reception tiene que ser 0 porque no hay novedad")
			}
		})
		reception := model.ReceptionArt{
			ArticleID: purchase.Articles[0].ID,
			Count:     purchase.Articles[0].Count,
		}

		t.Run("Insert(reception) check reception status 1", func(t *testing.T) {
			body, err := json.Marshal(reception)
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
			reception.CustomModel = res.CustomModel

			handleTest.DifferentMap(t,
				reception,
				res,
			)

			// chequeando estado del reception
			reception := model.Purchase{CustomModel: model.CustomModel{ID: purchase.ID}}

			handleTest.Db.Find(&reception)

			if reception.ReceptionStatus != 1 {
				t.Errorf("el estado de la reception tiene que ser 1 porque esta a medio completar")
			}
		})

		receptions := []model.ReceptionArt{
			{
				ArticleID: purchase.Articles[1].ID,
				Count:     10,
			},
			{
				ArticleID: purchase.Articles[1].ID,
				Count:     10,
			},
			{
				ArticleID: purchase.Articles[1].ID,
				Count:     30,
			},
			{
				ArticleID: purchase.Articles[1].ID,
				Count:     50,
			},
		}

		t.Run("Insert(reception) check reception status 2", func(t *testing.T) {
			body, err := json.Marshal(receptions)
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

			res := []model.ReceptionArt{}
			if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				t.Error(err.Error())
			}
			for i, recRes := range res {
				receptions[i].CustomModel = recRes.CustomModel
			}

			handleTest.DifferentMap(t,
				receptions,
				res,
			)

			// chequeando estado del reception
			reception := model.Purchase{CustomModel: model.CustomModel{ID: purchase.ID}}

			handleTest.Db.Find(&reception)

			if reception.ReceptionStatus != 2 {
				t.Errorf("el estado de la reception tiene que ser 2 porque esta a medio completar")
			}
		})

		t.Run("Update(purchase) check reception still status 2", func(t *testing.T) {
			purchase.Discount = 10
			body, err := json.Marshal(purchase)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPut,
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

			purchase.CustomModel = res.CustomModel
			purchase.CreatedAt = res.CreatedAt

			handleTest.DifferentMap(t,
				purchase,
				res,
			)

			// chequeando estado del reception
			reception := model.Purchase{CustomModel: model.CustomModel{ID: purchase.ID}}

			handleTest.Db.Find(&reception)

			if reception.ReceptionStatus != 2 {
				t.Errorf("el estado de la reception tiene que ser 2 porque esta a completo")
			}
		})

		t.Run("Update(purchase) check reception status 3", func(t *testing.T) {
			purchase.ReceptionStatus = 3

			body, err := json.Marshal(purchase)
			if err != nil {
				t.Errorf("Mal formato de json: %s", err.Error())
				return
			}

			request := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPut,
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

			purchase.CustomModel = res.CustomModel
			purchase.CreatedAt = res.CreatedAt

			handleTest.DifferentMap(t,
				purchase,
				res,
			)

			// chequeando estado del reception
			reception := model.Purchase{CustomModel: model.CustomModel{ID: purchase.ID}}

			handleTest.Db.Find(&reception)

			if reception.ReceptionStatus != 3 {
				t.Errorf("el estado de la reception tiene que ser 3 porque esta a cancelado")
			}
		})
	})
}

*/
