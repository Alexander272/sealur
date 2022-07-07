# Moment Service

Сервис отвечающий за расчеты момента затяжки

    protoc -I internal/transport/grpc internal/transport/grpc/proto/moment.proto --go_out=plugins=grpc:internal/transport/grpc/proto