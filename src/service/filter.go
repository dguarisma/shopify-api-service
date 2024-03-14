package service

import (
	"encoding/json"
	"strings"
)

type IFilters interface {
	StringToArrByte(string) []byte
	AdaptSearch(queryParams map[string]string) []byte
	IsArray(body []byte) bool
	IsEmpty(body []byte) bool
	IsObject(body []byte) bool
	IsQueryParamsEmpty(queryParams map[string]string) bool
	IsSearchById(queryParams map[string]string) (string, bool)
	IsSearchByMany(queryParams map[string]string) ([]byte, bool)
}

// This ID is for check search for ID, in case the change the primary key
type Filters struct {
	IDFormat string
}

// This ID is for check search for ID, in case the change the primary key(default ID)
func NewFilter(KeyIdFromTable string) IFilters {
	return &Filters{
		IDFormat: KeyIdFromTable,
	}
}

func (Filters) IsArray(body []byte) bool {
	request := string(body)
	first := request[0]
	last := request[len(body)-1]

	return first == '[' && last == ']'
}

func (Filters) IsObject(body []byte) bool {
	first := body[0]
	last := body[len(body)-1]
	return first == '{' && last == '}'
}

func (Filters) IsEmpty(body []byte) bool {
	if len(body) <= 2 {
		return true
	}
	request := string(body)
	bodyWithBrackets := request[1 : len(request)-1]
	return strings.TrimSpace(bodyWithBrackets) == ""
}

func (Filters) StringToArrByte(body string) []byte {
	return []byte(strings.TrimSpace(body))
}

func (Filters) IsQueryParamsEmpty(queryParams map[string]string) bool {
	return len(queryParams) == 0
}

func (f *Filters) IsSearchById(queryParams map[string]string) (string, bool) {
	id, ok := queryParams["ID"]
	return id, ok
}

func (f *Filters) IsSearchByMany(queryParams map[string]string) (j []byte, ok bool) {
	if len(queryParams) > 1 {
		j, _ := json.Marshal(queryParams)
		return j, true
	}
	return nil, false
}

func (Filters) AdaptSearch(queryParams map[string]string) []byte {
	j, _ := json.Marshal(queryParams)
	return j
}
