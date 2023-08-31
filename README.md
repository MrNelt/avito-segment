# Окружение

Необходимо подготовить окружение:

1. заполнить ```.env.example``` (там уже выставлены дефолты)
2. выполнить ```make env```


# Тесты

Должно быть подготовлено окружение

Запуск: ```make test```

# Запуск

Должно быть подготовлено окружение

Выполнить команду ```make run```

# Остановка

Выполнить команду ```make stop```

# Остановка и удаление данных из бд

Выполнить команду ```make down```

### Возможный сценарий использования api:

```
curl --location --request POST 'localhost:8080/segment/TEST1'
```

```
curl --location --request POST 'localhost:8080/segment/TEST2'
```

```
curl --location --request POST 'localhost:8080/segment/TEST3'
```

```
curl --location --request POST 'localhost:8080/segment/TEST4'
```

```
curl --location 'localhost:8080/user/1' \
--header 'Content-Type: application/json' \
--data '{
    "add": ["TEST1", "TEST2"],
    "TTL": 2,
	"TTLUnit": "DAYS"
}'
```

```
curl --location 'localhost:8080/user/2' \
--header 'Content-Type: application/json' \
--data '{
    "add": ["TEST1", "TEST2", "TEST3", "TEST4"],
    "TTL": 2,
	"TTLUnit": "DAYS"
}'
```

```
curl --location 'localhost:8080/user/2'
```

```
{
    "ID": 2,
    "segments": [
        "TEST1",
        "TEST2",
        "TEST3",
        "TEST4"
    ]
}
```
```
curl --location 'localhost:8080/user/2'
```

```
{
    "ID": 1,
    "segments": [
        "TEST1",
        "TEST2"
    ]
}
```

```
curl --location 'localhost:8080/user/1' \
--header 'Content-Type: application/json' \
--data '{
    "add": ["TEST3", "TEST4"],
    "delete": ["TEST1", "TEST2"],
    "TTL": 2,
	"TTLUnit": "DAYS"
}'
```

```
curl --location 'localhost:8080/user/2'
```

```
{
    "ID": 1,
    "segments": [
        "TEST3",
        "TEST4"
    ]
}
```

#### Обработка ошибок

Сегмент с таким именем уже существует:

```
curl --location --request POST 'localhost:8080/segment/TEST1'
```

```
{
    "error": "ErrSegmentAlreadyExists"
}
```

Сегмент не найден:

```
curl --location 'localhost:8080/user/1' \
--header 'Content-Type: application/json' \
--data '{
    "add": ["TEST100", "TEST2"]
}'
```
```
{
    "error": "ErrSegmentNotFound"
}
```

Остальные серверные ошибки также обработаны

## TTL

из дополнительных заданий реализован TTL у юзера в сегменте
необходимо передавать в запросе ```TTL``` и ```TTLUnit```

Пример:

```
"TTL": 2,
"TTLUnit": "DAYS"
```

```TTLUnit in ["DAYS", "HOURS", "MINUTES", "SECONDS"]```