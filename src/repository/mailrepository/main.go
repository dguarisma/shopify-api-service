package mailrepository

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository"
	"desarrollosmoyan/lambda/src/utils/set"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type MailRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.MailRepository {
	return &MailRepo{db: db}
}

const ErrSendedPrevius = "El mensaje se envio previamete"

type Producto struct {
	Name string
	Sku  string
	Ean  string

	Count     uint
	BasePrice float32
	Tax       float32
	Discount  float32
	Bonus     uint
	Subtotal  float32
	Total     float32
}

type info struct {
	Fecha    string
	Orden    uint
	Bodega   string
	Name     string
	NIT      string
	Email    string
	Telefono string
	//Productos *[]Producto
	Subtotal  float32
	Impuesto  float32
	Descuento float32
	Total     float32
	To        string
}

type mensajeCompleto struct {
	info
	Products []Producto `json:",omitempty"`
}

func (repo *MailRepo) GetMsg2(purchase *model.Purchase) (mail *model.EmailMsg, err error) {

	exa := mensajeCompleto{
		info: info{
			Name:  "FJM INVERSIONES S.A.S",
			NIT:   "901515179-9",
			Email: "compras.proveedores@farmu.com.co",
		},
	}

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(purchase).Error; err != nil {
			return err
		}
		query := `
		select
		  w.name as Bodega,
		  p.id as Orden,
		  s.phone_contact as Telefono,
		  p.tax as Impuesto,
		  p.discount_global as Descuento,
		  p.sub_total as Subtotal,
		  p.total as Total,
		  s.email_contact as 'To'
		from purchases p
		inner join warehouses w on w.id = p.warehouse_id
		inner join suppliers  s on s.id = p.supplier_id
		where p.id = ?`

		if err := tx.Raw(query, purchase.ID).
			Scan(&mail).Error; err != nil {
			return err
		}

		query2 := `
		select
		  pro.name,
		  pro.sku,
		  pro.ean,
		  a.count,
		  a.base_price,
		  a.tax,
		  a.discount,
		  a.bonus,
		  a.sub_total,
		  a.total
		from purchases p
		inner join articles a on p.id = a.purchase_id
		inner join products pro on pro.id = a.product_id
		where p.id = ?`

		if err := tx.Raw(query2, purchase.ID).
			Scan(&exa.Products).Error; err != nil {
			return err
		}

		if j, _ := json.MarshalIndent(exa, "", "\t"); true {
			fmt.Printf("\n\n%v\n\n", string(j))
		}
		return nil
	})
	return mail, err
}

func (repo *MailRepo) GetMsg(purchase *model.Purchase) (mail *model.EmailMsg, err error) {
	msg := model.EmailMsg{}

	err = repo.db.Transaction(func(tx *gorm.DB) error {
		purchase2 := &model.Purchase{}
		warehouse := &model.Warehouse{}
		supplier := &model.Supplier{}

		if purchase.ID != 0 {

			if err := tx.Preload("Articles").
				First(purchase2, purchase.ID).Error; err != nil {
				return fmt.Errorf(
					"Compra id:%d tubo un problema al ser enviada: %v",
					purchase2.ID, err.Error(),
				)
			}
		} else {
			if err := tx.Save(purchase).Error; err != nil {
				return err
			}
		}
		if purchase2.Status == 1 {
			return fmt.Errorf(ErrSendedPrevius)
		}

		warehouse.ID = purchase.WarehouseID
		supplier.ID = purchase.SupplierID

		res := tx.First(&warehouse)
		if err := handlerErr(res); err != nil {
			return fmt.Errorf(
				"Compra id: %d tubo un problema al obtener la bodega %d ser enviada: %v",
				purchase2.ID, purchase.WarehouseID, err.Error(),
			)
		}

		res = tx.First(&supplier)
		if err := handlerErr(res); err != nil {
			return fmt.Errorf(
				"Compra id: %d tubo un problema al obtener el proveedor %d ser enviada: %v",
				purchase2.ID, purchase.SupplierID, err.Error(),
			)
		}

		productsIds := set.New[uint]()
		for _, article := range purchase.Articles {
			productsIds.Add(article.ProductID)
		}

		products := []model.Product{}
		res = tx.Find(&products, productsIds.Get())
		if err := handlerErr(res); err != nil {
			return fmt.Errorf(
				"Compra id: %d tubo un problema al obtener los productos %d ser enviada: %v",
				purchase2.ID, purchase.SupplierID, err.Error(),
			)
		}

		msg = model.EmailMsg{
			Bodega: warehouse.Name,
			Orden:  purchase.ID,
			// Orden:  supplier.ID,
			Fecha:     supplier.CreatedAt.Format("02/01/2006"),
			Name:      "FJM INVERSIONES S.A.S",
			NIT:       "901515179-9",
			Email:     "compras.proveedores@farmu.com.co",
			Telefono:  supplier.PhoneContact,
			Impuesto:  purchase.Tax,
			Descuento: purchase.DiscountGlobal,
			Productos: make([]model.Producto, 0, len(purchase.Articles)),

			Subtotal: purchase.SubTotal,
			Total:    purchase.Total,
			To:       supplier.EmailContact,
		}

		for _, art := range purchase.Articles {
			for _, product := range products {
				if art.ProductID != product.ID {
					continue
				}

				msg.Productos = append(msg.Productos, model.Producto{
					Name:      product.Name,
					Sku:       product.Sku,
					Ean:       product.Ean,
					Count:     art.Count,
					BasePrice: art.BasePrice,
					Tax:       art.Tax,
					Discount:  art.Discount,
					Bonus:     art.Bonus,
					Subtotal:  art.SubTotal,
					Total:     art.Total,
				})
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func handlerErr(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return fmt.Errorf("Elemento no existente")
	}
	return nil
}
