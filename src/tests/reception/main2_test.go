package reception

/*
func TestSuccess2(t *testing.T) {
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

	t.Run("Simple CRUD", func(t *testing.T) {
		//example := model.Reception{}
		supplier := model.Supplier{}
		warehouse := model.Warehouse{}

		handleTest.Db.Save(&supplier)
		handleTest.Db.Save(&warehouse)
		purchase := model.Purchase{
			SupplierID:  supplier.ID,
			WarehouseID: warehouse.ID,
		}
		handleTest.Db.Save(&purchase)
		article := model.Article{
			PurchaseID: purchase.ID,
			Count:      4,
		}
		handleTest.Db.Save(&article)

		defer handleTest.DeleteInfo(t, deletes)

		t.Run("Insert(indirect)", func(t *testing.T) {
			example := model.ReceptionArt{
				ArticleID: article.ID,
				Count:     article.Count,
			}
			example2 := []model.ReceptionArt{}

			for i := 0; i < 2; i++ {
				example2 = append(example2, model.ReceptionArt{
					ArticleID: article.ID,
					Missing:   uint(i + 1),
					Refund:    uint(i + 2),
					Reason:    uint8(i + 3),
					Count:     uint(i + 4),
					Batch:     fmt.Sprint(i),
					Date:      time.Now().Add(time.Hour * 1),
				})
			}

			t.Run("one", func(t *testing.T) {
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
					server.Reception,
					request,
					http.StatusOK,
				)

				res := model.ReceptionArt{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Error(err.Error())
				}

				example.CustomModel = res.CustomModel
				handleTest.DifferentMap(t,
					example,
					res,
				)

				purchaseUpdate := model.Purchase{}
				handleTest.Db.Find(&purchaseUpdate, purchase.ID)

				// deberia ser 2 porque la compra esta completa
				if purchaseUpdate.ReceptionStatus != 2 {
					t.Errorf("deberia ser status 2 porque la compra estaria completa")
				}
			})
			t.Run("update", func(t *testing.T) {
				example.Count = article.Count + 1
				body, err := json.Marshal(example)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
					return
				}

				request := events.APIGatewayProxyRequest{
					HTTPMethod: http.MethodPut,
					Body:       string(body),
				}

				resBody := handleTest.UseHandleRequest(t,
					server.Reception,
					request,
					http.StatusInternalServerError,
				)
				// retorna un error
				// res := model.ReceptionArt{}
				// if err := json.Unmarshal([]byte(resBody), &res); err != nil {
				// 	t.Error(err.Error())
				// }
				fmt.Println(resBody)
				purchaseUpdate := model.Purchase{}
				handleTest.Db.Find(&purchaseUpdate, purchase.ID)

				if purchaseUpdate.ReceptionStatus != 2 {
					t.Errorf("deberia ser status 2 porque la compra estaria completa")
				}
			})

			t.Run("delete", func(t *testing.T) {
				body, err := json.Marshal(example)
				if err != nil {
					t.Errorf("Mal formato de json: %s", err.Error())
					return
				}

				request := events.APIGatewayProxyRequest{
					HTTPMethod: http.MethodDelete,
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
				fmt.Println(resBody)

				articleUpdate := model.Article{}
				handleTest.Db.
					Preload("ReceptionInfo").
					Find(&articleUpdate, article.ID)

				if j, _ := json.MarshalIndent(articleUpdate, "", "\t"); true {
					fmt.Printf("\n\n%v\n\n", string(j))
				}
				if len(articleUpdate.ReceptionInfo) != 0 {
					t.Errorf("deberia ser cero porque fue borrada la recepcion")
				}

			})

			/*
					t.Run("many", func(t *testing.T) {
						body, err := json.Marshal(example2)
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
						fmt.Println(resBody)
						/*
							purchase.ID = res.ID
							purchase.CreatedAt = res.CreatedAt
							purchase.Reception = res.Reception

							receptionsIds.Add(res.Reception.ID)
							example.ID = res.Reception.ID
							example.PurchaseID = res.ID
					})

				t.Run("GetByPurchaseId", func(t *testing.T) {

					path := map[string]string{
						"ArticleID": fmt.Sprint(article.ID),
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod:     http.MethodGet,
						PathParameters: path,
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
					fmt.Println(resBody)
					/*
						purchase.ID = res.ID
						purchase.CreatedAt = res.CreatedAt
						purchase.Reception = res.Reception

						receptionsIds.Add(res.Reception.ID)
						example.ID = res.Reception.ID
						example.PurchaseID = res.ID
				})

		})
	})
}

*/
