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

func NewCategoryThreeRepository(db *gorm.DB) controller.CRUD {
	return &CategoryThreeRepository{
		grepo: genericrepository.New[*model.CategoryThree](
			db, []string{},
		),
	}
}

type CategoryThreeRepository struct {
	grepo repository.GenericRepository[*model.CategoryThree]
}

// Get -----------------------------------------------------------
func (p *CategoryThreeRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *CategoryThreeRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *CategoryThreeRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.CategoryThree{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *CategoryThreeRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElementThree(reqBody)
	if err != nil {
		return response.NewFailResult(
			err, http.StatusBadRequest,
		)
	}
	return p.grepo.Insert(&element)
}

func (p *CategoryThreeRepository) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleAdaptThree(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if err := p.grepo.InsertMany(elements); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return controller.FormateBody(elements, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *CategoryThreeRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElementThree(reqBody)
	if err != nil {
		return response.NewFailResult(
			err, http.StatusBadRequest,
		)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *CategoryThreeRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(
			err, http.StatusBadRequest,
		)
	}
	return p.grepo.DeleteByID(id)
}
