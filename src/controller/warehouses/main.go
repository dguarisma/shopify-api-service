package warehouses

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/repository/genericrepository"
	"desarrollosmoyan/lambda/src/response"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func NewWarehouseRepository(db *gorm.DB) controller.CRUD {
	return &WarehouseRepository{
		grepo: genericrepository.New[*model.Warehouse](db, []string{}),
	}
}

type WarehouseRepository struct {
	grepo repository.GenericRepository[*model.Warehouse]
}

// Get -----------------------------------------------------------
func (p *WarehouseRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *WarehouseRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *WarehouseRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.Warehouse{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *WarehouseRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *WarehouseRepository) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleAdapt(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if err := p.grepo.InsertsAndUpdates(elements); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return controller.FormateBody(elements, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *WarehouseRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Update(&element)
}

// Delete ---------------------------------------------------------
func (p *WarehouseRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}
