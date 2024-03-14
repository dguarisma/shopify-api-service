La mayoria de los endpoints funcionan como crud
sin efectos secundarios
(packs, makers, suppliers, substances, trademarks, typesproducts,
categoryone, categorytwo, categorythree, warehouses)

Get:
    GetAll:
        trae todos los elementos ejemplo "/pack"
    GetBy:
        trae el elemento correpondiente con ese id ejemplo "/pack?ID=10"
    GetBy:
        trae los elementos que cumplan con un detemiado campo
        Por ejemplo si en el path termina "/pack?Name=example"
        Debera traer a todos los elementos que tengan como nombre example

Post:
    Insert:
        permite insertar un elemento
    InsertMany:
        permite insertar y actualizar varios elementos

Put:
    Update:
        Permite actualizar un elemento en concreto por defecto solo actualizar
        ese unico elemento.

Delete:
    Delete
        Permite borrar un elemento pasandole el id ejemplo `{"ID":10}`


inventory, products, purchases, reception
inventory trae las compras y ventas.
productos tiene otro tipo de busqueda igual que inventory.
esto se puede probar con los archivos de postman

para todos lo endpoints se necesita conexion a la base de datos principa `DB_URI`

para el caso compras `Purchase` se necesitan las credenciales de `aws_ses` para
enviar los mensajes a los clientes.

En el caso de `Reception` se necesitan las credenciales de las bodegas `BARRANQULLA` y `BOGOTA`.
y la base de datos se shopify `DB_SHOPIFY`
