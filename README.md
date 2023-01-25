# cmdb functions

## Type (ex. Profiletype)
- **`kafka topic (ingress input)- type`**
- **`kafka key - type name`**
- **`kafka message`**
```
{
    "method": "create",
    "name": "nodes",
    "payload": {}
}
```
- **method** - может принимать значения (create, createChild, read, update, delete)
- **payload** - произвольный объект, будет как единственный профиль в cmdb
- **name** - имя линка


### Object (ex. MgmtObject)
- **`kafka topic (ingress input)- object`**
- **`kafka key - object uuid`**
- **`kafka message`**
```
{
    "method": "createChild",
    "name": "nodes",
    "type": "node",
    "payload": {}
}
```
- **method** - может принимать значения (createChild, read, update, delete)
- **type** - произвольный тип
- **payload** - произвольный объект, будет как единственный профиль в cmdb (обязателен при createChild/update)
- **name** - имя линка


## Link
- **`kafka topic (ingress input)- link`**
- **`kafka key - from uuid or link uuid`**
- **`kafka message`**
```
{
    "method": "create",
    "name": "nvme0n1",
    "to": "aec5b90c-1922-4ec5-a448-76c4710fc4e2",
    "payload": {}
}
```
- **method** - может принимать значения (create, read, update, delete)
- **name** - имя линка
- **to** - uuid объекта к которому линк
- **payload** - произвольный объект, будет как единственный профиль в cmdb

