package model

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type AdaptPagination struct {
	Limit int    `json:"limit,string"`
	Page  int    `json:"page,string"`
	Sort  string `json:"sort"`
}

func GetPagination(b []byte) (*Pagination, error) {
	adapt := AdaptPagination{}
	pagination := &Pagination{}

	fmt.Print("3")
	if err := json.Unmarshal(b, &adapt); err != nil {
		return nil, err
	}

	fmt.Print("4")
	if adapt.Limit == 0 && adapt.Page == 0 && adapt.Sort == "" {
		return pagination, nil
	}

	pagination.Limit = adapt.Limit
	pagination.Sort = adapt.Sort
	pagination.Page = adapt.Page

	return pagination, nil
}

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,string,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"Rows"`
}

func (p *Pagination) IsEmpty() bool {
	if p.Limit != 0 {
		return false
	}

	if p.Page != 0 {
		return false
	}

	if p.Sort != "" {
		return false
	}
	return true
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func PaginateBy(value interface{}, where *gorm.DB, pagination *Pagination) (*Pagination, error) {
	if err := where.Count(&pagination.TotalRows).Error; err != nil {
		return nil, err
	}

	CalculatedTotalPages(pagination)
	if pagination.TotalPages == 0 {
		pagination.TotalPages = 1
	}

	err := where.
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order(pagination.GetSort()).Find(value).Error // luego ver

	if err != nil {
		return nil, err
	}

	pagination.Rows = value
	return pagination, nil
}

func CalculatedTotalPages(pag *Pagination) {
	if pag.TotalRows%int64(pag.GetLimit()) != 0 {
		pag.TotalPages = 1
	}
	pag.TotalPages += int(pag.TotalRows / int64(pag.GetLimit()))
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	db.Model(value).Debug().Count(&pagination.TotalRows)

	CalculatedTotalPages(pagination)
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Offset(pagination.GetOffset()).
			Limit(pagination.GetLimit()).
			Order(pagination.GetSort())
	}
}
