package purchases

import (
	"desarrollosmoyan/lambda/src/repository/purchaserepository"
	"desarrollosmoyan/lambda/src/tests/utils"
	"fmt"
	"testing"
)

const endpoint = "/pack"

var (
	handleTest *utils.HandleTest
)

func init() {
	handleTest, _ = utils.NewHandleTest()
	// handleTest.Show()
}

func TestGet(t *testing.T) {
	tx := handleTest.Begin()
	defer handleTest.Rollback()
	repo := purchaserepository.New(tx, nil)
	aa := repo.GetAll()
	fmt.Println(aa)

	//mapa := []map[string]interface{}{}
	/*
		fmt.Print("\n\n\n")
		mapa := model.PurchaseForGet{}
		q := `
			SELECT *
			FROM purchases p
			WHERE p.deleted_at IS NULL Limit 1
			`
		tx.Raw(q).Scan(&mapa)

		if j, _ := json.MarshalIndent(mapa.Purchase, "", "\t"); true {
			fmt.Printf("\n\n%v\n\n", string(j))
		}
	*/

}
