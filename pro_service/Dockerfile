FROM golang:1.17-alpine as builder
WORKDIR /build
COPY app/go.mod . 
RUN go mod download
COPY app/. .
RUN go build -o cmd/app/main main.go

FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]