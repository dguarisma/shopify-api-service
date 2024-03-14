package purchases

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/mailserv"
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/repository/purchaserepository"
	"desarrollosmoyan/lambda/src/response"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func NewPurchaseRepository(db *gorm.DB, mailservice mailserv.MailService) controller.CRUD {
	return &PurchaseServ{
		repo: purchaserepository.New(
			db, []string{"Articles"},
		),
		mailServ: mailservice,
		Db:       db,
	}
}

type PurchaseServ struct {
	repo     repository.PurchaseRepository
	mailServ mailserv.MailService
	// storeShopify shopifyserv.ShopifyService
	Db *gorm.DB
}

// Get -----------------------------------------------------------
func (p *PurchaseServ) GetAll() response.Result           { return p.repo.GetAll() }
func (p *PurchaseServ) GetByID(id string) response.Result { return p.repo.GetByID(id) }

func (p *PurchaseServ) GetManyBy(pathMap []byte) response.Result {
	search := &model.Purchase{}
	if err := json.Unmarshal(pathMap, search); err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.GetManyBy(search)
}

// Post -----------------------------------------------------------
func (p *PurchaseServ) InsertOne(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if element.Status == 1 {
		p.mailServ.HandleMsg(&element)
	}
	return p.repo.Insert(&element)
}

func (p *PurchaseServ) InsertMany(reqBody []byte) response.Result {
	elements, err := HandleNewElements(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.InsertMany(elements)
}

// Update ---------------------------------------------------------
func (p *PurchaseServ) Update(reqBody []byte) response.Result {
	element, err := HandleNewElement(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	if element.Status == 1 {
		p.mailServ.HandleMsg(&element)
	}
	return p.repo.Update(&element)
}

// Delete ---------------------------------------------------------
func (p *PurchaseServ) DeleteById(reqBody []byte) response.Result {
	purchase, err := HandleDelete(reqBody)
	if err != nil {
		return response.NewFailResult(err, http.StatusBadRequest)
	}
	return p.repo.DeleteByID(purchase)
}

func HandleDelete(body []byte) (*model.Purchase, error) {
	adapt := &model.Purchase{}
	if err := json.Unmarshal(body, adapt); err != nil {
		return nil, controller.ErrFormatForID
	}
	return adapt, nil
}
