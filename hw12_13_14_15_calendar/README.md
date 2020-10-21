#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:

-   [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
-   [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
-   [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
-   [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**

### Evironment variables to deploy (.env)

```dotenv
CAL_PRODUCTION=true         # set production flag to true; false if ommited
CAL_HTTP_PORT=9000  # HTTP API external port
CAL_GRPC_PORT=9001  # GPRC external port

# Postgres data
CAL_STORAGE_HOST=localhost
CAL_STORAGE_PORT=5432
CAL_STORAGE_USER=postgres
CAL_STORAGE_PASSWORD=password
CAL_STORAGE_DBNAME=postgres

# RabbitMQ data
CAL_RMQ_HOST=localhost
CAL_RMQ_PORT=5672
CAL_RMQ_USER=admin
CAL_RMQ_PASSWORD=secret_word
```
