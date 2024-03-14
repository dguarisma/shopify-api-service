package service_test

import (
	"desarrollosmoyan/lambda/src/response"
	"desarrollosmoyan/lambda/src/service"
	"net/http"
	"testing"
)

const (
	getAll        = "getAll"
	getByID       = "getByID"
	getManyBy     = "getManyBy"
	insertOne     = "insertOne"
	insertMany    = "insertMany"
	update        = "update"
	deleteById    = "deleteById"
	defaultStatus = 200
)

type mockRepo struct{}

func (m *mockRepo) GetByID(id string) response.Result {
	return response.NewResult(getByID, defaultStatus)
}

func (m *mockRepo) GetManyBy(pathParam []byte) response.Result {
	return response.NewResult(getManyBy, defaultStatus)
}

func (m *mockRepo) GetAll() response.Result {
	return response.NewResult(getAll, defaultStatus)
}

func (m *mockRepo) InsertOne(reqBody []byte) response.Result {
	return response.NewResult(insertOne, defaultStatus)
}

func (m *mockRepo) InsertMany(reqBody []byte) response.Result {
	return response.NewResult(insertMany, defaultStatus)
}

func (m *mockRepo) Update(reqBody []byte) response.Result {
	return response.NewResult(update, defaultStatus)
}

func (m *mockRepo) DeleteById(reqBody []byte) response.Result {
	return response.NewResult(deleteById, defaultStatus)
}

func errPrint2(t *testing.T, input interface{}, expected, output response.Result) {
	ErrMsg := "\n\tInput      :  %v"
	ErrMsg += "\n\tExpect Body: %s"
	ErrMsg += "\n\tOutput Body: %s"

	if expected.Body != output.Body {
		t.Errorf(ErrMsg, input, expected.Body, output.Body)
		return
	}

	ErrMsg = "\n\tInput        : %v"
	ErrMsg += "\n\tExpect Status: %d"
	ErrMsg += "\n\tOutput Status: %d"

	if expected.Status != output.Status {
		t.Errorf(ErrMsg, input, expected.Status, output.Status)
		return
	}
}

func TestService(t *testing.T) {
	idTable := "ID"
	var filter service.IFilters = service.NewFilter(idTable)

	repo := mockRepo{}

	serv := service.NewService(filter, &repo)
	t.Run("Check Gets", func(t *testing.T) {
		type CaseCheckGet struct {
			title       string
			queryParams map[string]string
			expected    response.Result
		}

		cases := []CaseCheckGet{
			{
				title:       "It's should to be getAll",
				queryParams: map[string]string{},
				expected:    response.NewResult(getAll, defaultStatus),
			},
			{
				title:       "It's should to be getByID",
				queryParams: map[string]string{idTable: "100"},
				expected:    response.NewResult(getByID, defaultStatus),
			},
			{
				title:       "It's should to be getManyBy",
				queryParams: map[string]string{"Name": "Maria"},
				expected:    response.NewResult(getManyBy, defaultStatus),
			},
			{
				title:       "It's should to be getByID",
				queryParams: map[string]string{"iD": "100"},
				expected:    response.NewResult(getManyBy, defaultStatus),
			},
		}

		for _, curCase := range cases {
			t.Run(curCase.title, func(t *testing.T) {
				res := serv.Get(curCase.queryParams)
				errPrint2(t, curCase.queryParams, curCase.expected, res)
			})
		}

	})

	t.Run("Check Insert Cases", func(t *testing.T) {
		type CaseInsert struct {
			title    string
			input    string
			expected response.Result
		}

		cases := []CaseInsert{
			{
				title:    "It's should return empty error",
				input:    "",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "{}",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "[]",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return format error",
				input:    "trying to insert something.",
				expected: response.NewFailResult(service.ErrNotJsonOrArray, http.StatusBadRequest),
			},
			{
				title:    "It's should return insertOne and defaultStatus",
				input:    `{"Name":"Maria"}`,
				expected: response.NewResult(insertOne, defaultStatus),
			},
			{
				title:    "It's should return insertMany and defaultStatus",
				input:    `[{"Name":"Maria"}, {"Name":"Maria2"}]`,
				expected: response.NewResult(insertMany, defaultStatus),
			},
		}

		for _, curCase := range cases {
			t.Run(curCase.title, func(t *testing.T) {
				res := serv.Insert(curCase.input)
				errPrint2(t, curCase.input, curCase.expected, res)
			})
		}
	})

	t.Run("Check Update Cases", func(t *testing.T) {
		type CaseUpdate struct {
			title    string
			input    string
			expected response.Result
		}

		cases := []CaseUpdate{
			{
				title:    "It's should return empty error",
				input:    "",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "{}",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "[]",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return format error json",
				input:    "trying to insert something.",
				expected: response.NewFailResult(service.ErrNotJson, http.StatusBadRequest),
			},
			{
				title:    "It's should return format error json",
				input:    `[{"Name":"Maria"}, {"Name":"Maria2"}]`,
				expected: response.NewFailResult(service.ErrNotJson, http.StatusBadRequest),
			},
			{
				title:    "It's should return update and defaultStatus",
				input:    `{"Name":"Maria"}`,
				expected: response.NewResult(update, defaultStatus),
			},
		}

		for _, curCase := range cases {
			t.Run(curCase.title, func(t *testing.T) {
				res := serv.Update(curCase.input)
				errPrint2(t, curCase.input, curCase.expected, res)
			})
		}
	})

	t.Run("Check Delete Cases", func(t *testing.T) {
		type CaseDelete struct {
			title    string
			input    string
			expected response.Result
		}

		cases := []CaseDelete{
			{
				title:    "It's should return empty error",
				input:    "",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "{}",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return empty error",
				input:    "[]",
				expected: response.NewFailResult(service.ErrEmptyBody, http.StatusBadRequest),
			},
			{
				title:    "It's should return format error json",
				input:    "trying to insert something.",
				expected: response.NewFailResult(service.ErrNotJson, http.StatusBadRequest),
			},
			{
				title:    "It's should return format error json",
				input:    `[{"Name":"Maria"}, {"Name":"Maria2"}]`,
				expected: response.NewFailResult(service.ErrNotJson, http.StatusBadRequest),
			},
			{
				title:    "It's should return update and defaultStatus",
				input:    `{"ID":"01"}`,
				expected: response.NewResult(deleteById, defaultStatus),
			},
		}

		for _, curCase := range cases {
			t.Run(curCase.title, func(t *testing.T) {
				res := serv.Delete(curCase.input)
				errPrint2(t, curCase.input, curCase.expected, res)
			})
		}
	})
}
