package purchaserepository

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/response"
	"desarrollosmoyan/lambda/src/utils/set"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func New(db *gorm.DB, preload []string) repository.PurchaseRepository {
	return &PurchaseRepo{
		db:       db,
		preloads: preload,
	}
}

type PurchaseRepo struct {
	db       *gorm.DB
	preloads []string
}

func (repo *PurchaseRepo) GetAll() response.Result {
	res := []model.PurchaseForGet{}
	query := repo.preload()
	//	pag := model.Pagination{Limit: 100}
	//result := query.Scopes(model.Paginate(&[]model.Purchase{}, &pag, query)).Find(&res)
	result := query.Model(&[]model.Purchase{}).Find(&res)
	//result := query.Scopes(model.Paginate(&[]model.Purchase{}, &pag, query)).Find(&res)

	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	result = query.Model(&[]model.Purchase{}).Find(&res)
	if resultErr, problem := controller.
		HandleResultSearchDB(result); problem {
		return resultErr
	}
	suppliersIds := set.New[uint](len(res))
	warehouseIds := set.New[uint](len(res))

	for _, purchase := range res {
		suppliersIds.Add(purchase.SupplierID)
		warehouseIds.Add(purchase.WarehouseID)
	}

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		warehouses := []model.Warehouse{}
		suppliers := []model.Supplier{}

		if err := tx.
			Find(&warehouses, warehouseIds.Get()).
			Error; err != nil {
			return err
		}

		if err := tx.
			Find(&suppliers, suppliersIds.Get()).
			Error; err != nil {
			return err
		}

		for _, supp := range suppliers {
			for i, purchase := range res {
				if purchase.SupplierID == supp.ID {
					res[i].Supplier = supp
				}
			}
		}

		for _, wh := range warehouses {
			for i, purchase := range res {
				if purchase.WarehouseID == wh.ID {
					res[i].Warehouse = wh
				}
			}
		}
		return nil
	})

	if err != nil {
		return response.NewFailResult(
			fmt.Errorf("fallo al formar respuesta:%s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) GetByID(id string) response.Result {
	query := repo.preload()
	//res := &model.PurchaseForGet{}
	purcha := &model.Purchase{}

	result := query.First(purcha, id)
	if resultErr, problem := controller.
		HandleResultSearchDB(result); problem {
		return resultErr
	}

	res := &model.PurchaseForGet{Purchase: *purcha}
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		warehouse := model.Warehouse{}
		supplier := model.Supplier{}
		products := []model.Product{}

		if err := tx.
			Find(&warehouse, res.WarehouseID).
			Error; err != nil {
			return err
		}

		if err := tx.
			Find(&supplier, res.SupplierID).
			Error; err != nil {
			return err
		}
		productsId := set.New[uint]()

		for _, art := range purcha.Articles {
			productsId.Add(art.ProductID)
		}

		if err := tx.
			Find(&products, productsId.Get()).
			Error; err != nil {
			return err
		}

		res.Supplier = supplier
		res.Products = products
		res.Warehouse = warehouse
		res.UpdatedAt = res.Purchase.UpdatedAt

		return nil
	})

	if err != nil {
		return response.NewFailResult(
			fmt.Errorf("fallo al formar respuesta:%s", err.Error()),
			http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) GetManyBy(search *model.Purchase) response.Result {
	res := &[]model.Purchase{}

	query := repo.preload()
	result := query.Find(res, search)
	if resultErr, problem := controller.
		HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) DeleteByID(res *model.Purchase) response.Result {
	result := repo.db.Delete(res)
	if resultErr, problem := controller.
		HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) Insert(res *model.Purchase) response.Result {
	if err := repo.
		db.
		Omit("ReceptionStatus").
		Save(res).Error; err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) InsertMany(elements *[]model.Purchase) response.Result {
	updateElement := make([]model.Purchase, 0, len(*elements))
	newElement := make([]model.Purchase, 0, len(*elements))

	for _, element := range *elements {
		if element.ID != 0 {
			updateElement = append(updateElement, element)
			continue
		}
		newElement = append(newElement, element)
	}

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if len(updateElement) != 0 {
			if err := tx.Save(updateElement).Error; err != nil {
				return err
			}
		}
		if len(newElement) != 0 {
			if err := tx.Create(newElement).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}

	res := append(updateElement, newElement...)
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) Update(res *model.Purchase) response.Result {
	deleteActiles, err := repo.searchAndDeleteArticles(res)
	if err != nil {
		return response.NewFailResult(
			fmt.Errorf("Error al comparar para eliminar:%s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		if len(deleteActiles) != 0 {
			if err := tx.
				Delete(&[]model.Article{}, deleteActiles).
				Error; err != nil {
				return err
			}
		}

		if res.ReceptionStatus != 3 {
			if err := tx.
				Session(&gorm.Session{FullSaveAssociations: true}).
				Omit("ReceptionStatus").
				Save(res).Error; err != nil {
				return err
			}
			return nil
		}

		if err := tx.
			Session(&gorm.Session{FullSaveAssociations: true}).
			Save(res).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return response.NewFailResult( // posiblemente necesite otro manejador
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *PurchaseRepo) searchAndDeleteArticles(res *model.Purchase) ([]uint, error) {
	purchaseDb := model.Purchase{}
	if err := repo.db.Preload("Articles").First(&purchaseDb, res.ID).Error; err != nil {
		return nil, err
	}
	articlesForDelete := set.New[uint]()
	for _, art := range purchaseDb.Articles {
		articlesForDelete.Add(art.ID)
	}
	for _, artRQ := range res.Articles {
		articlesForDelete.Delete(artRQ.ID)
	}
	return articlesForDelete.Get(), nil
}

func (repo *PurchaseRepo) preload() *gorm.DB {
	query := repo.db
	for _, preload := range repo.preloads {
		query = query.Preload(preload)
	}
	return query
}
