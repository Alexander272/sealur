FROM golang:1.17-alpine as builder
WORKDIR /build
COPY ./go.mod . 
RUN go mod download
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:3
COPY ./configs/config.yaml /configs/config.yaml
COPY ./cert /cert
COPY --from=builder /build/main /bin/main
ENTRYPOINT ["/bin/main"]