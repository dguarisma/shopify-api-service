package reception_test

/*
func TestSuccess(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}

		deletes := []interface{}{
			&model.ArticlesReception{},
			&model.Article{},
			&model.Product{},
			&model.Receiving{},
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
		article := model.Article{PurchaseID: purchase.ID}
		handleTest.Db.Save(&article)
		fmt.Println(article)
		defer handleTest.DeleteInfo(t, deletes)

		receptionsIds := set.New[uint]()

		example := model.Receiving{
			Articles: []model.ArticlesReception{},
		}

			// the reception create is indirect
			t.Run("Insert(indirect)", func(t *testing.T) {
				t.Run("one", func(t *testing.T) {
					supplier := model.Supplier{
						BusinessName: "BusinessName1",
					}
					warehouse := model.Warehouse{
						Name:       "Warehouse1",
						Department: "exmaple",
					}
					if err := handleTest.Db.Save(&supplier).Error; err != nil {
						t.Fatalf("a %v", err.Error())
					}
					if err := handleTest.Db.Save(&warehouse).Error; err != nil {
						t.Fatalf("a %v", err.Error())
					}

					products := []model.Product{
						{
							Name: "prod1",
							Sku:  "prod1",
							Ean:  "prod1",
							Iva:  10,
						},
						{
							Name: "prod2",
							Sku:  "prod2",
							Ean:  "prod2",
							Iva:  20,
						},
					}

					if err := handleTest.Db.Save(&products).Error; err != nil {
						t.Fatalf("a %v", err.Error())
					}

					pro := []uint{}
					for _, p := range products {
						pro = append(pro, p.ID)
					}

					products2 := &[]model.Product{}
					fmt.Println(pro)
					handleTest.Db.Find(products2, pro)

					purchase := model.Purchase{
						Discount: 10,
						Notes:    "esto es una nota",
						//Articles:          make([]model.Article, 2),
						Articles:          make([]model.Article, 2),
						SupplierID:        supplier.ID,
						WarehouseID:       warehouse.ID,
						DiscountEarliyPay: 100,
					}
					rand.Seed(time.Now().UnixNano())

					for i, art := range purchase.Articles {
						//purchase.Articles[i].ProductID = products[prodRand].ID
						art = model.Article{
							ProductID:          products[i].ID,
							BasePrice:          float32(rand.Intn(101)),
							Bonus:              uint(rand.Intn(102)),
							Discount:           float32(rand.Intn(103)),
							DiscountAdditional: float32(rand.Intn(104)),
							Count:              uint(rand.Intn(105)),
							SubTotal:           float32(rand.Intn(106)),
							Total:              float32(rand.Intn(107)),
							Tax:                float32(rand.Intn(107)),
						}
						purchase.Articles[i] = art
					}

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

					purchase.ID = res.ID
					purchase.CreatedAt = res.CreatedAt
					purchase.Reception = res.Reception

					receptionsIds.Add(res.Reception.ID)
					example.ID = res.Reception.ID
					example.PurchaseID = res.ID

					recep := model.Receiving{}
					handleTest.Db.Find(&recep, model.Receiving{
						PurchaseID: purchase.ID,
					})

					//for _, purArt := range purchase.Articles {
					for _, purArt := range res.Articles {
						recep.Articles = append(recep.Articles, model.ArticlesReception{

							Count:     purArt.Count - 1,
							ArticleID: purArt.ID,
						})
					}

					if j, _ := json.MarshalIndent(recep.Articles, "", "\t"); true {
						fmt.Printf("\n\n%v\n\n", string(j))
					}
					if err := handleTest.Db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&recep).Error; err != nil {
						t.Error(err.Error())
					}

					if j, _ := json.MarshalIndent(recep, "", "\t"); true {
						fmt.Printf("\n\n%v\n\n", string(j))
					}

						handleTest.DifferentMap(
							t,
							purchase,
							res,
						)
				})

		/*
				t.Run("many", func(t *testing.T) {
					supplier := model.Supplier{
						BusinessName: "BusinessName1",
					}
					supplier2 := model.Supplier{
						BusinessName: "BusinessName2",
					}
					warehouse := model.Warehouse{
						Name: "Warehouse1",
					}
					warehouse2 := model.Warehouse{
						Name: "Warehouse1",
					}
					handleTest.Db.Save(&supplier)
					handleTest.Db.Save(&warehouse)
					handleTest.Db.Save(&supplier2)
					handleTest.Db.Save(&warehouse2)

					purchases := []model.Purchase{
						{
							Articles:    []model.Article{},
							SupplierID:  supplier.ID,
							WarehouseID: warehouse.ID,
						},
						{
							Articles:    []model.Article{},
							SupplierID:  supplier.ID,
							WarehouseID: warehouse.ID,
						},
						{
							Articles:    []model.Article{},
							SupplierID:  supplier2.ID,
							WarehouseID: warehouse2.ID,
						},
						{
							Articles:    []model.Article{},
							SupplierID:  supplier2.ID,
							WarehouseID: warehouse.ID,
						},
					}

					body, err := json.Marshal(purchases)
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

					res := []model.Purchase{}

					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}
					for i, purchaRes := range res {
						purchases[i].ID = purchaRes.ID
						purchases[i].CreatedAt = purchaRes.CreatedAt
						purchases[i].Reception = purchaRes.Reception
						receptionsIds.Add(purchaRes.Reception.ID) // luego posiblemente borrado
					}
					handleTest.DifferentMap(
						t,
						purchases,
						res,
					)
				})
			})
			// se deberia realizar una peticion al la base de datos
			// para obtener el reception id
			t.Run("Get", func(t *testing.T) {

				t.Run("ById", func(t *testing.T) {
					path := map[string]string{
						"ID": fmt.Sprint(example.ID),
					}

					request := events.APIGatewayProxyRequest{
						HTTPMethod:            http.MethodGet,
						QueryStringParameters: path,
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Reception,
						request,
						http.StatusOK,
					)

					res := model.Receiving{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					fmt.Println(resBody)
					handleTest.DifferentMap(
						t,
						example,
						res,
					)
				})

				t.Run("All", func(t *testing.T) {
					request := events.APIGatewayProxyRequest{
						HTTPMethod: http.MethodGet,
					}

					resBody := handleTest.UseHandleRequest(t,
						server.Reception,
						request,
						http.StatusOK,
					)

					res := []model.Receiving{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					receptionIds := receptionsIds.Get()
					if len(res) != len(receptionIds) {
						t.Errorf("Deberia ser la misma cantidad de receptions")
					}

					for _, reception := range res {
						receptionsIds.Delete(reception.ID)

					}
					if len(receptionsIds.Get()) != 0 {
						t.Errorf("Deberia existir estos ids %v",
							receptionsIds.Get(),
						)
					}

						handleTest.DifferentMap(
							t,
							example,
							res,
						)
				})
			})

			t.Run("Update", func(t *testing.T) {
				supplier := model.Supplier{}
				warehouse := model.Warehouse{}
				handleTest.Db.Save(&supplier)
				handleTest.Db.Save(&warehouse)

				purchase := model.Purchase{
					Articles:    []model.Article{},
					SupplierID:  supplier.ID,
					WarehouseID: warehouse.ID,
				}
				handleTest.Db.Save(&purchase)

				example = model.Receiving{
					CustomModel: model.CustomModel{ID: example.ID},
					PurchaseID:  purchase.ID,
					//				Missing:     1,
					//				Refund:      2,
					//				Reason:      3,
					Status:   4,
					Articles: []model.ArticlesReception{},
				}
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
					http.StatusOK,
				)

				res := model.Receiving{}
				if err := json.Unmarshal([]byte(resBody), &res); err != nil {
					t.Error(err.Error())
				}

				handleTest.DifferentMap(
					t,
					example,
					res,
				)
				// check update
				temp := model.Receiving{}
				handleTest.Db.First(&temp, example.ID)

				temp.CreatedAt = example.CreatedAt
				temp.DeletedAt = example.DeletedAt
				temp.UpdatedAt = example.UpdatedAt
				// this is for dont have articles
				temp.Articles = example.Articles

				handleTest.DifferentMap(
					t,
					example,
					temp,
				)
			})

	})

	t.Run("Complex CRUD", func(t *testing.T) {
		/*
			defer handleTest.DeleteInfo(t, deletes)
			receptionsIds := set.New[uint]()

			supplier := model.Supplier{}
			warehouse := model.Warehouse{}
			handleTest.Db.Save(&supplier)
			handleTest.Db.Save(&warehouse)

			products := []model.Product{
				{
					Name:       "Name1",
					Sku:        "Sk1",
					Ean:        "Ean1",
					Warehouses: []*model.Warehouse{
						//	{CustomModel: model.CustomModel{ID: warehouse.ID}},
					},
				},
				{
					Name:       "Name2",
					Sku:        "Sk2",
					Ean:        "Ean2",
					Warehouses: []*model.Warehouse{
						//	{CustomModel: model.CustomModel{ID: warehouse.ID}},
					},
				},
			}
			handleTest.Db.Save(&products)

			example := model.Receiving{
				//Articles: []model.ArticlesReception{},
			}

			purchase := model.Purchase{
				Articles:    make([]model.Article, 0, len(products)+1),
				SupplierID:  supplier.ID,
				WarehouseID: warehouse.ID,
			}
			for _, p := range products {
				purchase.Articles = append(purchase.Articles, model.Article{
					ProductID: p.ID,
					Count:     10,
				})

			}

			purchase.Articles = append(purchase.Articles, model.Article{
				ProductID: products[0].ID,
				Count:     10,
			})

			t.Run("Insert(indirect)", func(t *testing.T) {
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

				// it's created indirect in the database
				purchase.ID = res.ID
				purchase.CreatedAt = res.CreatedAt
				purchase.Reception = res.Reception

				for i, artRq := range purchase.Articles {
					artRs := res.Articles[i]
					if artRq.ProductID != artRs.ProductID {
						t.Errorf(
							"El id del producto deberia ser el mismo %v - %v",
							artRq.ProductID,
							res.Articles[i].ProductID,
						)
					}
					artRq.ID = artRs.ID
					artRq.PurchaseID = artRs.PurchaseID
					purchase.Articles[i] = artRq
				}

				handleTest.DifferentMap(
					t,
					purchase,
					res,
				)

				receptionsIds.Add(res.Reception.ID)

				example = *purchase.Reception
				temp := model.Receiving{}
				handleTest.Db.First(&temp, model.Receiving{PurchaseID: purchase.ID})

				// this is because i dont need create, update, delete
				temp.CustomModel = model.CustomModel{ID: temp.ID}

				handleTest.DifferentMap(
					t,
					example,
					temp,
				)

			})
			t.Run("Update", func(t *testing.T) {
				t.Run("add_article", func(t *testing.T) {
					example.Articles = append(example.Articles, model.ArticlesReception{
						ArticleID: purchase.Articles[0].ID,
						Count:     purchase.Articles[0].Count - 1,
						Date:      time.Now(),
					})

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
						http.StatusOK,
					)

					res := model.Receiving{}
					if err := json.Unmarshal([]byte(resBody), &res); err != nil {
						t.Error(err.Error())
					}

					for i, art := range res.Articles {
						example.Articles[i].ID = art.ID
						example.Articles[i].ReceivingID = example.ID
					}

					// check update
					temp := model.Receiving{}
					handleTest.Db.Preload("Articles").First(&temp, example.ID)

					temp.CustomModel = model.CustomModel{ID: temp.ID}
					for i, art := range temp.Articles {
						temp.Articles[i].CustomModel = model.CustomModel{ID: art.ID}
					}
					example.Articles[0].Date = temp.Articles[0].Date

					handleTest.DifferentMap(
						t,
						example,
						temp,
					)
				})
				t.Run("update_article", func(t *testing.T) {
					t.Run("success", func(t *testing.T) {
						example.Articles[0].Count = example.Articles[0].Count + 1

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
							http.StatusOK,
						)

						res := model.Receiving{}
						if err := json.Unmarshal([]byte(resBody), &res); err != nil {
							t.Error(err.Error())
						}

						for i, art := range res.Articles {
							example.Articles[i].ID = art.ID
							example.Articles[i].ReceivingID = example.ID
						}

						// check update
						temp := model.Receiving{}
						handleTest.Db.Preload("Articles").First(&temp, example.ID)

						temp.CustomModel = model.CustomModel{ID: temp.ID}
						for i, art := range temp.Articles {
							temp.Articles[i].CustomModel = model.CustomModel{ID: art.ID}
						}
						example.Articles[0].Date = temp.Articles[0].Date

						handleTest.DifferentMap(
							t,
							example,
							temp,
						)
					})
					t.Run("fail", func(t *testing.T) {
						// este va a fallar por que la cantidad(count) de
						// articulos(reception) nunca debe ser mayor que la cantidad que hay en la
						// comprar(pruchase) de determinado articulo
						example.Articles[0].Count = example.Articles[0].Count + 1
						body, err := json.Marshal(example)
						if err != nil {
							t.Errorf("Mal formato de json: %s", err.Error())
							return
						}

						t.Log("hay q cambiarle el estado aca y la respuesta")
						request := events.APIGatewayProxyRequest{
							HTTPMethod: http.MethodPut,
							Body:       string(body),
						}

						resBody := handleTest.UseHandleRequest(t,
							server.Reception,
							request,
							http.StatusInternalServerError,
						)

						res := response.ErrMsg{}
						if err := json.Unmarshal([]byte(resBody), &res); err != nil {
							t.Error(err.Error())
						}

						data := response.ErrMsg{
							Error: model.ExceededTheLimit(
								purchase.ID,
								example.Articles[0].ArticleID,
								int(purchase.Articles[0].Count),
								int(example.Articles[0].Count),
							).Error(),
						}

						handleTest.DifferentMap(
							t,
							data,
							res,
						)
					})
					t.Run("chance purchase status", func(t *testing.T) {
						example.Articles[0].Count = purchase.Articles[0].Count
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
							http.StatusOK,
						)

						res := model.Receiving{}
						if err := json.Unmarshal([]byte(resBody), &res); err != nil {
							t.Error(err.Error())
						}

						handleTest.DifferentMap(
							t,
							example,
							res,
						)

						// check update
						articles := []model.Article{}
						handleTest.Db.Find(
							&articles, purchase.Articles[0].ID,
						)
					})
					t.Run("chance purchase status2", func(t *testing.T) {})
				})
				t.Run("delete_article", func(t *testing.T) {})
				t.Run("add_and_update", func(t *testing.T) {})
			})

			/*
				t.Run("Get", func(t *testing.T) {
					t.Run("ById", func(t *testing.T) {
						path := map[string]string{
							"ID": fmt.Sprint(example.ID),
						}

						request := events.APIGatewayProxyRequest{
							HTTPMethod:            http.MethodGet,
							QueryStringParameters: path,
						}

						resBody := handleTest.UseHandleRequest(t,
							server.Reception,
							request,
							http.StatusOK,
						)

						res := model.Receiving{}
						if err := json.Unmarshal([]byte(resBody), &res); err != nil {
							t.Error(err.Error())
						}

						handleTest.DifferentMap(
							t,
							example,
							res,
						)
					})
				})


	})

}
*/
