package product_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"desarrollosmoyan/lambda/src/tests/utils"
	"desarrollosmoyan/lambda/src/utils/set"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// / //go:embed products.json
//
//	//go:embed example.json
var (

	//go:embed products3.json
	productsArr []byte

// ErrSendedPrevius = "El mensaje se envio previamete"
)

type Agregar interface {
	Add(data uint)
}

func stringUint(set Agregar, value string, typekOf string) {
	if value != "" {
		id, err := strconv.ParseUint(value, 10, 63)
		if err != nil {
			log.Panicf("%v %v", typekOf, value)
		}
		set.Add(uint(id))
	}
}

func StringToArrUint(set Agregar, value string, typekOf string) {
	if value == "" {
		return
	}

	dataWithoutSpace := regexp.
		MustCompile(" ").
		ReplaceAllLiteralString(value, "")

	info := strings.Split(dataWithoutSpace, ",")

	for _, data := range info {
		stringUint(set, data, typekOf)
	}
}

func TestInsertMany(t *testing.T) {

	handleTest, err := utils.NewHandleTest()
	if err != nil {
		t.Fatalf(err.Error())
	}
	/*
		tables := []interface{}{
			&model.SubstanceProducts{},
			&model.WarehouseProducts{},
			&model.Product{},
			&model.Substance{},
			&model.Pack{},
			&model.Trademark{},
			&model.Maker{},
			&model.TypesProduct{},
			&model.Warehouse{},
		}

		handleTest.DeleteInfo(t, tables)
		defer handleTest.DeleteInfo(t, tables)
	*/

	products := []Product{}
	if err := json.Unmarshal(productsArr, &products); err != nil {
		log.Println(err.Error())
	}

	makersId := set.New[uint]()
	trademakersId := set.New[uint]()
	TypesProductsID := set.New[uint]()
	// CategoryOneID := set.New[uint]()
	// CategoryTwoID := set.New[uint]()
	// CategoryThreeID := set.New[uint]()
	PackID := set.New[uint]()
	//SubstancesIDS := set.New[uint]()
	WarehouseIDS := set.New[uint]()
	// SubstitutesIDS := set.New[uint]()

	for i, product := range products {

		products[i].CategoryOneID = ""
		products[i].CategoryTwoID = ""
		products[i].CategoryThreeID = ""
		products[i].PackID = ""
		products[i].SubstitutesIDS = ""

		stringUint(makersId, product.MakerID, "makeId")
		stringUint(trademakersId, product.TrademarkID, "tradeMark")
		stringUint(TypesProductsID, product.TypesProductID, "TypesProductID")
		// stringUint(CategoryOneID, product.CategoryOneID, "CategoryOneID")
		// stringUint(CategoryTwoID, product.CategoryTwoID, "CategoryTwoID")
		// stringUint(CategoryThreeID, product.CategoryThreeID, "CategoryThreeID")
		stringUint(PackID, product.PackID, "PackID")
		//StringToArrUint(SubstancesIDS, product.SubstancesIDS, "SubstancesIDS")
		StringToArrUint(WarehouseIDS, product.WarehouseIDS, "WarehouseIDS")
		// StringToArrUint(SubstitutesIDS, product.SubstitutesIDS, "SubstitutesIDS")
	}
	products = append(products, Product{})

	makers := make([]model.Maker, len(makersId.Get()))
	for i, id := range makersId.Get() {
		makers[i].ID = id
	}
	handleTest.Db.Save(&makers)

	tMarks := make([]model.Trademark, len(trademakersId.Get()))
	for i, id := range trademakersId.Get() {
		tMarks[i].ID = id
		tMarks[i].MakerID = makers[0].ID
	}
	handleTest.Db.Save(&tMarks)

	packs := make([]model.Pack, len(PackID.Get()))
	for i, id := range PackID.Get() {
		packs[i].ID = id
	}
	handleTest.Db.Save(&packs)
	tProducts := make([]model.TypesProduct, len(TypesProductsID.Get()))
	for i, id := range TypesProductsID.Get() {
		tProducts[i].ID = id
	}
	handleTest.Db.Save(&tProducts)

	/*
		whouse := make([]model.Warehouse, len(WarehouseIDS.Get()))
		for i, id := range WarehouseIDS.Get() {
			whouse[i].ID = id
		}
		handleTest.Db.Save(&whouse)
			sustancias := make([]model.Substance, len(SubstancesIDS.Get()))
			for i, id := range SubstancesIDS.Get() {
				sustancias[i].ID = id
			}
			handleTest.Db.Save(&sustancias)
	*/

	b, _ := json.Marshal(&products)

	//fmt.Println(string(b))
	request := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		//	HTTPMethod: http.MethodPut,
		Body: string(b),
	}

	resBody := handleTest.UseHandleRequest(t,
		server.Product,
		request,
		http.StatusOK,
	)
	//	fmt.Println(resBody)

	request = events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
	}

	resBody = handleTest.UseHandleRequest(t,
		server.Product,
		request,
		http.StatusOK,
	)

	queryParams := map[string]string{
		//"ID":    fmt.Sprint(example.ID),
		"Sku":   "CJ0ZGP9W",
		"page":  "1",
		"limit": "2",
	}
	request = events.APIGatewayProxyRequest{
		HTTPMethod:            http.MethodGet,
		QueryStringParameters: queryParams,
	}

	resBody = handleTest.UseHandleRequest(t,
		server.Product,
		request,
		http.StatusOK,
	)
	func(i ...interface{}) {

	}(resBody)
	fmt.Println(resBody)
}

type Product struct {
	Name            string `json:"Name"`
	Sku             string `json:"Sku"`
	Ean             string `json:"Ean"`
	MakerID         string `json:"MakerID"`
	TrademarkID     string `json:"TrademarkID"`
	TypesProductID  string `json:"TypesProductID"`
	Variation       string `json:"Variation"`
	CategoryOneID   string `json:"CategoryOneID"`
	CategoryTwoID   string `json:"CategoryTwoID"`
	CategoryThreeID string `json:"CategoryThreeID"`
	PackID          string `json:"PackID"`
	Quantity        string `json:"Quantity"`
	Weight          string `json:"Weight"`
	Width           string `json:"Width"`
	Height          string `json:"Height"`
	Depth           string `json:"Depth"`
	Keywords        string `json:"Keywords"`
	SubstancesIDS   string `json:"SubstancesIDS"`
	WarehouseIDS    string `json:"WarehouseIDS"`
	SubstitutesIDS  string `json:"SubstitutesIDS"`
	Status          bool   `json:"Status"`
	HandlesBaq      string `json:"HandlesBaq"`
	HandlesBog      string `json:"HandlesBog"`
	Iva             string `json:"iva"`
	Taxed           bool   `json:"Taxed"`
}
