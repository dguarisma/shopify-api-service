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

func NewCategoryTwoRepository(db *gorm.DB) controller.CRUD {
	return &CategoryTwoRepository{
		grepo: genericrepository.New[*model.CategoryTwo](
			db, []string{},
		),
	}
}

type CategoryTwoRepository struct {
	grepo repository.GenericRepository[*model.CategoryTwo]
}

// Get -----------------------------------------------------------
func (p *CategoryTwoRepository) GetAll() response.Result { return p.grepo.GetAll() }

func (p *CategoryTwoRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *CategoryTwoRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.CategoryTwo{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *CategoryTwoRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElementTwo(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *CategoryTwoRepository) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleAdaptTwo(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if err := p.grepo.InsertMany(elements); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return controller.FormateBody(elements, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *CategoryTwoRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElementTwo(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *CategoryTwoRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(
			err, http.StatusBadRequest,
		)
	}
	return p.grepo.DeleteByID(id)
}
