# Pro Service

Сервис отвечающий за логику работы серверной части Sealur Pro

#### шаблон файла config.yaml
    http:
        serviceName: "pro-service"
        maxHeaderBytes: 1
        readTimeout: 10s
        writeTimeout: 10s

    postgres:
        username: "username"
        host: "host"
        port: "port"
        dbname: "sealur_pro"
        sslmode: "disable"

#### миграции для базы данных

    docker run -v /home/martynov/Projects/SealurProBack/pro_service/app/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://postgres:qwerty@0.0.0.0:5436/sealur_pro?sslmode=disable' up
    docker run -v home/martynov/Projects/SealurProBack/pro_service/app/migrations:/migrations --network host migrate/migrate create -ext sql -dir ./migrations -seq init

#### генерация proto файлов

    protoc -I internal/transport/grpc internal/transport/grpc/proto/stand.proto --go_out=plugins=grpc:internal/transport/grpc/proto
    protoc -I internal/transport/grpc internal/transport/grpc/proto/pro.proto --go_out=plugins=grpc:internal/transport/grpc/proto
