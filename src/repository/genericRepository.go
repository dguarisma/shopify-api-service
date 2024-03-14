package repository

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/response"
)

type GenericRepository[k model.ICustom] interface {
	GetAll() response.Result
	GetByID(id string) response.Result
	GetManyBy(req k) response.Result

	// Pagination
	GetFirstPage() response.Result
	GetByPagination(pathMap []byte) response.Result

	Insert(res k) response.Result

	// este metodo inserta nuevos elementos
	// y actualiza otros pero borra los que no esten
	// dentro de la actualizacion
	InsertMany(updates []k) error

	// este metodo inserta nuevos elementos y actualiza otros
	InsertsAndUpdates(updates []k) error

	Update(res k) response.Result
	// UpdateFull(res k) response.Result

	DeleteByID(id string) response.Result
}
