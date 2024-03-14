package categories

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

func NewCategoryOneRepository(db *gorm.DB) controller.CRUD {
	return &CategoryOneRepository{
		grepo: genericrepository.New[*model.CategoryOne](
			db, []string{},
		),
	}
}

type CategoryOneRepository struct {
	grepo repository.GenericRepository[*model.CategoryOne]
}

// Get -----------------------------------------------------------
func (p *CategoryOneRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *CategoryOneRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *CategoryOneRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.CategoryOne{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *CategoryOneRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElementOne(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *CategoryOneRepository) InsertMany(reqBody []byte) response.Result {
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
func (p *CategoryOneRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElementOne(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *CategoryOneRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(
			err, http.StatusBadRequest,
		)
	}
	return p.grepo.DeleteByID(id)
}
