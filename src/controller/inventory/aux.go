package inventory

import (
	"desarrollosmoyan/lambda/src/model"
	"fmt"
	"math"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ErrNotSellOrPurchases = "No hay ventas ni compras registradas"
)

type InventoryInfo struct {
	model.Pagination
	Sku         string `json:"Sku"`
	WarehouseID uint   `json:"WarehouseID,string"`
	From        string `json:"From"`
	To          string `json:"To"`
}

type IHandlerInventory interface {
	GetPurchaseInventory(inventoryInfo *InventoryInfo) (*[]InventorySell, error)
	GetSelledInventory(inventoryInfo *InventoryInfo) (*[]InventorySell, error)
	SortTables(purchase, sell []InventorySell) (*[]InventorySell, error)
	CompleteTables(fullTable *[]InventorySell)
}

type HandlerTables struct{ Db *gorm.DB }

func NewHandlerTablesInventory(db *gorm.DB) IHandlerInventory {
	return &HandlerTables{db}
}

type Store struct {
	Units int64
	Cost  float32
}

func (HandlerTables) SetDate(date time.Time) string { return date.Format("02/01/2006") }

func (ht *HandlerTables) AvgCost(stores ...Store) float32 {
	var totalCost float32
	var totalUnits int64

	for _, store := range stores {
		totalCost += store.Cost * float32(store.Units)
		totalUnits += store.Units
	}

	if totalUnits == 0 {
		return 0
	}
	result := math.Round(float64(totalCost) / float64(totalUnits))

	return float32(result)
}

type InventorySell struct {
	Date          string
	ProductId     string // it's a shopify product id
	ProductName   string
	Sku           string
	InitInventory Store
	Purchase      Store
	Available     Store
	Selled        Store
	EndInventory  Store
}

func (ht *HandlerTables) SortTables(purchase, sell []InventorySell) (*[]InventorySell, error) {
	fulltable := []InventorySell{}
	sellHandler := NewHandlerItem(&sell)
	purchaseHandler := NewHandlerItem(&purchase)

	if sellHandler.IsEmpty() && purchaseHandler.IsEmpty() {
		return nil, fmt.Errorf(ErrNotSellOrPurchases)
	}

	for { // pasar el campo init al end
		if !sellHandler.IsExist() && !purchaseHandler.IsExist() {
			break
		}

		if !sellHandler.IsExist() {
			fulltable = append(fulltable, purchaseHandler.GetItem())
			purchaseHandler.Next()
			continue
		}
		if !purchaseHandler.IsExist() {
			fulltable = append(fulltable, sellHandler.GetItem())
			sellHandler.Next()
			continue
		}
		/*
			if purchaseHandler.GetDate().Equal(sellHandler.GetDate()) {
				purcha := purchaseHandler.GetItem()
				selled := sellHandler.GetItem()
				purcha.Selled = selled.Selled
				fulltable = append(fulltable, purcha)
				sellHandler.Next()
				purchaseHandler.Next()
				continue
			}
		*/

		if purchaseHandler.GetDate().After(sellHandler.GetDate()) {
			fulltable = append(fulltable, purchaseHandler.GetItem())
			purchaseHandler.Next()
			continue
		}

		fulltable = append(fulltable, sellHandler.GetItem())
		sellHandler.Next()
	}

	return &fulltable, nil
}

func (ht *HandlerTables) CompleteTables(fullTable *[]InventorySell) {
	for i, row := range *fullTable {

		row.Available.Cost = ht.AvgCost(row.InitInventory, row.Purchase)
		row.Available.Units = row.InitInventory.Units + row.Purchase.Units
		if row.Selled.Units == 0 {
			row.EndInventory = row.Available
		} else {
			row.EndInventory.Cost = row.Selled.Cost
			row.EndInventory.Units = row.Available.Units - row.Selled.Units
		}

		(*fullTable)[i] = row
	}
}

