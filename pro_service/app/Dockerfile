FROM golang:1.19.8-alpine3.17 as builder
WORKDIR /build
COPY ./go.mod . 
RUN go mod download
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:3
RUN apk add --no-cache tzdata
COPY ./configs/config.yaml /configs/config.yaml
COPY ./template /template
COPY ./cert /cert
COPY --from=builder /build/main /bin/main
ENTRYPOINT ["/bin/main"]