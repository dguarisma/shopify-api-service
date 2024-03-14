package productrepository

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/response"
	"net/http"

	"gorm.io/gorm"
)

func New(db *gorm.DB, grepo repository.GenericRepository[*model.Product]) repository.ProductRepository {
	return &ProductRepo{
		Db:                db,
		GenericRepository: grepo,
	}
}

type ProductRepo struct {
	Db       *gorm.DB
	Preloads []string // son para buscar datos anidados
	repository.GenericRepository[*model.Product]
}

func (pr *ProductRepo) GetByField(word string, pag model.Pagination, field repository.Field) response.Result {
	query := pr.Db.Model(&model.Product{})
	for _, preload := range pr.Preloads {
		query = query.Preload(preload)
	}

	items := []model.Product{}
	var result *gorm.DB

	switch field {
	case repository.Name:
		if pag.Limit == 0 {
			pag.Limit = 100
		}
		where := query.Where("name LIKE ?", "%"+word+"%")

		pagination, err := model.PaginateBy(&items, where, &pag)
		if err != nil {
			return response.NewFailResult(err, 500)
		}
		return controller.FormateBody(pagination, http.StatusOK)

	case repository.Sku:
		result = query.Where("sku LIKE ?", word+"%").Find(&items)
	case repository.Ean:
		result = query.Where("ean LIKE ?", word+"%").Find(&items)
	}

	if result, problem := controller.HandleResultSearchDB(result); problem {
		return result
	}
	return controller.FormateBody(items, http.StatusOK)
}

func (pr *ProductRepo) GetBySkus(skus []string) response.Result {
	query := pr.Db
	for _, preload := range pr.Preloads {
		query = query.Preload(preload)
	}
	re := []model.Product{}

	if err := pr.Db.Model(&re).Where("sku in (?)", skus).Find(&re).Error; err != nil {
		return response.NewFailResult(err, 500)
	}
	return controller.FormateBody(re, http.StatusOK)
}

/*
func (pr *ProductRepo) GetByPagination(pathMap []byte, body interface{}, res interface{}) response.Result {
	pagination, err := model.GetPagination(pathMap)
	if err != nil {
		return controller.ErrPagination(err)
	}
	if err := json.Unmarshal(pathMap, body); err != nil {
		return controller.ErrImposibleFormat()
	}

	query := pr.Db
	for _, preload := range pr.Preloads {
		query = query.Preload(preload)
	}

	result := query.Scopes(model.Paginate(res, pagination, pr.Db)).Find(res, body)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	pagination.Rows = res
	return controller.FormateBody(pagination, http.StatusOK)
}
*/
