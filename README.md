# level 0
## Задание для WildBerries

Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. Модель данных в формате JSON прилагается к заданию.	

### Что нужно сделать:
- Развернуть локально PostgreSQL

Создать свою БД

Настроить своего пользователя

Создать таблицы для хранения полученных данных

- Разработать сервис

Реализовать подключение и подписку на канал в nats-streaming

Полученные данные записывать в БД

Реализовать кэширование полученных данных в сервисе (сохранять in memory)

- В случае падения сервиса необходимо восстанавливать кэш из БД

Запустить http-сервер и выдавать данные по id из кэша

Разработать простейший интерфейс отображения полученных данных по id заказа

### Советы				
- Данные статичны, исходя из этого подумайте насчет модели хранения в кэше и в PostgreSQL. Модель в файле model.json
  
- Подумайте как избежать проблем, связанных с тем, что в канал могут закинуть что-угодно
  
- Чтобы проверить работает ли подписка онлайн, сделайте себе отдельный скрипт, для публикации данных в канал

- Подумайте как не терять данные в случае ошибок или проблем с сервисом

- Nats-streaming разверните локально (не путать с Nats)

Модель json-данных:
```console
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```
## Запуск приложения

Сначала клонируем репозиторий:

```console
git clone git@github.com:SaRu621/level0.git
```

Для запуска приложения необходимо запустить **nats streaming server**:
```console
cd tmp/nats-streaming-server-v0.14.2-linux-amd64
./nats-streaming-server
```
Далее запуск самого приложения:

```console
cd level0/cmd/
go run main.go
```

Запуск графического интерфейса:

```console
cd level0/interface/
python3 client.py
```

Запуск скрипта-генератора:

```console
cd level0/generation_script/
go run generator.go
```

Переменные окружения находятся в .env . В нем можно изменить параметры подключения к БД и изменить порт, на котором работает http-роутер.
