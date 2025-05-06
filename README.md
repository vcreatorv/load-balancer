# Load Balancer

Этот проект представляет собой HTTP-балансировщик нагрузки, реализованный на языке Go. Он поддерживает алгоритмы round-robin и least connections для распределения запросов между бэкенд-серверами.

## Основные функции

- **Балансировка нагрузки**: Поддержка алгоритмов round-robin и least connections.
- **Управление бэкенд-серверами**: Добавление и удаление серверов через API.
- **Управление алгоритмом балансировки**: Изменение алгоритма балансировки через API.
- **Конфигурация**: Настройка через параметры командной строки.
- **Логирование**: Стандартное логирование с помощью пакета `log`.
- **Graceful Shutdown**: Корректное завершение работы при получении сигналов SIGINT или SIGTERM.
- **Обработка ошибок**: Возврат структурированных JSON-сообщений об ошибках.

## Установка и запуск

### Шаги установки

1. **Склонируйте репозиторий:**

    ```bash
    git clone https://github.com/vcreatorv/load-balancer.git
    cd load-balancer
    ```

2. **Установите зависимости:**

    ```bash
    go mod tidy
    ```

3. **Запустите балансировщик:**

    ```bash
    go run ./cmd/app/main.go --port=8090 --servers="http://127.0.0.1:8081,http://127.0.0.1:8082"
    ```

   **Флаги командной строки:**

    - `--port`: Порт для прослушивания балансировщиком (по умолчанию: 8080).
    - `--servers`: Список URL-адресов бэкенд-серверов, разделенных запятыми.

4. **Запустите бэкенд-серверы:**

    ```bash
    go run ./cmd/pool/main.go
    ```

   Этот скрипт запускает три сервера на портах 8081, 8082 и 8083.

## API

- **Перенаправление запроса:**

    ```
    GET http://localhost:8090/loadbalancer/hello
    ```

  Пример ответа: `Hello from Server 1!`

- **Добавление сервера:**

    ```
    POST http://localhost:8090/loadbalancer/backend/add
    ```

  Тело запроса:

    ```json
    {
        "server_url": "http://localhost:8083/"
    }
    ```

- **Удаление сервера:**

    ```
    DELETE http://localhost:8090/loadbalancer/backend/delete
    ```

  Тело запроса:

    ```json
    {
        "server_url": "http://localhost:8081/"
    }
    ```

- **Изменение алгоритма балансировки:**

    ```
    POST http://localhost:8090/loadbalancer/algorithm/set
    ```

  Тело запроса:

    ```json
    {
        "algorithm": "round-robin"
    }
    ```

## Обработка ошибок

Примеры сообщений об ошибках:

- **Backend не найден:**

    ```json
    {
        "status": 404,
        "message": "Backend not found (http://127.0.0.1:8081)"
    }
    ```

- **Backend уже существует:**

    ```json
    {
        "status": 409,
        "message": "A backend with this url already exists (http://127.0.0.1:8083)"
    }
    ```
  
- **Неправильный формат данных:**
    ```json
    {
        "status": 400,
        "message": "server_url: not_a_url does not validate as url"
    }
    ```