package inventory

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/response"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"gorm.io/gorm"
)

func NewInventoryRepository(db *gorm.DB) controller.CRUD {
	return &InventoryRepository{
		handlerInventory: NewHandlerTablesInventory(db),
	}
}

type InventoryRepository struct{ handlerInventory IHandlerInventory }

// Get -----------------------------------------------------------
func (p *InventoryRepository) GetManyBy(pathMap []byte) response.Result {
	inventoryInfoSearch := InventoryInfo{}

	if err := json.Unmarshal(pathMap, &inventoryInfoSearch); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}

	if inventoryInfoSearch.Limit == 0 {
		inventoryInfoSearch.Limit = 20
	}

	data, err := p.handlerInventory.GetPurchaseInventory(&inventoryInfoSearch)
	if err != nil {
		return response.NewFailResult(err, 500)
	}

	shopifyData, err := p.handlerInventory.GetSelledInventory(&inventoryInfoSearch)
	if err != nil {
		return response.NewFailResult(err, 500)
	}

	data2, err := p.handlerInventory.SortTables(*data, *shopifyData)
	if err != nil {
		return response.NewFailResult(err, 500)
	}
	inventoryInfoSearch.Pagination.Rows = data2

	inventoryInfoSearch.TotalPages = int(math.Ceil(float64(inventoryInfoSearch.TotalRows / int64(inventoryInfoSearch.GetLimit()))))

	return controller.FormateBody(inventoryInfoSearch.Pagination, http.StatusOK)
}

func (p *InventoryRepository) GetAll() response.Result                   { return MethodUseIncorret() }
func (p *InventoryRepository) GetByID(id string) response.Result         { return MethodUseIncorret() }
func (p *InventoryRepository) InsertOne(reqBody []byte) response.Result  { return MethodDontAllow() }
func (p *InventoryRepository) InsertMany(reqBody []byte) response.Result { return MethodDontAllow() }
func (p *InventoryRepository) Update(reqBody []byte) response.Result     { return MethodDontAllow() }
func (p *InventoryRepository) DeleteById(reqBody []byte) response.Result { return MethodDontAllow() }

func MethodDontAllow() response.Result {
	return response.NewFailResult(
		fmt.Errorf("este metodo no esta disponible"), http.StatusForbidden,
	)
}

func MethodUseIncorret() response.Result {
	return response.NewFailResult(
		fmt.Errorf("este metodo requiere[ProductId(midas), WarehouseID(midas)]"),
		http.StatusBadRequest,
	)
}
