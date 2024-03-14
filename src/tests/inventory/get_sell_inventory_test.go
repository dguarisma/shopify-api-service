package inventory_test

/*
func TestGetSellinventory(t *testing.T) {

	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	tx := handleTest.Begin()
	defer handleTest.Rollback()

	wh := model.Warehouse{City: "Bogotá"}
	err = tx.Save(&wh).Error
	assert.NoError(t, err)

	handlerTables := inventory.HandlerTables{Db: tx}
	i := inventory.InventoryInfo{
		From:        "2023/01/11",
		To:          "2023/01/12",
		WarehouseID: wh.ID,
		Pagination:  model.Pagination{Limit: 10},
	}
	data, err := handlerTables.GetSelledInventory(i)
	assert.NoError(t, err)

	resInventory := inventory.HandlerTables{}
	fmt.Println(data)
}

func TestGetAllinventory(t *testing.T) {

	handleTest, err := utils.NewHandleTest()
	if err != nil {
		log.Fatal(err.Error())
	}
	tx := handleTest.Begin()
	defer handleTest.Rollback()

	wh := model.Warehouse{City: "Bogotá"}
	err = tx.Save(&wh).Error
	assert.NoError(t, err)

	handlerTables := inventory.HandlerTables{Db: tx}
	i := inventory.InventoryInfo{
		From:        "2023/01/11",
		To:          "2023/01/12",
		WarehouseID: wh.ID,
	}
	data, err := handlerTables.GetSelledInventory(i)
	assert.NoError(t, err)
	fmt.Println(data)
}
*/
