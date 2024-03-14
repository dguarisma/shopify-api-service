package mailserv_test

import (
	"desarrollosmoyan/lambda/src/mailserv"
	"desarrollosmoyan/lambda/src/model"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic("No se cargaron la variables de entorno")
	}
}

func TestAll(t *testing.T) {
	t.Run("Format_msg_html", func(t *testing.T) {
		// al ejecutarse este test crea un html con el formato del mensaje llamado test.html
		mailServ := mailserv.New(nil, &mailserv.MailerConfig{})

		example := model.EmailMsg{
			Fecha:     time.Now().Format("02/01/2006"),
			Orden:     10,
			Bodega:    "Bodega",
			Name:      "Name",
			NIT:       "NIT",
			Email:     "example@hotmail.com",
			Telefono:  "45831900",
			Subtotal:  90,
			Impuesto:  9,
			Descuento: 1,
			Total:     100,
		}

		for i := 0; i < 10; i++ {
			example.Productos = append(example.Productos, model.Producto{
				Name:      fmt.Sprintf("Name%d", i),
				Sku:       fmt.Sprintf("Sku%d", i),
				Ean:       fmt.Sprintf("Ean%d", i),
				Count:     uint(rand.Intn(100)),
				BasePrice: float32(rand.Int31n(30)),
				Tax:       float32(rand.Int31n(30)),
				Discount:  float32(rand.Int31n(30)),
				Bonus:     uint(rand.Int31n(30)),
				Subtotal:  float32(rand.Int31n(30)),
				Total:     float32(rand.Int31n(100)),
			})
		}

		data, e := mailServ.FormatMsg(&example)
		if e != nil {
			log.Fatalln(e)
		}
		ioutil.WriteFile("test.html", data, 0644)
	})

	t.Run("Send msg", func(t *testing.T) {

		config := &mailserv.MailerConfig{
			Host:     os.Getenv("HOST"),
			Username: os.Getenv("AWS_SES_ACCESS_KEY"),
			Password: os.Getenv("AWS_SES_SECRET_KEY"),
			Sender:   os.Getenv("EMAIL"),
			Port:     587,
			Timeout:  5 * time.Second,
		}

		mailServ := mailserv.New(nil, config)
		mail := os.Getenv("HOST")

		if err := mailServ.Send(mail, "example", "hello world"); err != nil {
			t.Fatal(err.Error())
		}
	})
}
