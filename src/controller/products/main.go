package products

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/repository/genericrepository"
	"desarrollosmoyan/lambda/src/repository/productrepository"
	"desarrollosmoyan/lambda/src/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func NewProductRepository(db *gorm.DB) controller.CRUD {
	return &ProductRepository{
		repo: productrepository.New(
			db,
			genericrepository.New[*model.Product](db, []string{
				"Maker.Trademarks",
				"Pack",
				"CategoryOne",
				"CategoryTwo",
				"CategoryThree",
				"Warehouses",
				"Substance",
				"Substitutes",
			}),
		),
	}
}

type ProductRepository struct {
	repo repository.ProductRepository
}

// Get -----------------------------------------------------------
func (p *ProductRepository) GetAll() response.Result           { return p.repo.GetFirstPage() }
func (p *ProductRepository) GetByID(id string) response.Result { return p.repo.GetByID(id) }

func (p *ProductRepository) GetManyBy(pathMap []byte) response.Result {
	query := map[string]string{}
	if err := json.Unmarshal(pathMap, &query); err != nil {
		return response.NewFailResult(err, 400)
	}
	if name, ok := query["Name"]; ok {
		pag := model.Pagination{}
		if err := json.Unmarshal(pathMap, &pag); err != nil {
			return response.NewFailResult(err, 400)
		}
		return p.repo.GetByField(name, pag, repository.Name)
	}
	if sku, ok := query["Sku"]; ok {
		return p.repo.GetByField(sku, model.Pagination{}, repository.Sku)
	}
	if Ean, ok := query["Ean"]; ok {
		return p.repo.GetByField(Ean, model.Pagination{}, repository.Ean)
	}

	if values, ok := query["Skus"]; ok {
		skus := strings.Split(values, ",")
		return p.repo.GetBySkus(skus)
	}

	fmt.Print("1")
	return p.repo.GetByPagination(pathMap)
}

// Post -----------------------------------------------------------
func (p *ProductRepository) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.Insert(&element)
}

func (p *ProductRepository) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleAdapt(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if err := p.repo.InsertsAndUpdates(elements); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return controller.FormateBody(elements, http.StatusOK)
}

// Update ---------------------------------------------------------
func (p *ProductRepository) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.Insert(&element)
}

// Delete ---------------------------------------------------------
func (p *ProductRepository) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.DeleteByID(id)
}
