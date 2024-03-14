package makers

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

func NewMakerRepository(db *gorm.DB) controller.CRUD {
	return &MakerRepository{
		grepo: genericrepository.New[*model.Maker](
			db, []string{"Trademarks"},
		),
	}
}

type MakerRepository struct {
	grepo repository.GenericRepository[*model.Maker]
}

// Get -----------------------------------------------------------
func (p *MakerRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *MakerRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }
func (p *MakerRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.Maker{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *MakerRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *MakerRepository) InsertMany(reqBody []byte) response.Result {
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
func (p *MakerRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *MakerRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}
