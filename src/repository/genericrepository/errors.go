package genericrepository

import "fmt"

const (
	ElementsNotFounds = "No se ha encontrado registros en la base de datos"
	elementNotFound   = "No se ha encontrado registro del elemento con el id(%v)"
)

func ElementNotFoundById(id string) error {
	return fmt.Errorf(elementNotFound, id)
}
