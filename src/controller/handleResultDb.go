package controller

import (
	"desarrollosmoyan/lambda/src/response"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrNotFound  = errors.New("Elemento no encontrado")
	ErrNotUpdate = errors.New("No se ha actualizado ningun elemento")
)

func HandleResultSearchDB(result *gorm.DB) (res response.Result, problem bool) {
	if result.Error != nil {
		return response.NewFailResult(
			result.Error,
			http.StatusInternalServerError,
		), true
	}

	if result.RowsAffected == 0 {
		return response.NewFailResult(
			ErrNotFound,
			http.StatusNotFound,
		), true
	}
	return response.Result{}, false
}
