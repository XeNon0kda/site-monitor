# Site Availability Monitor

Real-time мониторинг доступности веб-сайтов с веб-интерфейсом.

## Технологии
- Go 1.25.1
- Gorilla Mux, Gorilla WebSocket
- In-memory хранилище
- Docker / docker-compose

## Архитектура
- Чистая архитектура: domain, repository, service, handler, worker
- Фоновый воркер для периодической проверки
- WebSocket для real-time обновлений
- Graceful shutdown

## Запуск

```bash
docker-compose up --build
