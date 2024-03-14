package inventory_test

import (
	"desarrollosmoyan/lambda/src/controller/inventory"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetDate(t *testing.T) {
	handleTablesInventory := inventory.HandlerTables{}
	t.Log("este seria el formato de la fecha de todas las tablas de inventory")

	type currentCase struct {
		Input    time.Time
		Expected string
	}

	cases := []currentCase{
		{time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC), "01/11/2023"},
		{time.Date(2023, 02, 3, 0, 0, 0, 0, time.UTC), "03/02/2023"},
		{time.Date(2023, 05, 23, 0, 0, 0, 0, time.UTC), "23/05/2023"},
		{time.Date(2023, 02, 20, 0, 0, 0, 0, time.UTC), "20/02/2023"},
		{time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC), "03/10/2023"},
		{time.Date(2023, 12, 29, 0, 0, 0, 0, time.UTC), "29/12/2023"},
	}
	for i, currentCase := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			output := handleTablesInventory.SetDate(currentCase.Input)
			assert.Equal(t, currentCase.Expected, output)
		})
	}
}
