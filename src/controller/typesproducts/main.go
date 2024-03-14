package typesproducts

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

func NewTypesProductRepository(db *gorm.DB) controller.CRUD {
	return &TypesProductRepository{
		grepo: genericrepository.New[*model.TypesProduct](db, []string{}),
	}
}

type TypesProductRepository struct {
	grepo repository.GenericRepository[*model.TypesProduct]
}

// Get -----------------------------------------------------------
func (p *TypesProductRepository) GetAll() response.Result { return p.grepo.GetAll() }

func (p *TypesProductRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }
func (p *TypesProductRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.TypesProduct{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *TypesProductRepository) InsertOne(pathMap []byte) response.Result {
	element := &model.TypesProduct{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

func (p *TypesProductRepository) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleAdapt(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if err := p.grepo.InsertMany(elements); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return controller.FormateBody(elements, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *TypesProductRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Update(&element)
}

// Delete ---------------------------------------------------------
func (p *TypesProductRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}
