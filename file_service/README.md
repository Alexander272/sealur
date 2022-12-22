# File Service

Сервис отвечающий за хранение файлов

    protoc -I internal/transport/grpc internal/transport/grpc/proto/file.proto --go_out=plugins=grpc:internal/transport/grpc/proto