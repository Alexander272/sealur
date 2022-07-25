# Proto

    protoc -I internal/transport/grpc internal/transport/grpc/proto/moment.proto --go_out=plugins=grpc:internal/transport/grpc/proto

    protoc -I src/moment_api --go_out=. --go-grpc_out=. src/moment_api/moment.proto

    protoc -I=../src/ --go_out=. --go-grpc_out=. ../src/*/*.proto