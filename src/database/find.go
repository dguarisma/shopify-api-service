package database

/*
type Repo struct {
	Db       *gorm.DB
	Preloads []string // son para buscar datos anidados
}

func (repo *Repo) GetAllByPagination(res interface{}, pagination *model.Pagination) response.Result {
	query := repo.Db
	for _, preload := range repo.Preloads {
		query = query.Preload(preload)
	}

	result := query.Scopes(model.Paginate(res, pagination, repo.Db)).Find(res)

	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}
	pagination.Rows = res

	return controller.FormateBody(pagination, http.StatusOK)
}

func (repo *Repo) GetByPagination(pathMap []byte, body interface{}, res interface{}) response.Result {
	pagination, err := model.GetPagination(pathMap)
	if err != nil {
		return controller.ErrPagination(err)
	}
	if err := json.Unmarshal(pathMap, body); err != nil {
		return controller.ErrImposibleFormat()
	}

	query := repo.Db
	for _, preload := range repo.Preloads {
		query = query.Preload(preload)
	}

	result := query.Scopes(model.Paginate(res, pagination, repo.Db)).Find(res, body)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	pagination.Rows = res
	return controller.FormateBody(pagination, http.StatusOK)
}

func (repo *Repo) GetAll(res interface{}) response.Result {
	query := repo.Db
	for _, preload := range repo.Preloads {
		query = query.Preload(preload)
	}

	result := query.Find(res)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}

	return controller.FormateBody(res, http.StatusOK)
}

func (repo *Repo) GetByID(res interface{}, idString string) response.Result {
	query := repo.Db
	for _, preload := range repo.Preloads {
		query = query.Preload(preload)
	}

	result := query.Find(res, idString)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *Repo) GetManyBy(pathMap []byte, model interface{}, res interface{}) response.Result {
	if err := json.Unmarshal(pathMap, model); err != nil {
		return controller.ErrImposibleFormat()
	}

	query := repo.Db
	for _, preload := range repo.Preloads {
		query = query.Preload(preload)
	}

	result := query.Find(res, model)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *Repo) DeleteByID(res interface{}, id string) response.Result {
	result := repo.Db.Delete(res, id)
	if resultErr, problem := controller.HandleResultSearchDB(result); problem {
		return resultErr
	}
	return controller.FormateBody(res, http.StatusOK) // ver si se quiere cambiar el tipo de respuesta
}

func (repo *Repo) Insert(res interface{}) response.Result {
	if err := repo.Db.Save(res).Error; err != nil {
		return response.NewFailResult( // posiblemente necesite otro manejador
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}

func (repo *Repo) InsertMany(newElement, updateElement interface{}) error {
	return repo.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(updateElement).Error; err != nil {
			return err
		}
		if err := tx.Create(newElement).Error; err != nil {
			return err
		}
		return nil
	})
}

func (repo *Repo) Update(res interface{}) response.Result {
	if err := repo.Db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(res).Error; err != nil {
		return response.NewFailResult( // posiblemente necesite otro manejador
			err, http.StatusInternalServerError,
		)
	}
	return controller.FormateBody(res, http.StatusOK)
}
*/
