package shopifyserv

import (
	"fmt"
)

const (
	ErrEmptyCredentials      = "no se han aportado credenciales para la bodega de shopify [Nombre, Url, token]"
	ErrNotEnoughtCredentials = "las credenciales de la bodega de shopify no son suficiente se requiere [Nombre, Url, token]"
)

func ErrDoesntExistWarehouse(locationName string) error {
	return fmt.Errorf(
		"no existe la bodega(%s) entre las disponibles",
		locationName,
	)
}
func ErrNewRequestFormat(err error) error {
	return fmt.Errorf(
		"hubo un error en el formato de la request %s",
		err.Error(),
	)
}

func ErrResponse(err error) error {
	return fmt.Errorf(
		"hubo un error en la respuesta %s",
		err.Error(),
	)
}

func ErrResponseFormat(err error) error {
	return fmt.Errorf(
		"hubo un error al darle formato a la respuesta %s",
		err.Error(),
	)
}

func ErrStatusUnexpected(expect, got int) error {
	return fmt.Errorf(
		"hubo un error en un status inesperado esperado(%v) obtenido(%v)",
		expect, got,
	)
}
