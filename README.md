# Chat Service

Сервис чата, построенный на **gRPC-стриминге** для обмена сообщениями в реальном времени.  
Работает в Docker, использует **собственную PostgreSQL (`pg-chat`)** и **общую Kafka-инфраструктуру** от [
`auth-service`](https://github.com/ne4chelovek/auth-service).

---

## 🚀 Функционал

- **Создание чата** — с указанием участников
- **Получение информации о чате** по ID
- **Удаление чата**
- **Подключение к чату через gRPC-стриминг** — односторонний (серверный) поток
- **Отправка и получение сообщений в реальном времени**
- **REST API** через `grpc-gateway` (авто-генерация из `.proto`)
- **WebSocket-шлюз** — преобразует gRPC-стрим в WebSocket для доступа с браузера
- **JWT-аутентификация** через `auth-service`
- **Обработка событий из Kafka** — чтение топика `user_session_events` (вход/выход пользователей)
- **Swagger UI** — документация REST API

> ⚠️ **WebSocket-шлюз находится в стадии доработки**

---

## 📦 Технологии

- Go 1.23+
- gRPC + gRPC-Gateway
- gRPC-Streaming (односторонний)
- WebSocket
- PostgreSQL (своя БД — `pg-chat`)
- Kafka (consumer топика `user_session_events`)
- TLS для gRPC
- Cobra CLI — консольный интерфейс для использования чата
- Docker + Docker Compose

---

## 🧰 Перед запуском

Убедитесь, что:

- Установлены [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)
- `auth-service` запущен — он предоставляет:
    - Kafka (`kafka1`)
    - ZooKeeper
    - `auth-service` (для аутентификации)

---

## 🛠 Запуск

### Запустите `chat-service` (инфраструктура)

```bash
git clone https://github.com/ne4chelovek/chat-service
cd ../chat-service
docker-compose up --build -d
```

### После запуска сервисы доступны:

- gRPC: localhost:9070
- REST API: http://localhost:8080
- Swagger UI: http://localhost:8090
- WebSocket: ws://localhost:8080/ws/chat?chat_id=1

---

## 💻 CLI 

Расположение: chat-service/cli.

Сборка: `go build -o bin/mikle main.go`

### Доступные команды

- Помощь по CLI:

````
  bin/mikle --help
````

- Регистрация нового пользователя:

````
  bin/mikle register --username <имя> --email <email> --password <пароль>
````

- Вход (аутентификация):

````
  bin/mikle login
````

- Выход (удаление токена):

````
  bin/mikle logout
````

- Создание чата:

````
  bin/mikle create --users user1,user2,user3
````

- Удаление чата:

````
  bin/mikle delete --chat-id 5
````

- Подключение к чату (просмотр сообщений):

````
  bin/mikle connect --chat-id 5
````
 
