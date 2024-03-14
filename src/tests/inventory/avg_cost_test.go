package inventory_test

import (
	"desarrollosmoyan/lambda/src/controller/inventory"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvgCost(t *testing.T) {
	handleTablesInventory := inventory.HandlerTables{}
	t.Log("Esta funcion calcularia el promedio entre el estado del almacen inicial y la compra")
	t.Log("almacen inicial unidades(iu),almacen inicial costo(ic)")
	t.Log("compra unidades(cu),compra costo(cc)")
	t.Log("(cu*cc)+(iu*ic)/(iu+cu)")

	type currentCase struct {
		input    []inventory.Store
		expected float32
	}

	cases := []currentCase{
		{
			input:    []inventory.Store{},
			expected: 0,
		}, {
			input: []inventory.Store{
				{Units: 10, Cost: 1000},
			},
			expected: 1000,
		}, {
			input: []inventory.Store{
				{Units: 10, Cost: 1000},
				{Units: 10, Cost: 1500},
			},
			expected: 1250,
		}, {
			input: []inventory.Store{
				{Units: 10, Cost: 1000},
				{Units: 10, Cost: 1500},
			},
			expected: 1250,
		}, {
			input: []inventory.Store{
				{Units: 974, Cost: 2182},
				{Units: 100, Cost: 2000},
			},
			expected: 2165,
		}, {
			input: []inventory.Store{
				{Units: 1074, Cost: 2165},
				{Units: 200, Cost: 2100},
			},
			expected: 2155,
		}, {
			input: []inventory.Store{
				{Units: 274, Cost: 2155},
				{Units: 300, Cost: 2000},
			},
			expected: 2074,
		}, {
			input: []inventory.Store{
				{Units: 74, Cost: 2074},
				{Units: 400, Cost: 2050},
			},
			expected: 2054,
		}, {
			input: []inventory.Store{
				{Units: 74, Cost: 2054},
				{Units: 0, Cost: 0},
			},
			expected: 2054,
		},
	}

	for i, curCase := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			output := handleTablesInventory.AvgCost(curCase.input...)
			assert.Equal(t, curCase.expected, output)
		})
	}
}
