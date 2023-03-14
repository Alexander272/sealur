# Pro Service

Сервис отвечающий за логику работы серверной части Sealur Pro

go mod edit -replace github.com/Alexander272/sealur_proto/api ../../../sealur_proto/api
replace github.com/Alexander272/sealur_proto/api => ../../../sealur_proto/api

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
    docker run -v home/martynov/Projects/SealurProBack/pro_service/app/migrations:/migrations --network host migrate/migrate create -ext sql -dir /migrations -seq init

#### генерация proto файлов

    protoc -I internal/transport/grpc internal/transport/grpc/proto/pro.proto --go_out=plugins=grpc:internal/transport/grpc/proto

    openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout localhost.key -out localhost.pem -subj "/C=RU/CN=localhost"
    openssl x509 -outform pem -in localhost.pem -out localhost.crt
