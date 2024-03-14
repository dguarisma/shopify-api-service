package packs

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

func NewPackRepository(db *gorm.DB) controller.CRUD {
	return &PackRepository{
		grepo: genericrepository.New[*model.Pack](
			db, []string{},
		),
	}
}

type PackRepository struct {
	grepo repository.GenericRepository[*model.Pack]
}

// Get -----------------------------------------------------------
func (p *PackRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *PackRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *PackRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.Pack{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *PackRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *PackRepository) InsertMany(reqBody []byte) response.Result {
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
func (p *PackRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *PackRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}
