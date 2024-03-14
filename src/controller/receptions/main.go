package receptions

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/repository/genericrepository"
	"desarrollosmoyan/lambda/src/repository/receptionrepository"
	"desarrollosmoyan/lambda/src/response"
	"desarrollosmoyan/lambda/src/shopifyserv"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func NewReceptionServ(db *gorm.DB, store shopifyserv.ShopifyService) controller.CRUD {
	return &ReceptionServ{
		repo: receptionrepository.New(
			db,
			genericrepository.New[*model.ReceptionArt](db, []string{}),
		),
		db:                    db,
		handlerPurchaseStatus: New(db, store),
	}
}

type ReceptionServ struct {
	repo                  repository.ReceptionRepository
	handlerPurchaseStatus IHandlerPurchaseStatus
	db                    *gorm.DB
}

func (r *ReceptionServ) GetAll() response.Result           { return r.repo.GetAll() }
func (r *ReceptionServ) GetByID(id string) response.Result { return r.repo.GetByID(id) }

func (r *ReceptionServ) GetManyBy(pathReq []byte) response.Result {
	type ForArticleId struct {
		ArticleID uint `json:"ArticleID,string"`
	}

	temp2 := ForArticleId{}
	if err := json.Unmarshal(pathReq, &temp2); err != nil {
		return response.NewFailResult(
			fmt.Errorf("parametros en la peticion no aceptado, solo esta permidito[ArticleID]"),
			http.StatusBadRequest,
		)
	}

	if temp2.ArticleID == 0 {
		return response.NewFailResult(
			fmt.Errorf("Forma de busqueda no permitida"),
			http.StatusForbidden,
		)
	}
	search := &model.ReceptionArt{ArticleID: temp2.ArticleID}
	return r.repo.GetManyBy(search)
}

func (r *ReceptionServ) InsertOne(reqBody []byte) response.Result {

	element, err := handleNewElement2(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {

		reception, err := r.handlerPurchaseStatus.InsertOne(tx, &element)
		if err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.IsArticleOverflow(tx, reception.ArticleID); err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.SendUpdateShopify(tx, reception); err != nil {
			// fallo la actualizacion en shopify
			return err
		}

		return nil
	})

	if err != nil {
		return response.NewFailResult(err, 500)
	}

	return controller.FormateBody(element, http.StatusOK)
}

func (r *ReceptionServ) InsertMany(reqBody []byte) response.Result {

	elements, err := handleNewElements2(reqBody)
	if err != nil {
		return response.NewFailResult(
			fmt.Errorf("error:%v", err.Error()), http.StatusBadRequest,
		)
	}

	receptions := []model.ReceptionArt{}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		receptions, err := r.handlerPurchaseStatus.InsertMany(tx, elements)
		if err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.IsArticleOverflow(tx, receptions[0].ArticleID); err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.SendUpdatesShopify(tx, receptions); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.NewFailResult(err, 500)
	}

	return controller.FormateBody(receptions, http.StatusOK)
}

func (p *ReceptionServ) Update(reqBody []byte) response.Result {
	return response.NewFailResult(fmt.Errorf("No se pueden actualizar"), http.StatusBadRequest)
}

func (r *ReceptionServ) DeleteById(reqBody []byte) response.Result {
	id, err := controller.HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)

	}
	reception := &model.ReceptionArt{}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Find(reception, id)
		articleId := reception.ArticleID
		if result.Error != nil {
			return fmt.Errorf("No se pudo encontrar la reception que se intenta borrar: %v", err.Error())
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("No se pudo encontrar la reception que se intenta borrar")
		}
		if err := tx.Delete(&reception, reception.ID).Error; err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.IsArticleOverflow(tx, articleId); err != nil {
			return err
		}

		if err := r.handlerPurchaseStatus.SendUpdateDeleteShopify(r.db, reception); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return response.NewFailResult(err, 500)
	}
	return controller.FormateBody(reception, 200)
}
