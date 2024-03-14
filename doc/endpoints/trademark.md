# Get:
mark trae con sigo todas sus trademarks,
## GetAll: 

endpoint: `/trademark`

Response:
```sh
[        
    {
        "ID": 725,
        "MakerID": 2037,
        "Name": "par-0",
        "Status": true
    },
    {
        "ID": 726,
        "MakerID": 2037,
        "Name": "inpar-1",
        "Status": false
    },
    {
        "ID": 727,
        "MakerID": 2037,
        "Name": "par-2",
        "Status": true
    }    
]
```

## GetById: 
endpoint: `/trademark?ID=100` 

Response:
```json
{
    "ID": 727,
    "MakerID": 2037,
    "Name": "par-2",
    "Status": true
} 
```
## GetByName (nombre exacto)

endpoint: `/trademark?Name=example` 

Response:
```json
[{
    "ID": 100,
    "MakerID": 2037,
    "Name": "example",
    "Status": true
}]
```

# Post
El campo MakerID es obligatorio que hace referencia a la tabla Maker, los demas son opcionales.
## Insert One

endpoint: `/trademark`

request body:
```json
{"ID":0,"MakerID":2060,"Name":"example","Status":true}
```
response:
```json
{
    "ID": 1259,
    "MakerID": 2060,
    "Name": "example",
    "Status": true
}
```
## InsertMany
Permite insertar nuevos items y tambien actualizar ya existentes, para actualizar un item necesita el ID

endpoint: `/trademark`


request body:
```json
[
    {"ID":1289,"MakerID":2078,"Name":"example2","Status":false},
    {"ID":1290,"MakerID":2078,"Name":"example3","Status":true},
    {"ID":0,"MakerID":2078,"Name":"example-2","Status":true},
    {"ID":0,"MakerID":2078,"Name":"example-2","Status":false}
]
```
response:
```json
[
    {
        "ID": 1289,
        "MakerID": 2078,
        "Name": "example2",
        "Status": false
    },
    {
        "ID": 1290,
        "MakerID": 2078,
        "Name": "example3",
        "Status": true
    },
    {
        "ID": 1291,
        "MakerID": 2078,
        "Name": "example-2",
        "Status": true
    },
    {
        "ID": 1292,
        "MakerID": 2078,
        "Name": "example-2",
        "Status": false
    }
]
```

# Put

## Update
Permite actualizar items ya existentes, para actualizar un item necesita el ID
Y tambien necesita pasarse el MakerID

endpoint: `/trademark`

request body:
```json
{"ID":1736,"MakerID":2105,"Name":"example-update","Status":true}
```
response:
```json
{
    "ID": 1736,
    "MakerID": 2105,
    "Name": "example-update",
    "Status": true
}
```

# Delete:
## DeleteById
Permite borrar(soft delete) un elemento pasandole el id ejemplo 

endpoint: `/trademark`

request body:
```json
{"ID":100}
```
response:
```json
{
    "ID": 100,
    "MakerID": 0,
    "Name": "",
    "Status": false
}
```