package inventory_test

import (
	"desarrollosmoyan/lambda/src/controller/inventory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortTables(t *testing.T) {
	handlerInvent := inventory.HandlerTables{}
	t.Run("Empty list", func(t *testing.T) {
		list, err := handlerInvent.SortTables(nil, nil)
		assert.NotNil(t, err)
		assert.Equal(t, inventory.ErrNotSellOrPurchases, err.Error())
		assert.Nil(t, list)
	})

	t.Run("One Item", func(t *testing.T) {
		purchase := []inventory.InventorySell{
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		sell := []inventory.InventorySell{}

		expected := []inventory.InventorySell{
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		list, err := handlerInvent.SortTables(purchase, sell)
		assert.NoError(t, err)
		assert.Equal(t, &expected, list)
	})

	t.Run("Two Items", func(t *testing.T) {
		purchase := []inventory.InventorySell{
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		sell := []inventory.InventorySell{
			{Date: "13/02/2023", Selled: inventory.Store{Units: 10, Cost: 100}},
		}

		expected := []inventory.InventorySell{
			{Date: "13/02/2023", Selled: inventory.Store{Units: 10, Cost: 100}},
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		list, err := handlerInvent.SortTables(purchase, sell)
		assert.NoError(t, err)
		assert.Equal(t, &expected, list)
	})

	t.Run("Three Items and two dates", func(t *testing.T) {
		purchase := []inventory.InventorySell{

			{Date: "13/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		sell := []inventory.InventorySell{
			{Date: "13/02/2023", Selled: inventory.Store{Units: 10, Cost: 100}},
		}

		expected := []inventory.InventorySell{
			{
				Date:     "13/02/2023",
				Selled:   inventory.Store{Units: 10, Cost: 100},
				Purchase: inventory.Store{Units: 10, Cost: 100},
			},
			{Date: "03/02/2023", Purchase: inventory.Store{Units: 10, Cost: 100}},
		}

		list, err := handlerInvent.SortTables(purchase, sell)
		assert.NoError(t, err)
		assert.Equal(t, &expected, list)
	})

}
