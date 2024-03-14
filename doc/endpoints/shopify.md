se opto por utilizar los metodos graphql de shopify

Esta funcionalidad es activada cada vez que se realiza una o varias inserciones
en el endpoint `reception` o cuando se borra algun elemento en el mismo(metodo `delete`)

Lo primero que hace es cargar las credenciales
cada bodega necesita tres datos para enviar actualizaciones al inventory se shopify,
el nombre, la url y XshopifyAccessToken.

sin estas credenciales directamente no permite ninguna accion el endpoint `reception`

Desde reception se llama el metodo `SendUpdateShopify` que pertenece la interfaz
`IHandlerPurchaseStatus` esta es implementada por `HandlerPurchaseStatus`

este metodo llama a otro metodo interno llamado `getCityAndProductShopifyId`
que lo que hace es obtener la ciudad de la bodega `warehouse.City` correspondiente
a la compra, luego toma el ID correspondiente al producto dependiendo de la bodega
lo tomara de `product.HandlesBaq` o `product.HandlesBog`,
el primero seria si la ciudad de la bodega es `Barranquilla` y el segundo si fuera `Bogotá`.

Si no falta ningun dato tienen que existir, la compra(purchase), la bodega(warehouse), el producto(product)

se genera una estructura con la ciudad de la bodega,
el id correspondiente al producto en shopify y la cantidad en la recepcion.
que sera enviada al metodo `ShopifyServ.SendUpdateInventory`
este chequea si existe la ciudad, entre los datos de las credenciales.
si no esta, aborta el proceso y devuelve `"no existe la bodega([nombre no admitido]) entre las disponibles"`
Si existe prepara las credenciales para las request.

1. utiliza el id obtenido de `product.HandlesBaq` o `product.HandlesBog`,
para hacerle una peticion al metodo `GetItemInventory` de shopify para obtener
el `InventoryItemID` porque es necesario para actualizar el inventory de shopify.
Si hay un error devolvera el mensaje `"hubo un error en la peticion seguramente
este ProductId no corresponde a la bodega"`,
puede ser que el id no sea el correcto, o que no exista
entonces no se podra actualizar en el inventario.

2. se formatea otra consulta con el `InventoryItemID` y la cantidad, para el metodo
`InventoryBulkAdjustQuantitiesAtLocationMutation` de shopify.

pasado estos dos paso se actualiza el inventory del item en shopify.