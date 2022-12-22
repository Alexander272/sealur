# Moment Service

Сервис отвечающий за расчеты момента затяжки

    replace github.com/Alexander272/sealur_proto/api => ../../../sealur_proto/api

    protoc -I internal/transport/grpc internal/transport/grpc/proto/moment.proto --go_out=plugins=grpc:internal/transport/grpc/proto