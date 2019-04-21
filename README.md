#### Используемые по умолчанию порты
- 24001 - HTTP API

#### Docker
Сборка образа:

    docker build -t pocoz/auto-builder .

Push образа:

    docker push pocoz/auto-builder

Запуск образа:

    docker run -d --network host --restart always --name auto-builder pocoz/auto-builder
