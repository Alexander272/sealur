# Email Service

Сервис отвечающий за отправление email

    protoc -I internal/transport/grpc internal/transport/grpc/proto/email.proto --go_out=plugins=grpc:internal/transport/grpc/proto