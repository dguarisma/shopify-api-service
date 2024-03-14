# Get:
mark trae con sigo todas sus trademarks,
## GetAll: 

endpoint: `/maker`

Response:
```sh
[
    {
        "ID": 2,
        "Name": "par-0",
        "Status": true,
        "Trademarks": [
            {
                "ID": 4,
                "MakerID": 1012,
                "Name": "trademark-04",
                "Status": false
            }, {
                "ID": 7,
                "MakerID": 1012,
                "Name": "trademark-07",
                "Status": false
            }
        ]
    },
    {
        "ID": 3,
        "Name": "inpar-1",
        "Status": false,
        "Trademarks": [
            {
                "ID": 5,
                "MakerID": 3,
                "Name": "trademark-05",
                "Status": false
            }
        ]
    },
    {
        "ID": 4,
        "Name": "par-2",
        "Status": true,
        "Trademarks": [
            {
                "ID": 6,
                "MakerID": 4,
                "Name": "trademark-06",
                "Status": false
            }
        ]
    }
]
```

## GetById: 
endpoint: `/maker?ID=100` \
Response:
```json
{
    "ID": 100,
    "Name": "par-50",
    "Status": true,
    "Trademarks": [
        {
            "ID": 200,
            "MakerID": 100,
            "Name": "trademark-200",
            "Status": true
        }
    ]
}
```
## GetByName (nombre exacto)

endpoint: `/maker?Name=example` 

Response:
```json
[{
    "ID": 1513,
    "Name": "example",
    "Status": false,
    "Trademarks": []
}]
```

# Post
## Insert One
Ningun campo es obligatorio pero debe tener al menos uno(Name, Status)


endpoint: `/maker`

request body:
```json
{"ID":0,"Name":"example","Status":false,"Trademarks":null}
```
response:
```json
{
    "ID": 2012,
    "Name": "example",
    "Status": false,
    "Trademarks": null
}
```
## InsertMany
Permite insertar nuevos items y tambien actualizar ya existentes, para actualizar un item necesita el ID

endpoint: `/maker`


request body:
```json
[
    {"ID":2018,"Name":"example2","Status":false,"Trademarks":null},
    {"ID":2019,"Name":"example3","Status":true,"Trademarks":null},
    {"ID":0,"Name":"example-2","Status":true,"Trademarks":null},
    {"ID":0,"Name":"example-2","Status":false,"Trademarks":null}
]
```
response:
```json
[
    {
        "ID": 2018,
        "Name": "example2",
        "Status": false,
        "Trademarks": null
    },
    {
        "ID": 2019,
        "Name": "example3",
        "Status": true,
        "Trademarks": null
    },
    {
        "ID": 2020,
        "Name": "example-2",
        "Status": true,
        "Trademarks": null
    },
    {
        "ID": 2021,
        "Name": "example-2",
        "Status": false,
        "Trademarks": null
    }
]
```

# Put

## Update
Permite actualizar items ya existentes, para actualizar un item necesita el ID

endpoint: `/maker`

request body:
```json
{"ID":2022,"Name":"example-update","Status":true,"Trademarks":null}
```
response:
```json
{
    "ID": 2022,
    "Name": "example-update",
    "Status": true,
    "Trademarks": null
}
```

# Delete:
## DeleteById
Permite borrar(soft delete) un elemento pasandole el id ejemplo 

endpoint: `/maker`

request body:
```json
{"ID":100}
```
response:
```json
{
    "ID": 100,
    "Name": "",
    "Status": false,
    "Trademarks": null
}
```