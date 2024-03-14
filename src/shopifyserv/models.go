package shopifyserv

import "net/http"

type ShopifyService interface {
	SendUpdateInventory(items ProductsUpdates) error
}

type ProductsUpdate struct {
	InventoryItemAdjustments []InventoryItemAdjustments `json:"inventoryItemAdjustments"`
	LocationId               string                     `json:"-"`
}

type InventoryItemAdjustments struct {
	InventoryItemId string `json:"inventoryItemId"`
	AvailableDelta  int64  `json:"availableDelta"`
}

type ShopifyCredentials struct {
	Name                string
	UrlBase             string
	XShopifyAccessToken string
}

type WareHouseShopify struct {
	UrlBase             string
	XShopifyAccessToken string
	AdminGraphqlApiId   string
	//Id                  uint64
}

type ShopifyServ struct {
	warehouses map[string]WareHouseShopify
	urlBase    string
	version    string
	client     *http.Client
}

type ShopifyConfig struct {
	UrlBase             string
	XShopifyAccessToken string
	Version             string
}

type Item struct {
	ProductShopifyID string
	Available        int64
}

type ProductsUpdates struct {
	Items     []Item
	WareHouse string
}

type ErrorElement struct {
	errorDescription string
	element          uint64
	count            uint64
}

type Location struct {
	Locations []struct {
		ID                uint64 `json:"id"`
		Name              string `json:"name"`
		AdminGraphqlApiId string `json:"admin_graphql_api_id"`
	} `json:"locations"`
}

type Product struct {
	Product struct {
		Variants []struct {
			InventoryItemId   uint64 `json:"inventory_item_id"`
			AdminGraphqlApiId string `json:"admin_graphql_api_id"`
		} `json:"variants"`
	} `json:"product"`
}

type Query struct {
	Query     string `json:"query"`
	Variables any    `json:"variables,omitempty"`
}

type RootEntity struct {
	Data struct {
		Product *struct {
			Variants struct {
				Edges []struct {
					Node struct {
						InventoryItem struct {
							Id string `json:"id"`
						} `json:"inventoryItem"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"variants"`
		} `json:"product"`
	} `json:"data"`
}

type RootEntity2 struct {
	Data struct {
		InventoryBulkAdjustQuantityAtLocation struct {
			UserErrors []struct {
				Field   []string `json:"field"`
				Message string   `json:"message"`
			} `json:"userErrors"`
		} `json:"inventoryBulkAdjustQuantityAtLocation"`
	} `json:"data"`
}
