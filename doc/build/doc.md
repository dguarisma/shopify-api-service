Primero se tiene que elegir la lambda

se selecciona en el archivo `src/server/server.go`

```go
// ./src/server/server.go
package server

import ...

func Lambda() {
	db, errDb := database.GetDb()
	lambda.Start(HandleRequest2(
		db,
		// Product,
		// Reception,
		// Inventory,
		// CategoryOne,
		// CategoryTwo,
		// CategoryThree,
		Reception,
		// Substance,
		// Supplier,
		// Typesproduct,
		// Warehouse,
		// Pack,
		// Trademark,
		// Maker,
		// Purchase,
		errDb,
	))
}

...
```
en el ejemplo esta seleccionado reception y yo quiero product solo tengo q descomentar `Product` y comentar 
`Reception`
```go
// ./src/server/server.go
package server

import ...

func Lambda() {
	db, errDb := database.GetDb()
	lambda.Start(HandleRequest2(
		db,
		Product,
		// Reception,
		// Inventory,
		// CategoryOne,
		// CategoryTwo,
		// CategoryThree,
		// Reception,
		// Substance,
		// Supplier,
		// Typesproduct,
		// Warehouse,
		// Pack,
		// Trademark,
		// Maker,
		// Purchase,
		errDb,
	))
}
...
```
# En linux

luego estando al nivel del `comprimir.sh` ejecutar sus comandos 
para poder ejecutarlo en linux
```sh
$ chmod +x comprimir.sh
```
```sh
$ ./comprimir.sh
```
y compilara el codigo en un binario de nombre `main` y un `main.zip` que sera el que se envie a la lambda.

En las fotos estan los paso a paso para subir el `main.zip`