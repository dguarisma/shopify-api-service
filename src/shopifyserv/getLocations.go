package shopifyserv

import (
	"encoding/json"
	"fmt"
)

func (ss *ShopifyServ) GetWareHouses(locationName string) error {
	wh, ok := ss.warehouses[locationName]
	if !ok {
		return ErrDoesntExistWarehouse(locationName)
	}

	// si ya esta seteado no hace falta volver a llamar
	if wh.AdminGraphqlApiId != "" {
		return nil
	}

	url := fmt.Sprintf("%s/admin/api/%s/locations.json",
		wh.UrlBase,
		ss.version,
	)

	data, err := ss.getRequest(url, wh.XShopifyAccessToken)
	if err != nil {
		return err
	}

	locations := Location{}
	if err := json.Unmarshal(data, &locations); err != nil {
		return err
	}

	for _, location := range locations.Locations {
		if wh, ok := ss.warehouses[location.Name]; ok {
			wh.AdminGraphqlApiId = location.AdminGraphqlApiId
			ss.warehouses[location.Name] = wh
		}
	}
	return nil
}
