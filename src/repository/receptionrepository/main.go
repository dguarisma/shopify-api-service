package receptionrepository

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/response"
	"net/http"

	"gorm.io/gorm"
)

type ReceptionRepo struct {
	Db *gorm.DB
	//grepo repository.GenericRepository[*model.ReceptionArt]
	repository.GenericRepository[*model.ReceptionArt]
}

func New(db *gorm.DB, grepo repository.GenericRepository[*model.ReceptionArt]) repository.ReceptionRepository {
	return &ReceptionRepo{
		Db:                db,
		GenericRepository: grepo,
		//grepo: grepo,
	}
}

func (repo *ReceptionRepo) Insert(element *model.ReceptionArt) response.Result {
	tx := repo.Db.Begin()
	if err := tx.Save(element).Error; err != nil {
		tx.Rollback()
		return response.NewFailResult(err, 500)
	}
	if err := repo.IsArticleOverflow(tx, element.ArticleID); err != nil {
		tx.Rollback()
		return response.NewFailResult(err, 500) // cambia este error
	}

	return controller.FormateBody(element, http.StatusOK)
}

func (repo *ReceptionRepo) IsArticleOverflow(tx *gorm.DB, articleId uint) error {
	count := Count{}
	query := `
	select
	  sum(r.count) as Reception,
	  a.count as Articles
	from articles a
	inner join reception_arts r on r.article_id = a.id
	inner join purchases p on p.id = a.purchase_id
	where a.id = ?` // agregar que es posible que tenga null
	if err := tx.Raw(query, articleId).Scan(&count).Error; err != nil {
		return err
	}

	if count.Articles < count.Reception {
		return ErrOverflowArticles(articleId, count.Articles, count.Reception)
	}
	return nil
}

func (repo *ReceptionRepo) TransactionInsert(element *model.ReceptionArt) (tx *gorm.DB, err error) {
	tx = repo.Db.Begin()
	if err := tx.Save(element).Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (repo *ReceptionRepo) InsertMany(elements []*model.ReceptionArt) error {
	updateElement := make([]model.ReceptionArt, 0, len(elements))
	newElement := make([]model.ReceptionArt, 0, len(elements))

	for _, element := range elements {
		if element.ID != 0 {
			updateElement = append(updateElement, *element)
			continue
		}
		newElement = append(newElement, *element)
	}

	err := repo.Db.Transaction(func(tx *gorm.DB) error {
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
		return err

	}

	//res := append(updateElement, newElement...)
	return nil
}

/*

func (repo *ReceptionRepo) GetAll() response.Result           { return repo.grepo.GetAll() }
func (repo *ReceptionRepo) GetByID(id string) response.Result { return repo.grepo.GetByID(id) }
func (repo *ReceptionRepo) GetManyBy(search *model.ReceptionArt) response.Result {
	return repo.grepo.GetManyBy(search)
}
func (repo *ReceptionRepo) Update(res *model.ReceptionArt) response.Result {
	if err := repo.Db.Save(res).Error; err != nil {
		return response.NewFailResult(
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *ReceptionRepo) DeleteByID(id string) response.Result {
	//return repo.grepo.DeleteByID(id)
	res := &model.ReceptionArt{}
	result := repo.Db.Delete(res, id)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK)
}
*/
