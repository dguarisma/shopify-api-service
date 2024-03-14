package service_test

import (
	"desarrollosmoyan/lambda/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var filt service.IFilters = &service.Filters{}

	type CaseCheck struct {
		expected bool
		reason   string
		body     string
	}

	t.Run("Check Emptys", func(t *testing.T) {
		casesFilt := []CaseCheck{
			{true, "it's empty", ""},
			{true, "it's empty", " "},
			{true, "it's empty", "  "},
			{true, "it's empty", "   "},
			{true, "it's empty", "[]"},
			{true, "it's empty", "{}"},
			{false, "it is not empty", `{"id":""}`},
		}

		for _, curCase := range casesFilt {
			body := filt.StringToArrByte(curCase.body)
			isEmpty := filt.IsEmpty(body)
			if isEmpty != curCase.expected {
				assert.Equal(t,
					curCase.expected,
					isEmpty,
					curCase.reason,
				)

			}
		}
	})

	t.Run("Check Arrays", func(t *testing.T) {
		casesFilt := []CaseCheck{
			{false, "it is not array", ""},
			{false, "it is not array", "  {}"},
			{false, "it is not array", "[  }"},
			{false, "it is not array", "{] "},
			{false, "it is not array", "["},
			{false, "it is not array", "[]"},
			{false, "it is not array", "[                  ]"},
			{false, "it is not array", `{"id":""}`},
			{true, "it's array", `[{"id":""}]`},
		}

		for _, curCase := range casesFilt {
			body := filt.StringToArrByte(curCase.body)
			if filt.IsEmpty(body) {
				continue
			}
			isArray := filt.IsArray(body)

			assert.Equal(t,
				curCase.expected,
				isArray,
				curCase.reason,
			)
		}
	})

	t.Run("Check Object", func(t *testing.T) {

		casesFilt := []CaseCheck{
			{false, "it's not an object", ""},
			{false, "it's not an object", "  {}"},
			{false, "it's not an object", "[  }"},
			{false, "it's not an object", "{] "},
			{false, "it's not an object", "["},
			{false, "it's not an object", "[]"},
			{false, "it's not an object", `[{"id":""}]`},
			{true, "it's an object", `{"id":""}`},
		}

		for _, curCase := range casesFilt {
			body := filt.StringToArrByte(curCase.body)
			if filt.IsEmpty(body) {
				continue
			}
			IsObject := filt.IsObject(body)

			assert.Equal(t,
				curCase.expected,
				IsObject,
				curCase.reason,
			)
		}
	})

	t.Run("Check QueryParamsEmpty", func(t *testing.T) {
		type CaseQueryParamsEmpty struct {
			expected bool
			input    map[string]string
		}

		queryCases := []CaseQueryParamsEmpty{
			{true, map[string]string{}},
			{false, map[string]string{"ID": "1"}},
			{false, map[string]string{"Name": "Maria"}},
			{false, map[string]string{"Sku": "0101"}},
		}

		for _, curCase := range queryCases {
			isEmpty := filt.IsQueryParamsEmpty(curCase.input)

			assert.Equal(t,
				curCase.expected,
				isEmpty,
			)
		}
	})
}
