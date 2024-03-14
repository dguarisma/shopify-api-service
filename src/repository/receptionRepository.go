package repository

import (
	"desarrollosmoyan/lambda/src/model"

	"gorm.io/gorm"
)

type ReceptionRepository interface {
	GenericRepository[*model.ReceptionArt]
	// GetAll() response.Result
	// GetByID(id string) response.Result
	// GetManyBy(search *model.ReceptionArt) response.Result
	TransactionInsert(element *model.ReceptionArt) (tx *gorm.DB, err error)
	// Insert(res *model.ReceptionArt) response.Result
	InsertMany(elements []*model.ReceptionArt) error
	//Update(res *model.ReceptionArt) response.Result
	//DeleteByID(id string) response.Result
}
