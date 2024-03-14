package shopifyserv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	_ "embed"
)

var (
	//go:embed getItemInventory.graphql
	getItemInventoryWithProductId string

	//go:embed mutationInventoryAjustQuantities.graphql
	mutationInventoryAjustQuantities string
)

func (ss *ShopifyServ) handlerItem(item Item, warehouse WareHouseShopify) (productRes *InventoryItemAdjustments, err error) {
	t := template.
		Must(template.New("template").
			Parse(getItemInventoryWithProductId))

	var b bytes.Buffer
	if err := t.Execute(&b, item); err != nil {
		return nil, err
	}

	query := Query{
		Query: b.String(),
	}
	res, status, err := ss.graphqlQuery(query, warehouse)
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, ErrStatusUnexpected(200, status)
	}

	Res := RootEntity{}
	if err := json.Unmarshal(res, &Res); err != nil {
		return nil, err //////
	}

	if Res.Data.Product == nil {
		return nil, fmt.Errorf("hubo un error en la peticion seguramente este ProductId no corresponde a la bodega")
	}

	productRes = &InventoryItemAdjustments{
		InventoryItemId: Res.Data.Product.Variants.Edges[0].Node.InventoryItem.Id,
		AvailableDelta:  int64(item.Available),
	}
	return productRes, nil
}

func (ss *ShopifyServ) SendUpdateInventory(products ProductsUpdates) error {

	if err := ss.GetWareHouses(products.WareHouse); err != nil {
		return err
	}

	example := ProductsUpdate{
		InventoryItemAdjustments: make(
			[]InventoryItemAdjustments,
			0, len(products.Items),
		),
	}

	items := map[string]uint64{}

	warehouseInfo := ss.warehouses[products.WareHouse]
	example.LocationId = warehouseInfo.AdminGraphqlApiId

	for _, i := range products.Items {
		curItem := items[i.ProductShopifyID]
		curItem += uint64(i.Available)
		items[i.ProductShopifyID] = curItem
	}

	for id, item := range items {
		i := Item{
			Available:        int64(item),
			ProductShopifyID: id,
		}
		itemAdjustment, err := ss.handlerItem(i, warehouseInfo)
		if err != nil {
			return fmt.Errorf("ProductoId(%v) Bodega(%v): error:%v ", i.ProductShopifyID, products.WareHouse, err)
		}
		example.InventoryItemAdjustments = append(example.InventoryItemAdjustments, *itemAdjustment)
	}

	t := template.
		Must(template.New("template").
			Parse(mutationInventoryAjustQuantities))

	var b bytes.Buffer
	if err := t.Execute(&b, example); err != nil {
		return err
	}
	query := Query{
		Query:     b.String(),
		Variables: example,
	}
	res, status, err := ss.graphqlQuery(query, warehouseInfo)
	if err != nil {
		return err
	}
	if status != 200 {
		return ErrStatusUnexpected(200, status)
	}
	res2 := RootEntity2{}

	if err := json.Unmarshal(res, &res2); err != nil {
		return err
	}
	if res2.Data.InventoryBulkAdjustQuantityAtLocation.UserErrors != nil {
		err := ""
		for _, v := range res2.Data.InventoryBulkAdjustQuantityAtLocation.UserErrors {
			err += fmt.Sprintf("%v ", v.Message)
		}
		if strings.TrimSpace(err) == "" {
			return nil
		}
		return fmt.Errorf("error al actualizar shopify: %v", err)
	}
	// agregar el chequeo de error
	return nil
}