func (ht *HandlerTables) GetPurchaseInventory(inventoryInfo *InventoryInfo) (*[]InventorySell, error) {
	type PurchaseList struct {
		Count       uint64
		Date        time.Time
		HandlesBaq  string
		HandlesBog  string
		Item        uint64
		Name        string
		City        string
		Sku         string
		PriceTotal  float32
		CostPerItem float32
	}

	purchaseList := []PurchaseList{}

	query := `
	select
	  count(r.id)
	from reception_arts r
	inner join articles a on a.id = r.article_id
	inner join products p on p.id = a.product_id
	inner join purchases pur on pur.id = a.purchase_id
	inner join warehouses w on w.id = pur.warehouse_id `

	query += fmt.Sprintf(`where w.id = %v `, inventoryInfo.WarehouseID)

	if inventoryInfo.Sku != "" {
		query += fmt.Sprintf(`and sku = '%v' `, inventoryInfo.Sku)
	}

	if inventoryInfo.From != "" && inventoryInfo.To != "" {
		query += fmt.Sprintf(`and r.date between '%v' and '%v' `, inventoryInfo.From, inventoryInfo.To)
	}

	result := ht.Db.Debug().
		Raw(query).Scan(&inventoryInfo.TotalPages)
	query = `
	select
	  r.count as item,
	  round(a.base_price / a.count, 2) as CostPerItem,
	  r.date,
	  p.name,
	  p.sku,
	  p.handles_bog,
	  p.handles_baq,
	  w.city
	from reception_arts r
	inner join articles a on a.id = r.article_id
	inner join products p on p.id = a.product_id
	inner join purchases pur on pur.id = a.purchase_id
	inner join warehouses w on w.id = pur.warehouse_id
	`
	query += fmt.Sprintf(`where w.id = %v `, inventoryInfo.WarehouseID)

	if inventoryInfo.Sku != "" {
		query += fmt.Sprintf(`and sku = '%v' `, inventoryInfo.Sku)
	}

	if inventoryInfo.From != "" && inventoryInfo.To != "" {
		query += fmt.Sprintf(`and r.date between '%v' and '%v' `, inventoryInfo.From, inventoryInfo.To)
	}

	query += `order by r.date desc `
	query += fmt.Sprintf(`limit %v offset %v;`, inventoryInfo.GetLimit(), inventoryInfo.GetOffset())

	// agregar el null
	result = ht.Db.Raw(query).Scan(&purchaseList)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return &[]InventorySell{}, nil
	}

	purchaseTable := make([]InventorySell, len(purchaseList))
	for i, purchaseRow := range purchaseList {

		item := InventorySell{
			Sku:         purchaseRow.Sku,
			ProductName: purchaseRow.Name,
			Date:        ht.SetDate(purchaseRow.Date),
			Purchase: Store{
				// Cost:  purchaseRow.PriceTotal / float32(purchaseRow.Count),
				Cost:  purchaseRow.CostPerItem,
				Units: int64(purchaseRow.Item),
			},
		}

		if purchaseRow.City == "Bogotá" {
			item.ProductId = purchaseRow.HandlesBog
		} else if purchaseRow.City == "Barranquilla" {
			item.ProductId = purchaseRow.HandlesBaq
		} else {
			return nil, fmt.Errorf("no existe la bodega(%v) en ", purchaseRow.City)
		}

		// ver si tiene acento
		purchaseTable[i] = item
	}

	return &purchaseTable, nil
}

func (ht *HandlerTables) GetSelledInventory(inventoryInfo *InventoryInfo) (*[]InventorySell, error) {

	var city string
	ht.Db.Raw(
		`select w.city from warehouses w where w.id = ?`,
		inventoryInfo.WarehouseID,
	).Scan(&city)
	if city == "Bogotá" {
		city = "Bogota"
	}

	if city == "" {
		return nil, fmt.Errorf("no existe una bodega con el id(%v)", inventoryInfo.WarehouseID)
	}

	uri := os.Getenv("DB_SHOPIFY")
	db, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//example := []map[string]interface{}{}
	example := []struct {
		//Warehouse string
		Sku      string
		Price    float32
		Quantity int
		Date     time.Time
		Stock    int64
	}{}

	var total int64

	query := `select count(sku) from customer_orders_products `
	query += fmt.Sprintf(`where warehouse = '%v' `, city)

	if inventoryInfo.Sku != "" {
		query += fmt.Sprintf(`and sku = '%v' `, inventoryInfo.Sku)
	}

	if inventoryInfo.From != "" && inventoryInfo.To != "" {
		query += fmt.Sprintf(`and order_date between '%v' and '%v'`, inventoryInfo.From, inventoryInfo.To)
	}
	/*
		query = `
		select
		  count(sku)
		from
		  customer_orders_products
		where warehouse = ?
		and sku = ?
		and order_date between ? and ?;`
	*/

	if err := db.Debug().Raw(query). //city,
		//inventoryInfo.Sku,
		//inventoryInfo.From,
		//inventoryInfo.To,

		Scan(&total).Error; err != nil {
		return nil, err
	}
	inventoryInfo.TotalRows += total

	query = `select
	  sku,
	  price,
	  quantity,
	  order_date as date,
	  gmv_by_product as stock
	from
	  customer_orders_products `

	query += fmt.Sprintf(`where warehouse = '%v' `, city)

	if inventoryInfo.Sku != "" {
		query += fmt.Sprintf(`and sku = '%v' `, inventoryInfo.Sku)
	}

	if inventoryInfo.From != "" && inventoryInfo.To != "" {
		query += fmt.Sprintf(`and order_date between '%v' and '%v'`, inventoryInfo.From, inventoryInfo.To)
	}
	query += `order by date desc `

	query += fmt.Sprintf(`limit %v offset %v;`, inventoryInfo.GetLimit(), inventoryInfo.GetOffset())

	/*
		query = `
		select
		  sku,
		  price,
		  quantity,
		  order_date as date,
		  gmv_by_product as stock
		from
		  customer_orders_products
		where warehouse = ?
		and sku = ?
		and order_date between ? and ?
		order by date desc
		limit ? offset ?;`
	*/

	if err := db.Debug().Raw(query). //city,
		// inventoryInfo.Sku,
		// inventoryInfo.From,
		// inventoryInfo.To,
		// inventoryInfo.Pagination.GetLimit(),
		// inventoryInfo.Pagination.GetOffset(),

		Scan(&example).Error; err != nil {
		return nil, err
	}
	fmt.Println("\n\n\n\n", example)

	res := make([]InventorySell, len(example))
	for i := range res {
		shopifyItem := example[i]
		res[i] = InventorySell{
			Date: ht.SetDate(shopifyItem.Date),
			Sku:  shopifyItem.Sku, // luego buscar el producto id
			//ProductId: shopifyItem.Sku, // luego buscar el producto id
			Selled: Store{
				Cost:  shopifyItem.Price,
				Units: int64(shopifyItem.Quantity),
			},
			EndInventory: Store{
				Units: shopifyItem.Stock,
			},
		}
	}
	return &res, nil
}
