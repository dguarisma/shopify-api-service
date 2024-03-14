package genericrepository

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/response"
	"desarrollosmoyan/lambda/src/utils/set"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func New[k model.ICustom](db *gorm.DB, preload []string) repository.GenericRepository[k] {
	return &GenericRepo[k]{
		db:       db,
		preloads: preload,
	}
}

type GenericRepo[k model.ICustom] struct {
	db       *gorm.DB
	preloads []string // son para buscar datos anidados
}

func (gr *GenericRepo[k]) GetFirstPage() response.Result {
	var items []k
	pag := model.Pagination{}
	query := gr.db
	for _, preload := range gr.preloads {
		query = query.Preload(preload)
	}

	result := query.Scopes(model.Paginate(items, &pag, gr.db)).Find(&items)

	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	pag.Rows = items
	return controller.FormateBody(pag, http.StatusOK)
}

func (gr *GenericRepo[k]) GetByPagination(pathMap []byte) response.Result {
	var body k
	var res []k

	fmt.Print("2")
	pagination, err := model.GetPagination(pathMap)
	if err != nil {
		return controller.ErrPagination(err)
	}

	fmt.Print("5")
	if err := json.Unmarshal(pathMap, &body); err != nil {
		return response.NewFailResult(err, 500)
	}

	query := gr.db
	for _, preload := range gr.preloads {
		query = query.Preload(preload)
	}

	fmt.Print("6")
	result := query.Scopes(model.Paginate(&res, pagination, gr.db)).Find(&res, &body)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	fmt.Print("7")
	pagination.Rows = res
	return controller.FormateBody(pagination, http.StatusOK)
}

func (gr *GenericRepo[k]) GetAll() response.Result {
	var all []k

	query := gr.db
	for _, preload := range gr.preloads {
		query = query.Preload(preload)
	}

	result := query.Find(&all)
	if result, problem := handlerGetErr(result); problem {
		return result
	}
	return controller.FormateBody(all, http.StatusOK)
}

func (gr *GenericRepo[k]) GetByID(id string) response.Result {
	var item k

	query := gr.db
	for _, preload := range gr.preloads {
		query = query.Preload(preload)
	}

	result := query.Find(&item, id)
	if result, problem := handlerGetErr(result); problem {
		return result
	}
	return controller.FormateBody(item, http.StatusOK)
}

func (gr *GenericRepo[k]) GetManyBy(req k) response.Result {
	res := new(*[]k)
	query := gr.db
	for _, preload := range gr.preloads {
		query = query.Preload(preload)
	}

	result := query.Find(res, req)
	if result, problem := handlerGetErr(result); problem {
		return result
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (gr *GenericRepo[k]) Insert(res k) response.Result {
	if err := gr.db.Save(res).Error; err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (gr *GenericRepo[k]) InsertTransaction(res k) response.Result {
	if err := gr.db.Save(res).Error; err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (gr *GenericRepo[k]) InsertMany(updates []k) error {
	return gr.db.Transaction(func(tx *gorm.DB) error {
		news := make([]k, 0, len(updates))
		update := make([]k, 0, len(updates))
		ref := new(k)
		ids := set.New[uint]()

		dbIds := []struct {
			ID uint `json:"id"`
		}{}

		if err := tx.Model(ref).
			Find(&dbIds).Error; err != nil {
			return err
		}

		for _, curId := range dbIds {
			ids.Add(curId.ID)
		}

		for _, cur := range updates {
			curId := cur.GetId()
			if curId == 0 {
				news = append(news, cur)
				continue
			}
			update = append(update, cur)
			ids.Delete(curId)
		}

		if len(news) != 0 {
			if err := tx.Save(news).Error; err != nil {
				return err
			}
		}

		if len(updates) != 0 {
			if err := tx.Save(updates).Error; err != nil {
				return err
			}
		}
		deletesIds := ids.Get()
		if len(deletesIds) != 0 {
			if err := tx.Delete(ref, deletesIds).
				Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (gr *GenericRepo[k]) InsertsAndUpdates(updates []k) error {
	return gr.db.Transaction(func(tx *gorm.DB) error {
		news := make([]k, 0, len(updates))
		update := make([]k, 0, len(updates))

		for _, cur := range updates {
			if cur.GetId() == 0 {
				news = append(news, cur)
				continue
			}
			update = append(update, cur)
		}

		if len(news) != 0 {
			if err := tx.Save(news).Error; err != nil {
				return err
			}
		}

		if len(updates) != 0 {
			if err := tx.Save(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (gr *GenericRepo[k]) Update(res k) response.Result {
	if err := gr.db.Save(res).Error; err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (gr *GenericRepo[k]) UpdateFull(res k) response.Result {
	if err := gr.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(res).Error; err != nil {
		return response.NewFailResult( // posiblemente necesite otro manejador
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func handlerGetErr(result *gorm.DB) (_ response.Result, problem bool) {
	if result.Error != nil {
		return response.NewFailResult(
			fmt.Errorf("%v", result.Error.Error()),
			http.StatusInternalServerError,
		), true
	}

	if result.RowsAffected == 0 {
		return response.NewFailResult(
			fmt.Errorf("no se ha encontrado ningun elemento"),
			http.StatusInternalServerError,
		), true
	}
	return response.Result{}, false
}

func (gr *GenericRepo[k]) DeleteByID(id string) response.Result {
	item := new(k)
	result := gr.db.Delete(item, id)
	if resultErr, problem := handlerGetErr(result); problem {
		return resultErr
	}
	return controller.FormateBody(item, http.StatusOK)
}
