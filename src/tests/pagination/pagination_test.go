package main

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/tests/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")
	count := 10000
	packs := make([]model.Pack, count)
	status := false
	for i := range packs {
		packs[i].Name = fmt.Sprintf("pack-%v", i)
		packs[i].Status = status
		status = !status
	}

	tx := handleTest.Begin()
	defer handleTest.Rollback()

	err = tx.CreateInBatches(&packs, len(packs)).Error
	assert.NoError(t, err)

	t.Run("Find True", func(t *testing.T) {
		packsResult := []model.Pack{}
		pag := model.Pagination{Limit: 10}
		where := tx.Model(&packsResult).Where("status = ?", true)
		pag2, err := model.PaginateBy(&packsResult, where, &pag)
		assert.NoError(t, err)

		assert.Equal(t, pag2.TotalPages, count/2/pag.Limit)

		pag = model.Pagination{Page: 500}
		where = tx.Model(&packsResult).Where("status = ?", true)
		pag2, err = model.PaginateBy(&packsResult, where, &pag)
		assert.NoError(t, err)
	})
	t.Run("Find All", func(t *testing.T) {
		packsResult := []model.Pack{}
		where := tx.Model(&packsResult)

		t.Run("first page", func(t *testing.T) {
			pag := &model.Pagination{Limit: 10}
			pag2, err := model.PaginateBy(&packsResult, where, pag)
			assert.NoError(t, err)

			assert.Equal(t, pag2.TotalPages, count/pag2.Limit)
			assert.Equal(t, pag2.TotalRows, int64(count))
			assert.Equal(t, pag2.Page, 1)
			rows, ok := (pag2.Rows).([]model.Pack)
			if ok {
				assert.Equal(t, len(rows), pag2.Limit)
			}
		})

		t.Run("last page", func(t *testing.T) {
			pag := &model.Pagination{Limit: 10, Page: count / 10}
			pag2, err := model.PaginateBy(&packsResult, where, pag)
			assert.NoError(t, err)

			assert.Equal(t, pag2.TotalPages, count/pag2.Limit)
			assert.Equal(t, pag2.TotalRows, int64(count))
			assert.Equal(t, pag2.Page, 1000)
			rows, ok := (pag2.Rows).([]model.Pack)
			if ok {
				assert.Equal(t, len(rows), pag2.Limit)
			}
		})
	})

}
func TestCalcuatedTotalPages(t *testing.T) {
	type Case struct {
		input  model.Pagination
		expect model.Pagination
	}
	cases := []Case{
		{
			input:  model.Pagination{TotalRows: 1000, Limit: 100},
			expect: model.Pagination{TotalPages: 10},
		},
		{
			input:  model.Pagination{TotalRows: 1001, Limit: 100},
			expect: model.Pagination{TotalPages: 11},
		},
		{
			input:  model.Pagination{TotalRows: 475, Limit: 100},
			expect: model.Pagination{TotalPages: 5},
		},
	}
	for i, curCase := range cases {
		t.Run(fmt.Sprintf("Case %v", i+1), func(t *testing.T) {
			model.CalculatedTotalPages(&curCase.input)
			assert.Equal(t,
				curCase.input.TotalPages,
				curCase.expect.TotalPages,
			)
		})
	}
}
