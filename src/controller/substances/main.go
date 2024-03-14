package substances

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

func NewSubstanceRepository(db *gorm.DB) controller.CRUD {
	return &SubstanceRepository{
		grepo: genericrepository.New[*model.Substance](db, []string{}),
	}
}

type SubstanceRepository struct {
	grepo repository.GenericRepository[*model.Substance]
}

// Get -----------------------------------------------------------
func (p *SubstanceRepository) GetAll() response.Result           { return p.grepo.GetAll() }
func (p *SubstanceRepository) GetByID(id string) response.Result { return p.grepo.GetByID(id) }

func (p *SubstanceRepository) GetManyBy(pathMap []byte) response.Result {
	element := &model.Substance{}
	if err := json.Unmarshal(pathMap, element); err != nil {
		return controller.ErrImposibleFormat()
	}
	return p.grepo.GetManyBy(element)
}

// Post -----------------------------------------------------------
func (p *SubstanceRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

func (p *SubstanceRepository) InsertMany(reqBody []byte) response.Result {
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
func (p *SubstanceRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *SubstanceRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.grepo.DeleteByID(id)
}

/*
func NewSubstanceRepository(db *gorm.DB) controller.CRUD {
	return &SubstanceRepository{
		repo: database.Repo{
			Db: db,
		},
	}
}

type SubstanceRepository struct {
	repo database.Repo
}

// Get -----------------------------------------------------------
func (p *SubstanceRepository) GetAll() response.Result {
	return p.repo.GetAll(&[]model.Substance{})
}

func (p *SubstanceRepository) GetByID(id string) response.Result {
	return p.repo.GetByID(&model.Substance{}, id)
}

func (p *SubstanceRepository) GetManyBy(pathMap []byte) response.Result {
	return p.repo.GetManyBy(pathMap, &model.Substance{}, &[]model.Substance{})
}

// Post -----------------------------------------------------------
func (p *SubstanceRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.Insert(&element)
}

func (p *SubstanceRepository) InsertMany(reqBody []byte) response.Result {
	newElement, updateElement, err := HandleNewElements(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}

	if len(newElement) == 0 || len(updateElement) == 0 {
		if len(newElement) != 0 {
			return p.repo.Insert(&newElement)
		}

		if len(updateElement) != 0 {
			return p.repo.Insert(&updateElement)
		}

		return response.NewFailResult(
			fmt.Errorf("Fallo de formato"), http.StatusBadRequest,
		)
	}

	if err := p.repo.InsertMany(&newElement, &updateElement); err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}

	responseBody := updateElement
	for _, update := range newElement {
		responseBody = append(responseBody, update)
	}
	return controller.FormateBody(responseBody, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *SubstanceRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *SubstanceRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.DeleteByID(&model.Substance{}, id)
}

*/
