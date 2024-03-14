package repository

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/response"
)

type Field uint

const (
	Name Field = iota
	Sku
	Ean
)

type ProductRepository interface {
	GenericRepository[*model.Product]
	GetByField(world string, pag model.Pagination, field Field) response.Result
	GetBySkus(skus []string) response.Result
}
