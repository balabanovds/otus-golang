#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:

-   [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
-   [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
-   [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
-   [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**

### Evironment variables to run calendar and scheduler

```dotenv
CAL_PRODUCTION=true         # set production flag to true; false if ommited
CAL_STORAGE_DSN="some DSN"  # full Data Source Name to connect to PostgreSQL
```

### Additional environment variables for scheduler

```dotenv
CAL_RMQ_LOGIN=login         # default 'guest'
CAL_RMQ_PASSWORD=password   # default 'guest'
```
