# Api Service

Сервис API Gateway (Gateway Offloading)

    replace github.com/Alexander272/sealur_proto/api => ../../../sealur_proto/api
    scp -r ./build administrator@pro:/home/administrator/app

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

Команда для генерации документации

    swag init -g ./cmd/app/main.go
