# Api Service

Сервис API Gateway (Gateway Offloading)

    replace github.com/Alexander272/sealur_proto/api => ../../../sealur_proto/api
    scp -r ./build administrator@pro:/home/administrator/app

#### шаблон файла config.yaml

    http:
        serviceName: 'api-service'
        host: localhost
        port: 8080
        maxHeaderBytes: 1
        readTimeout: 10s
        writeTimeout: 10s
        domain: localhost
        link: http://localhost

    cache:
        ttl: 3600s

    auth:
        accessTokenTTL: ttl
        refreshTokenTTL: ttl
        limitAuthTTL: ttl
        countAttempt: 5
        confirmTTL: ttl
        secure: true

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
