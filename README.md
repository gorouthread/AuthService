# AuthService
Ммикросервис аутентификации, реализующий единый источник истины для идентификации пользователей и выдачи токенов.

## Ключевые особенности
- Stateless‑подход с короткими access‑токенами и серверными refresh‑токенами в Redis 
- Поддержка ротации refresh‑токенов; 
- Готовность к горизонтальному масштабированию в Kubernetes.

## Аудитория
Web‑приложения, мобильные клиенты, внутренние API‑сервисы

## Запуск микросервиса через docker-compose 
Запуск
```bash
    make env-up
    make migrate-up
    make auth-up
```
Быстрые curl тесты (осторожно, данные попадают в бд, собираются метрики)
```bash
    make bash-test
```
Unit тесты service слоя с mock зависимостями
```bash
    make auth-test
```

## Переменные окружения
Переменные окружения находятся в файле .env.example.

## REST API endpoints
|      Endpoint         | Method |          Description           |
|-----------------------|--------|--------------------------------|
| /api/v1/auth/register | POST   | регистрация пользователя       |   
| /api/v1/auth/login    | POST   | аутентификация пользователя    |    
| /api/v1/auth/refresh  | POST   | получение нового access токена |
| /api/v1/auth/logout   | POST   | завершить сеанс пользователя   |  

## Метрики
GUI Grafana доступен на localhost:3000
GUI Prometheus доступен на localhost:9090
Приложение отображает в Grafana на dashbourd следующие метрики: rps, total request, request by status.

## Логирование
Логи приложения доступны в Grafana Loki: Grafana -> Explore -> Источник:Локи -> Код {job="docker"}
Сбор логов осущетсвляется с помощью агента alloy из всех docker контейнеров.

## Swagger 
TODO

## Используемые технологии
- Go (основной язык программирования)
- PostgreSQL (основная база данных)
- Redis (кеширование)
- Docker (контейнеризация приложения)
- Make (автоматизация сборки)
- Bash (скриптовые мини тесты curl)
- Prometheus (метрики)
- Grafana (dashbourd метрик)
- Loki (логи)