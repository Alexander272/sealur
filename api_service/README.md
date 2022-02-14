# Api Service

Сервис API Gateway (Gateway Offloading)

<b>папка proto</b>
содержит сгнерирование proto-файлы (появилась потому что, не работают импроты, и как заставить их работать я пока не знаю)

#### шаблон файла config.yaml

    http:
    serviceName: "api-service"
    host: "http://localhost"
    port: "port"
    maxHeaderBytes: 1
    readTimeout: 10s
    writeTimeout: 10s

    cache:
        ttl: 3600s

    auth:
        accessTokenTTL: ttl
        refreshTokenTTL: ttl

    redis:
        host: "host"
        port: "port"
        db: 0

    limiter:
        rps: 10
        burst: 20
        ttl: ttl
