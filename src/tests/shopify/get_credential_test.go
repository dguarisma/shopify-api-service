package shopify_test

import (
	"desarrollosmoyan/lambda/src/shopifyserv"
	"desarrollosmoyan/lambda/src/tests/utils"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredential(t *testing.T) {
	_, err := utils.NewHandleTest()
	// lo uso porque carga las variables de entorno
	assert.NoError(t, err, "handleTest don't return error")

	t.Run("Error", func(t *testing.T) {

		type caseErr struct {
			Input    string
			Expected string
		}
		examples := []caseErr{
			{"", shopifyserv.ErrEmptyCredentials},
			{"example,", shopifyserv.ErrNotEnoughtCredentials},
		}

		for _, curCase := range examples {

			_, err := shopifyserv.GetCredential(curCase.Input)
			assert.EqualError(t,
				err,
				curCase.Expected,
			)
		}
	})

	t.Run("Success", func(t *testing.T) {

		type caseErr struct {
			Input    string
			Expected shopifyserv.ShopifyCredentials
		}
		examples := []caseErr{
			{"a,b,c", shopifyserv.ShopifyCredentials{Name: "a", UrlBase: "b", XShopifyAccessToken: "c"}},
			{"a,c,d", shopifyserv.ShopifyCredentials{Name: "a", UrlBase: "c", XShopifyAccessToken: "d"}},
		}
		for _, str := range []string{os.Getenv("BARRANQULLA"), os.Getenv("BOGOTA")} {
			cred := strings.Split(str, ",")
			credential := shopifyserv.ShopifyCredentials{
				Name:                cred[0],
				UrlBase:             cred[1],
				XShopifyAccessToken: cred[2],
			}
			examples = append(examples, caseErr{str, credential})
		}

		for _, curCase := range examples {
			output, err := shopifyserv.GetCredential(curCase.Input)
			assert.NoError(t, err)
			assert.Equal(t,
				curCase.Expected.Name,
				output.Name,
			)

		}
	})

}
