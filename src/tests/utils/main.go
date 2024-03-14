package utils

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/server"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic("No se cargaron la variables de entorno")
	}
}

type HandleTest struct {
	Db   *gorm.DB
	tx   *gorm.DB
	show bool

	Get    string
	Post   string
	Put    string
	Delete string
}

func NewHandleTest() (*HandleTest, error) {
	//	dbUri := os.Getenv("DB_TEST") // para test
	dbUri := os.Getenv("DB_URI") // para test
	db, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &HandleTest{
		Db:     db,
		tx:     db.Begin(),
		show:   false,
		Get:    "GET",
		Post:   "POST",
		Put:    "PUT",
		Delete: "DELETE",
	}, nil
}
func (ht *HandleTest) Begin() *gorm.DB {
	ht.tx = ht.Db.Begin()
	return ht.tx
}

func (ht *HandleTest) Rollback() { ht.tx.Rollback() }

func (ht *HandleTest) Show() { ht.show = true }

func (ht *HandleTest) SeeCustomModel(custom *model.CustomModel) {
	fmt.Printf("id: %v\n", custom.ID)
	fmt.Printf("created: %v\n", custom.CreatedAt)
	fmt.Printf("update: %v\n", custom.UpdatedAt)
	fmt.Printf("delete: %v\n", custom.DeletedAt)
}

func (ht *HandleTest) DifferentMap(t *testing.T, expected interface{}, got interface{}) {
	if areEqual := reflect.DeepEqual(expected, got); !areEqual {
		t.Errorf("Body response\nexpected: %#v\ngot:      %#v\n\n", expected, got)
	}
}

func (ht *HandleTest) UseHandleRequest(t *testing.T, method uint8, request events.APIGatewayProxyRequest, statusExpected int) string {
	res, err := server.HandleRequest2(ht.tx, method, nil)(request)
	if err != nil {
		t.Errorf("la peticion nunca debe devolver error: %s", err.Error())
		return ""
	}
	if statusExpected != res.StatusCode {
		msg := "\nStatus"
		msg += "\n\tExpected: %d"
		msg += "\n\tGot: %d\n"
		t.Errorf(msg, statusExpected, res.StatusCode)
	}
	return res.Body
}

func (ht *HandleTest) UseHandleRequestTx(t *testing.T, method uint8, request events.APIGatewayProxyRequest, statusExpected int) string {
	res, err := server.HandleRequest2(ht.tx, method, nil)(request)
	if err != nil {
		t.Errorf("la peticion nunca debe devolver error: %s", err.Error())
		return ""
	}

	if statusExpected != res.StatusCode {
		msg := "\nStatus"
		msg += "\n\tExpected: %d"
		msg += "\n\tGot: %d\n"
		t.Errorf(msg, statusExpected, res.StatusCode)
	}
	return res.Body
}

func (ht *HandleTest) ShowRequest(t *testing.T, endpoint, method string, body []byte, response string) {
	if !ht.show {
		return
	}
	msg := "\nmethod: %s"
	msg += "\nendpoint: %s"
	msg += "\nbody: \n%s"
	msg += "\nresponse:\n%s"
	t.Logf(msg, method, endpoint, string(body), response)
}

func (ht *HandleTest) ShowGetRequest(t *testing.T, endpoint, response string) {
	if !ht.show {
		return
	}
	msg := "\nmethod: %s"
	msg += "\nendpoint: %s"
	msg += "\nresponse:\n%s"
	t.Logf(msg, http.MethodGet, endpoint, response)
}

func (ht *HandleTest) ShowGetByRequest(t *testing.T, endpoint, response string, mapa map[string]string) {
	if !ht.show {
		return
	}
	msg := "\nmethod: %s"
	msg += "\nendpoint: %s"
	for key, value := range mapa {
		msg += fmt.Sprintf("?%v=%v", key, value)
	}

	msg += "\nresponse:\n%s"
	t.Logf(msg, http.MethodGet, endpoint, response)
}

func (ht HandleTest) ShowRequestEmpty(t *testing.T, method, response string) {
	t.Logf("\nmethod: %s\nresponse:\n%s", method, response)
}

type DeleteById struct{ ID uint }
