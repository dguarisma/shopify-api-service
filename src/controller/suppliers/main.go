package suppliers

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

func NewSupplierRepository(db *gorm.DB) controller.CRUD {
	return &SupplierRepository{
		grepo: genericrepository.New[*model.Supplier](
			db, []string{},
		),
	}
}

type SupplierRepository struct {
	grepo repository.GenericRepository[*model.Supplier]
}

// Get -----------------------------------------------------------
func (p *SupplierRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *SupplierRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *SupplierRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.Supplier{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *SupplierRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *SupplierRepository) InsertMany(reqBody []byte) response.Result {
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
func (p *SupplierRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *SupplierRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}
