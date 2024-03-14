package receptions

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/shopifyserv"
	"fmt"

	"gorm.io/gorm"
)

func New(db *gorm.DB, shopifyserv shopifyserv.ShopifyService) *HandlerPurchaseStatus {
	return &HandlerPurchaseStatus{
		db:          db,
		shopifyserv: shopifyserv,
	}
}

type IHandlerPurchaseStatus interface {
	InsertOne(tx *gorm.DB, reception *model.ReceptionArt) (*model.ReceptionArt, error)
	InsertMany(tx *gorm.DB, reception []model.ReceptionArt) ([]model.ReceptionArt, error)
	IsArticleOverflow(tx *gorm.DB, articleId uint) (err error)
	SendUpdateShopify(tx *gorm.DB, reception *model.ReceptionArt) (err error)
	SendUpdatesShopify(tx *gorm.DB, reception []model.ReceptionArt) (err error)
	SendUpdateDeleteShopify(tx *gorm.DB, reception *model.ReceptionArt) (err error)
}

type HandlerPurchaseStatus struct {
	db          *gorm.DB
	shopifyserv shopifyserv.ShopifyService
}

func (handler *HandlerPurchaseStatus) InsertOne(tx *gorm.DB, reception *model.ReceptionArt) (*model.ReceptionArt, error) {
	if err := tx.Save(&reception).Error; err != nil {
		return nil, err
	}
	return reception, nil
}

func (handler *HandlerPurchaseStatus) InsertMany(tx *gorm.DB, receptions []model.ReceptionArt) ([]model.ReceptionArt, error) {
	if err := tx.Save(&receptions).Error; err != nil {
		return nil, err
	}
	return receptions, nil
}

func (handler *HandlerPurchaseStatus) IsArticleOverflow(tx *gorm.DB, articleId uint) (err error) {
	var articlesCompletes uint

	article, err := getById[model.Article](tx, articleId, "Article")
	if err != nil {
		return err
	}

	if err = tx.Model(&model.ReceptionArt{}).
		Where(&model.ReceptionArt{ArticleID: articleId}).
		Select("sum(count)").
		Find(&articlesCompletes).Error; err != nil {

		return err
	}

	if article.Count < articlesCompletes {
		return fmt.Errorf(
			"error al intentar ingresar en el articulo(%v) la cantidad de %v items cuado deberia ser %v",
			articleId, articlesCompletes, article.Count,
		)
	}
	purchase, err := getById[model.Purchase](tx, article.PurchaseID, "purchase")
	if err != nil {
		return err
	}

	if article.Count == articlesCompletes {
		purchase.ReceptionStatus = 2
	}
	if article.Count > articlesCompletes {
		purchase.ReceptionStatus = 1
	}
	if articlesCompletes == 0 {
		purchase.ReceptionStatus = 0
	}

	if err := tx.Save(&purchase).Error; err != nil {
		return fmt.Errorf("error al actualiza el estado de la compra")
	}

	return nil
}

func (handler *HandlerPurchaseStatus) SendUpdateShopify(tx *gorm.DB, reception *model.ReceptionArt) (err error) {
	shopiUpdate, err := handler.getCityAndProductShopifyId(tx, reception.ArticleID)
	if err != nil {
		return err
	}
	shopiUpdate.Items[0].Available = int64(reception.Count)

	return handler.shopifyserv.SendUpdateInventory(*shopiUpdate)
}

func (handler *HandlerPurchaseStatus) SendUpdatesShopify(tx *gorm.DB, receptions []model.ReceptionArt) (err error) {
	shopiUpdate, err := handler.getCityAndProductShopifyId(tx, receptions[0].ArticleID)
	if err != nil {
		return err
	}

	for _, reception := range receptions {
		shopiUpdate.Items[0].Available += int64(reception.Count)
	}

	return handler.shopifyserv.SendUpdateInventory(*shopiUpdate)
}

func (handler *HandlerPurchaseStatus) getCityAndProductShopifyId(tx *gorm.DB, articleId uint) (*shopifyserv.ProductsUpdates, error) {
	/*
		query := `
		select pur.id
		from articles a
		inner join purchases pur on pur.id = a.purchase_id
		inner join warehouses w  on w.id   = pur.warehouse_id
		inner join products p    on p.id   = a.product_id
		where a.id = ?


		`
	*/
	ware := &model.Warehouse{}
	product := &model.Product{}

	purchaseId := tx.Model(&model.Article{}).
		Select((&model.Article{}).GetPurchaseId()).
		Where("id = ?", articleId)

	warehouseId := tx.Model(&model.Purchase{}).
		Select((&model.Purchase{}).GetWarehouseId()).
		Where("id = (?)", purchaseId)

	productId := tx.Model(&model.Article{}).
		Select((&model.Article{}).GetProductId()).
		Where("id = ?", articleId)

	result := tx.Find(ware, "id = (?)", warehouseId)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no se ha encontrado la bodega perteneciente al articulo")
	}

	if err := tx.Find(product, "id = (?)", productId).Error; err != nil {
		return nil, err
	}

	shopiUpdate := shopifyserv.ProductsUpdates{
		Items:     []shopifyserv.Item{{}},
		WareHouse: ware.City,
	}

	if shopiUpdate.WareHouse == "Barranquilla" {
		shopiUpdate.Items[0].ProductShopifyID = product.HandlesBaq
	}

	if shopiUpdate.WareHouse == "Bogotá" {
		shopiUpdate.Items[0].ProductShopifyID = product.HandlesBog
	}

	return &shopiUpdate, nil
}

func (handler *HandlerPurchaseStatus) SendUpdateDeleteShopify(tx *gorm.DB, reception *model.ReceptionArt) (err error) {
	shopiUpdate, err := handler.getCityAndProductShopifyId(tx, reception.ArticleID)
	if err != nil {
		return err
	}
	shopiUpdate.Items[0].Available = -1 * int64(reception.Count)
	return handler.shopifyserv.SendUpdateInventory(*shopiUpdate)
}

func getById[item any](tx *gorm.DB, id uint, typeOfItem string) (i *item, err error) {
	result := tx.Find(&i, id)

	if result.Error != nil {
		return i, fmt.Errorf(
			"error al buscar el %v con el id(%v): %v",
			typeOfItem, id, result.Error.Error(),
		)
	}

	if result.RowsAffected == 0 {
		return i, fmt.Errorf(
			"error al buscar (%v): el %v no existe",
			id, typeOfItem,
		)
	}
	return i, nil
}
