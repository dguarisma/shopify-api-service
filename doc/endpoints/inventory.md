En el endpoint `inventory` solo funciona el metodo `GetManyBy`

Get:
el `WarehouseID` es obligatorio.

path: `/inventory?WarehouseID=20`
devuelve las compras y ventas de la bodega `20` en forma paginada.

path: `/inventory?WarehouseID=20&Page=2`
devuelve las compras y ventas de la bodega `20` en forma paginada, pag 2


path: `/inventory?WarehouseID=20&Page=2&From=01/11/2023&To=02/11/2023`
devuelve las compras y ventas de la bodega `20` en forma paginada, pag 2 entre
las fechas `01/11/2023` y `02/11/2023`
si esta el campo from es obligatorio el campo to, en caso contrario sera
ignorado.

path: `/inventory?WarehouseID=20&Page=2&From=01/11/2023&To=02/11/2023&Sku=a1s230s`
devuelve las compras y ventas de la bodega `20` en forma paginada, pag 2 entre
las fechas `01/11/2023` y `02/11/2023` del producto con el sku `a1s230s`
