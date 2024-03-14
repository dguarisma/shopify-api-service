# shopify-api-service

## Serverless para Actualización de Inventario y Órdenes en Shopify
Este proyecto utiliza arquitectura serverless para actualizar el inventario y gestionar órdenes en una tienda Shopify. La implementación se realiza utilizando servicios en la nube, lo que permite una escalabilidad eficiente y un manejo simplificado de la infraestructura.

#### Funcionalidades Principales
El sistema proporciona las siguientes funcionalidades principales:

Actualización de Inventario: Utiliza eventos de cambio de inventario para mantener actualizados los niveles de stock en la tienda Shopify.

Gestión de Órdenes: Procesa las órdenes entrantes de manera automatizada, actualizando el estado de las órdenes en la tienda Shopify y realizando acciones adicionales según sea necesario.



# Primero installar las dependecias
```go
$ go mod tidy
```
despues installar docker y docker-compose
url -> `https://docs.docker.com/engine/install/`
url -> `https://docs.docker.com/compose/install/linux/#install-using-the-repository`

ejecutar los contenedores en segundo plano
```bash
$ docker-compose up -d
```
en caso de querer detenerlos
```bash
$ docker-compose stop
```

en el puerto 8080 -> `http://localhost:8080/`
estara corriendo un gestor de base de datos lijero
servidor:      "mysql"
usuario:       "root"
contraseña:    "123456"
base de datos: "defaultdb"

primero hay q ingresar manualmente

para montar la base de datos,
en la raiz del proyecto ejecutar
```sh
$ go run ./src/tests/main.go
```
una ver realizado

para dar un ejemplo testearemos
si se quiere probar los test ir a
`./src/tests/<nombre del metodo que quiera hacer el test>`
y una vez en la carpeta ejecutar go test ./<nombre_de_los_test>.go -v
asi tambien se puede ver metodos(post, get, put, delete) y como deben ser
enviados

si se quiere habilitar ejemplos para ver las request en los test hay que
habilitar el show.
