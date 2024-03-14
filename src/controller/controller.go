package controller

import (
	"desarrollosmoyan/lambda/src/response"
)

type CRUD interface {
	GetAll() response.Result
	GetByID(id string) response.Result
	GetManyBy(pathReq []byte) response.Result

	InsertOne(reqBody []byte) response.Result
	InsertMany(reqBody []byte) response.Result

	Update(reqBody []byte) response.Result
	DeleteById(reqBody []byte) response.Result
}
