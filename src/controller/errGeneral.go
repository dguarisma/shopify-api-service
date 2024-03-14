package controller

import (
	"desarrollosmoyan/lambda/src/response"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrFormatForID     = errors.New("Se requiere la ID")
	errImposible       = errors.New("(Imposible error) de mapeo de datos")
	errFormatIncorrect = errors.New("El formato es incorrecto")
	ErrFormat          = "Formato incorrecto"
)

func ErrIDExpected() response.Result {
	return response.NewFailResult(
		ErrFormatForID,
		http.StatusMethodNotAllowed,
	)
}

func ErrImposibleFormat() response.Result {
	return response.NewFailResult(
		errImposible,
		http.StatusInternalServerError,
	)
}

func ErrPagination(err error) response.Result {
	return response.NewFailResult(
		fmt.Errorf("error al formatear la paginacion: %v", err.Error()),
		http.StatusInternalServerError,
	)
}

func ErrFormatIncorrect(err error, expected, obtein []byte) error {
	return fmt.Errorf("%v: %v \nexpected: %s \ngot: %s",
		errFormatIncorrect,
		err,
		string(expected),
		string(obtein),
	)
}

func ErrTypeFormat(field string, example, value interface{}) error {
	return fmt.Errorf(
		"%s: el campo %s deberia tener un formato %v en vez de %v",
		ErrFormat, field, example, value,
	)
}
