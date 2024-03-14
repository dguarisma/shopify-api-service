package repository

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/response"
)

type PurchaseRepository interface {
	GetAll() response.Result
	GetByID(id string) response.Result
	GetManyBy(search *model.Purchase) response.Result
	Insert(res *model.Purchase) response.Result
	InsertMany(elements *[]model.Purchase) response.Result
	Update(res *model.Purchase) response.Result
	DeleteByID(res *model.Purchase) response.Result
}
